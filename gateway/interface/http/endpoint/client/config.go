package client

import "github.com/gin-gonic/gin"

func WithEnvId(branchId string, c *gin.Context) {
	c.Set("envId", branchId)
}

func UseEnvId(c *gin.Context) string {
	envId, _ := c.Get("envId")

	return envId.(string)
}
