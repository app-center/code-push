package endpoints

import (
	"github.com/funnyecho/code-push/pkg/gin_res"
	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username []byte `form:"user_name" binding:"required"`
	Password []byte `form:"user_pwd" binding:"required"`
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

	authorizeErr := e.uc.Auth(auth.Username, auth.Password)
	if authorizeErr != nil {
		res.Error(c, authorizeErr)
		return
	}

	token, tokenErr := e.uc.SignToken()
	if tokenErr != nil {
		res.Error(c, tokenErr)
		return
	}

	res.Success(c, authResponse{string(token)})
	return
}
