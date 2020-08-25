package prometheus_gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(r *gin.Engine) {
	r.GET("/metrics", func() gin.HandlerFunc {
		promHttpHandler := promhttp.Handler()
		return func(context *gin.Context) {
			promHttpHandler.ServeHTTP(context.Writer, context.Request)
		}
	}())
}
