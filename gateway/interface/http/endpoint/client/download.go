package client

import (
	"fmt"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type fileDownloadRequest struct {
	FileId string `uri:"fileId" binding:"required"`
}

func DownloadFile(c *gin.Context) {
	var request fileDownloadRequest

	if err := c.ShouldBindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	downloadUri, downloadUriErr := endpoint.UseUC(c).FileDownload(c.Request.Context(), []byte(request.FileId))
	if downloadUriErr != nil {
		res.Error(c, downloadUriErr)
		return
	}

	c.Writer.Header().Set("content-disposition", fmt.Sprint("attachment; filename=\"nextly.pkg\""))
	http.Redirect(c.Writer, c.Request, string(downloadUri), http.StatusSeeOther)
}
