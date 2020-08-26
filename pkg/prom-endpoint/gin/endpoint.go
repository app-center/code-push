package prometheus_gin

import (
	prom_endpoint "github.com/funnyecho/code-push/pkg/prom-endpoint"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(r *gin.Engine) {
	r.GET(prom_endpoint.Endpoint, func() gin.HandlerFunc {
		promHttpHandler := promhttp.Handler()
		return func(context *gin.Context) {
			promHttpHandler.ServeHTTP(context.Writer, context.Request)
		}
	}())
}
