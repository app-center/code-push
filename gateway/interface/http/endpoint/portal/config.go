package portal

import "github.com/gin-gonic/gin"

func WithBranchId(branchId string, c *gin.Context) {
	c.Set("branchId", branchId)
}

func UseBranchId(c *gin.Context) string {
	branchId, _ := c.Get("branchId")

	return branchId.(string)
}
