package portal

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	"github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRequest struct {
	BranchId  string `form:"branch_id" binding:"required"`
	Timestamp string `form:"timestamp" binding:"required"`
	Nonce     string `form:"nonce" binding:"required"`
	Sign      string `form:"sign" binding:"required"`
}

type jwtAuthRequest struct {
	Token string `form:"token" binding:"required"`
}

type refreshAuthRequest struct {
	Token string `uri:"token" binding:"required"`
}

type authResponse struct {
	Token string `json:"token"`
}

type refreshAuthResponse struct {
	BranchId string `json:"branchId"`
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

func AuthWithJwt(c *gin.Context) {
	var auth jwtAuthRequest

	if err := c.Bind(&auth); err != nil {
		ginkit_res.Error(c, err)
		return
	}

	branchId, authorizeErr := endpoint.UseUC(c).AuthBranchWithJWT(c.Request.Context(), auth.Token)
	if authorizeErr != nil {
		ginkit_res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := endpoint.UseUC(c).SignTokenForBranch(c.Request.Context(), branchId)
	if tokenErr != nil {
		ginkit_res.Error(c, tokenErr)
		return
	}

	ginkit_res.Success(c, authResponse{string(token)})
	return
}

func RefreshAuthorization(c *gin.Context) {
	var request refreshAuthRequest

	if err := c.ShouldBindUri(&request); err != nil {
		ginkit_res.Error(c, err)
		return
	}

	branchId, verifyErr := endpoint.UseUC(c).VerifyTokenForBranch(c.Request.Context(), []byte(request.Token))
	if verifyErr != nil {
		ginkit_res.ErrorWithStatusCode(c, http.StatusUnauthorized, gateway.ErrUnauthorized)
		return
	}

	ginkit_res.Success(c, refreshAuthResponse{BranchId: string(branchId)})
	return
}
