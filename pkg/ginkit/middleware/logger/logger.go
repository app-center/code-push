package ginkit_middleware_logger

import (
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/gin-gonic/gin"
)

func UseLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		clientIP := c.ClientIP()
		method := c.Request.Method

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()
		errorMessage := c.Errors.String()

		logger.Debug(
			"gin-request",
			"path", path,
			"clientIP", clientIP,
			"method", method,
			"statusCode", statusCode,
			"errorMessage", errorMessage,
		)
	}
}
