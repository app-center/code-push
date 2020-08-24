package endpoints

import (
	"github.com/funnyecho/code-push/pkg/ginResponse"
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

func (e *Endpoints) Auth(c *gin.Context) {
	var auth authRequest

	if err := c.Bind(&auth); err != nil {
		res.Error(c, err)
		return
	}

	authorizeErr := e.uc.Auth(c.Request.Context(), []byte(auth.BranchId), []byte(auth.Timestamp), []byte(auth.Nonce), []byte(auth.Sign))
	if authorizeErr != nil {
		res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := e.uc.SignToken(c.Request.Context(), []byte(auth.BranchId))
	if tokenErr != nil {
		res.Error(c, tokenErr)
		return
	}

	res.Success(c, authResponse{string(token)})
	return
}
