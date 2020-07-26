package res

import (
	stderr "errors"
	"github.com/funnyecho/code-push/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorWithStatusCode(c *gin.Context, status int, err error) {
	res(c, bodyErrorCodeMiddleware(err), errorStacksMiddleware(err), statusCodeMiddleware(status))
}

func Error(c *gin.Context, err error) {
	res(c, bodyErrorCodeMiddleware(err), errorStacksMiddleware(err), statusCodeMiddleware(http.StatusBadRequest))
}

func Success(c *gin.Context, data interface{}) {
	res(c, bodyDataMiddleware(data), statusCodeMiddleware(http.StatusOK))
}

func res(c *gin.Context, fns ...resOptionsFn) {
	statusCode := http.StatusOK
	body := make(gin.H)

	for _, fn := range fns {
		fn(c, &statusCode, body)
	}

	c.JSON(statusCode, body)
}

type resOptionsFn func(c *gin.Context, statusCode *int, body gin.H)

func errorStacksMiddleware(err error) resOptionsFn {
	return func(c *gin.Context, statusCode *int, body gin.H) {
		for cause := err; cause != nil; cause = stderr.Unwrap(err) {
			parsedError, ok := err.(*gin.Error)

			if ok {
				_ = c.Error(parsedError)
			} else {
				var reasonableErr *errors.Error

				isReasonableCause := stderr.As(cause, reasonableErr)
				causeType := gin.ErrorTypePrivate

				if isReasonableCause {
					causeType = gin.ErrorTypePublic
				}

				_ = c.Error(&gin.Error{
					Err:  err,
					Type: causeType,
				})
			}
		}
	}
}

func bodyErrorCodeMiddleware(err error) resOptionsFn {
	return func(c *gin.Context, statusCode *int, body gin.H) {
		var reasonableErr *errors.Error

		isReasonableErr := stderr.As(err, reasonableErr)
		if !isReasonableErr {
			*reasonableErr = "FA_INTERNAL_ERROR"
		}

		body["code"] = reasonableErr.Error()
	}
}

func bodyDataMiddleware(data interface{}) resOptionsFn {
	return func(c *gin.Context, _ *int, body gin.H) {
		if data == nil {
			return
		}

		body["data"] = data
	}
}

func statusCodeMiddleware(code int) resOptionsFn {
	return func(c *gin.Context, statusCode *int, _ gin.H) {
		*statusCode = code
	}
}
