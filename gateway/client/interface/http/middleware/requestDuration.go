package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

func (m *Middleware) RequestDuration(c *gin.Context) {
	//startTime := time.Now()

	c.Next()

	//m.uc.RequestDuration(c.Request.URL.Path, c.Writer.Status() == http.StatusOK, durationToSeconds(time.Since(startTime)))
}

func durationToSeconds(duration time.Duration) float64 {
	return float64(duration.Nanoseconds()/1000) / 1000 / 1000
}
