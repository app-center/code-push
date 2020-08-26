package ginkit_middleware_metrics

import (
	ginkit_metrics "github.com/funnyecho/code-push/pkg/ginkit/metrics"
	"github.com/gin-gonic/gin"
	"time"
)

func UseMetrics(metrics *ginkit_metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()

		metrics.ObserveHttpRequestDuration(durationToSeconds(time.Since(startTime)), method, path, status)
		if (status < 200 || status >= 400) || len(c.Errors) > 0 {
			metrics.IncHttpRequestFailed(method, path, status)
		} else {
			metrics.IncHttpRequestSucceed(method, path, status)
		}
	}
}

func durationToSeconds(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()/1000) / 1000 / 1000
}
