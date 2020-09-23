package portal

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	"github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	BranchId  string `form:"branch_id" binding:"required"`
	Timestamp string `form:"timestamp" binding:"required"`
	Nonce     string `form:"nonce" binding:"required"`
	Sign      string `form:"sign" binding:"required"`
}

type authResponse struct {
	Token string `json:"token"`
}

func Auth(c *gin.Context) {
	var auth authRequest

	if err := c.Bind(&auth); err != nil {
		ginkit_res.Error(c, err)
		return
	}

	authorizeErr := endpoint.UseUC(c).AuthBranch(c.Request.Context(), []byte(auth.BranchId), []byte(auth.Timestamp), []byte(auth.Nonce), []byte(auth.Sign))
	if authorizeErr != nil {
		ginkit_res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := endpoint.UseUC(c).SignTokenForBranch(c.Request.Context(), []byte(auth.BranchId))
	if tokenErr != nil {
		ginkit_res.Error(c, tokenErr)
		return
	}

	ginkit_res.Success(c, authResponse{string(token)})
	return
}
