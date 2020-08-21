package endpoints

import (
	res "github.com/funnyecho/code-push/pkg/gin_res"
	"github.com/gin-gonic/gin"
)

type createBranchRequest struct {
	BranchName string `form:"branch_name" binding:"required"`
}

type createBranchResponse struct {
	BranchId       string `form:"branch_id" binding:"required"`
	BranchName     string `form:"branch_name" binding:"required"`
	BranchEncToken string `form:"branch_enc_token" binding:"required"`
}

func (e *Endpoints) CreateBranch(c *gin.Context) {
	var input createBranchRequest
	if err := c.ShouldBind(&input); err != nil {
		res.Error(c, err)
		return
	}

	branch, createErr := e.uc.CreateBranch(c.Request.Context(), []byte(input.BranchName))
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
