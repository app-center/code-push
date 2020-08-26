package endpoints

import (
	"github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username string `form:"user_name" binding:"required"`
	Password string `form:"user_pwd" binding:"required"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (e *Endpoints) Auth(c *gin.Context) {
	var auth authRequest

	if err := c.Bind(&auth); err != nil {
		ginkit_res.Error(c, err)
		return
	}

	authorizeErr := e.uc.Auth(c.Request.Context(), auth.Username, auth.Password)
	if authorizeErr != nil {
		ginkit_res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := e.uc.SignToken(c.Request.Context())
	if tokenErr != nil {
		ginkit_res.Error(c, tokenErr)
		return
	}

	ginkit_res.Success(c, authResponse{string(token)})
	return
}
