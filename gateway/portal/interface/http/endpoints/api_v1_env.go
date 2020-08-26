package endpoints

import (
	"github.com/funnyecho/code-push/gateway/portal/interface/http/middleware"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createEnvRequest struct {
	EnvName string `form:"env_name" binding:"required"`
}

type createEnvResponse struct {
	EnvId       string `json:"env_id"`
	EnvEncToken string `json:"env_enc_token"`
}

func (e *Endpoints) CreateEnv(c *gin.Context) {
	var request createEnvRequest

	if err := c.Bind(&request); err != nil {
		res.Error(c, err)
		return
	}

	branchId, authErr := middleware.AuthorizedWithReturns(e.uc, c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	response, err := e.uc.CreateEnv(c.Request.Context(), branchId, []byte(request.EnvName))
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, createEnvResponse{
		EnvId:       response.ID,
		EnvEncToken: response.EncToken,
	})
	return
}
