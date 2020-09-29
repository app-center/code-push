package sys

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type createBranchRequest struct {
	BranchName string `form:"branchName" binding:"required"`
}

type createBranchResponse struct {
	BranchId       string `form:"branchId" binding:"required"`
	BranchName     string `form:"branchName" binding:"required"`
	BranchEncToken string `form:"branchEncToken" binding:"required"`
}

func CreateBranch(c *gin.Context) {
	var input createBranchRequest
	if err := c.ShouldBind(&input); err != nil {
		res.Error(c, err)
		return
	}

	branch, createErr := endpoint.UseUC(c).CreateBranch(c.Request.Context(), []byte(input.BranchName))
	if createErr != nil {
		res.Error(c, createErr)
		return
	}

	res.Success(c, createBranchResponse{
		BranchId:       branch.ID,
		BranchName:     branch.Name,
		BranchEncToken: branch.EncToken,
	})
}
