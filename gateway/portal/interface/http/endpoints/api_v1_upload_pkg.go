package endpoints

import (
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type uploadPkgRequest struct {
	Pkg *multipart.FileHeader `form:"pkg" binding:"required"`
}

type uploadPkgResponse struct {
	FileKey string `json:"file_key"`
}

func (e *Endpoints) UploadPkg(c *gin.Context) {
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

	fileKey, uploadErr := e.uc.UploadPkg(c.Request.Context(), uploadStream)
	if uploadErr != nil {
		res.Error(c, uploadErr)
		return
	}

	res.Success(c, uploadPkgResponse{FileKey: string(fileKey)})
	return
}
