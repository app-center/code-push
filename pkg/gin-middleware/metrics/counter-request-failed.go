package metrics

import "github.com/gin-gonic/gin"

func counterRequestFailed(c *gin.Context) {
	c.Next()
}
