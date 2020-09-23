package portal

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type createEnvRequest struct {
	EnvName string `form:"env_name" binding:"required"`
}

type createEnvResponse struct {
	EnvId       string `json:"env_id"`
	EnvEncToken string `json:"env_enc_token"`
}

func CreateEnv(c *gin.Context) {
	var request createEnvRequest

	if err := c.Bind(&request); err != nil {
		res.Error(c, err)
		return
	}

	response, err := endpoint.UseUC(c).CreateEnv(c.Request.Context(), []byte(UseBranchId(c)), []byte(request.EnvName))
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
