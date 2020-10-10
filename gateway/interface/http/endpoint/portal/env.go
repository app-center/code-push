package portal

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type createEnvRequest struct {
	EnvName string `form:"envName" binding:"required"`
}

type getEnvRequest struct {
	EnvId string `uri:"envId" binding:"required"`
}

type envResponse struct {
	Id       string `json:"id"`
	BranchId string `json:"branchId"`
	Name     string `json:"name"`
	EncToken string `json:"encToken"`
	CreateAt int64  `json:"createAt"`
}

type createEnvResponse = envResponse

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
		Id:       response.ID,
		BranchId: response.BranchId,
		Name:     response.Name,
		EncToken: response.EncToken,
		CreateAt: response.CreateTime.UnixNano() / int64(time.Millisecond),
	})
	return
}

func GetEnvs(c *gin.Context) {
	response, err := endpoint.UseUC(c).GetEnvsWithBranchId(c.Request.Context(), UseBranchId(c))
	if err != nil {
		res.Error(c, err)
		return
	}

	list := make([]envResponse, len(response))
	for i, v := range response {
		list[i] = envResponse{
			Id:       v.ID,
			BranchId: v.BranchId,
			Name:     v.Name,
			EncToken: v.EncToken,
			CreateAt: v.CreateTime.UnixNano() / int64(time.Millisecond),
		}
	}

	res.Success(c, list)
	return
}

func GetEnv(c *gin.Context) {
	var request getEnvRequest

	if err := c.ShouldBindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	if request.EnvId == "" {
		res.Error(c, errors.WithMessage(gateway.ErrParamsInvalid, "envId is required"))
		return
	}

	env, err := endpoint.UseUC(c).GetEnv(c.Request.Context(), []byte(request.EnvId))
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, createEnvResponse{
		Id:       env.ID,
		BranchId: env.BranchId,
		Name:     env.Name,
		EncToken: env.EncToken,
		CreateAt: env.CreateTime.UnixNano() / int64(time.Millisecond),
	})
	return
}
