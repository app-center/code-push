package endpoints

import (
	"fmt"
	"github.com/funnyecho/code-push/gateway/client/interface/http/middleware"
	res "github.com/funnyecho/code-push/pkg/gin_res"
	"github.com/gin-gonic/gin"
	"net/http"
)

type versionDownloadRequest struct {
	AppVersion string `uri:"version" binding:"required"`
	FileName   string `uri:"filename"`
}

func (e *Endpoints) DownloadVersionPkg(c *gin.Context) {
	var request versionDownloadRequest

	if err := c.ShouldBindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	envId, authErr := middleware.AuthorizedWithReturns(e.uc, c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	downloadUri, downloadUriErr := e.uc.VersionDownloadPkg(envId, []byte(request.AppVersion))
	if downloadUriErr != nil {
		res.Error(c, downloadUriErr)
		return
	}

	attachmentName := request.FileName
	if len(attachmentName) == 0 {
		attachmentName = fmt.Sprintf("code-push__%s.pkg", request.AppVersion)
	}

	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachmentName))
	http.Redirect(c.Writer, c.Request, string(downloadUri), http.StatusSeeOther)
}
