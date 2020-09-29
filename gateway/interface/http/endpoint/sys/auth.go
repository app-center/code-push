package sys

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	"github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username string `form:"userName" binding:"required"`
	Password string `form:"userPwd" binding:"required"`
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

	uc := endpoint.UseUC(c)

	authorizeErr := uc.AuthRootUser(c.Request.Context(), auth.Username, auth.Password)
	if authorizeErr != nil {
		ginkit_res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := uc.SignTokenForRootUser(c.Request.Context())
	if tokenErr != nil {
		ginkit_res.Error(c, tokenErr)
		return
	}

	ginkit_res.Success(c, authResponse{string(token)})
	return
}
