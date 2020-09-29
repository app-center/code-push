package portal

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
)

type releaseVersionRequest struct {
	EnvId            string `form:"envId" binding:"required"`
	AppVersion       string `form:"appVersion" binding:"required"`
	CompatAppVersion string `form:"compatAppVersion"`
	Changelog        string `form:"changelog" binding:"required"`
	PackageFileKey   string `form:"packageFileKey" binding:"required"`
	MustUpdate       bool   `form:"mustUpdate"`
}

func ReleaseVersion(c *gin.Context) {
	var request releaseVersionRequest

	if err := c.Bind(&request); err != nil {
		res.Error(c, err)
		return
	}

	releaseErr := endpoint.UseUC(c).ReleaseVersion(c.Request.Context(), &gateway.VersionReleaseParams{
		EnvId:            []byte(request.EnvId),
		AppVersion:       []byte(request.AppVersion),
		CompatAppVersion: []byte(request.CompatAppVersion),
		Changelog:        []byte(request.Changelog),
		PackageFileKey:   []byte(request.PackageFileKey),
		MustUpdate:       request.MustUpdate,
	})

	if releaseErr != nil {
		res.Error(c, releaseErr)
		return
	}

	res.Success(c, nil)
}
