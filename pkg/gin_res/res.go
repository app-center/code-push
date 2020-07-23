package res

import (
	stderr "errors"
	"github.com/funnyecho/code-push/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, err error) {
	Res(c, err, nil)
}

func Success(c *gin.Context, data interface{}) {
	Res(c, nil, data)
}

func Res(c *gin.Context, err error, data interface{}) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "S_OK",
			"data": data,
		})
		return
	}

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

	var reasonableErr *errors.Error

	isReasonableErr := stderr.As(err, reasonableErr)
	if !isReasonableErr {
		*reasonableErr = "FA_INTERNAL_ERROR"
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"code": reasonableErr.Error(),
	})
}
