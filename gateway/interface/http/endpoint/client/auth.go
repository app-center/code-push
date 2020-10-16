package client

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	ginkit_res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	EnvId     string `form:"env_id" binding:"required"`
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

	authorizeErr := endpoint.UseUC(c).AuthEnv(c.Request.Context(), []byte(auth.EnvId), []byte(auth.Timestamp), []byte(auth.Nonce), []byte(auth.Sign))
	if authorizeErr != nil {
		ginkit_res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := endpoint.UseUC(c).SignTokenForEnv(c.Request.Context(), []byte(auth.EnvId))
	if tokenErr != nil {
		ginkit_res.Error(c, tokenErr)
		return
	}

	ginkit_res.Success(c, authResponse{string(token)})
	return
}
