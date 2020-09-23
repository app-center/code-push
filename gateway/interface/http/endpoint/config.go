package endpoint

import (
	"fmt"
	"github.com/funnyecho/code-push/gateway/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/gin-gonic/gin"
)

func WithLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
	}
}

func UseLogger(c *gin.Context) log.Logger {
	logger, existed := c.Get("logger")

	if !existed {
		return log.PanicLogger
	}

	return logger.(log.Logger)
}

func WithUseCase(uc usecase.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("usecase", uc)
	}
}

func UseUC(c *gin.Context) usecase.UseCase {
	uc, existed := c.Get("usecase")
	if !existed {
		panic(fmt.Errorf("usecase not existed in gin context"))
	}

	return uc.(usecase.UseCase)
}