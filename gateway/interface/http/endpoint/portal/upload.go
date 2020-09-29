package portal

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type uploadPkgRequest struct {
	Pkg *multipart.FileHeader `form:"pkg" binding:"required"`
}

type uploadPkgResponse struct {
	FileKey string `json:"fileKey"`
}

func UploadPkg(c *gin.Context) {
	var request uploadPkgRequest

	if err := c.Bind(&request); err != nil {
		res.Error(c, err)
		return
	}

	uploadStream, uploadStreamErr := request.Pkg.Open()
	if uploadStreamErr != nil {
		res.Error(c, uploadStreamErr)
		return
	}
	defer uploadStream.Close()

	fileKey, uploadErr := endpoint.UseUC(c).UploadPkg(c.Request.Context(), uploadStream)
	if uploadErr != nil {
		res.Error(c, uploadErr)
		return
	}

	res.Success(c, uploadPkgResponse{FileKey: string(fileKey)})
	return
}
