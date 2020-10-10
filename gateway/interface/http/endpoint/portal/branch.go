package portal

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"time"
)

type getBranchResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CreateAt int64  `json:"createAt"`
}

func GetBranch(c *gin.Context) {
	response, err := endpoint.UseUC(c).GetBranch(c.Request.Context(), UseBranchId(c))
	if err != nil {
		res.Error(c, err)
		return
	}

	res.Success(c, getBranchResponse{
		Id:       response.ID,
		Name:     response.Name,
		CreateAt: response.CreateTime.UnixNano() / int64(time.Millisecond),
	})
}
