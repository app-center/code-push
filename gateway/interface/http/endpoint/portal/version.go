package portal

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

type releaseVersionRequest struct {
	EnvId            string `form:"envId" binding:"required"`
	AppVersion       string `form:"appVersion" binding:"required"`
	CompatAppVersion string `form:"compatAppVersion"`
	Changelog        string `form:"changelog" binding:"required"`
	PackageFileKey   string `form:"packageFileKey" binding:"required"`
	MustUpdate       bool   `form:"mustUpdate"`
}

type getVersionListRequest struct {
	EnvId string `uri:"envId" binding:"required"`
}

type versionResponse struct {
	EnvId            string `json:"envId"`
	AppVersion       string `json:"appVersion"`
	CompatAppVersion string `json:"compatAppVersion"`
	Changelog        string `json:"changelog"`
	PackageFileKey   string `json:"packageFileKey"`
	MustUpdate       bool   `json:"mustUpdate"`
	PublishAt        int64  `json:"publishAt"`
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

func GetVersionList(c *gin.Context) {
	var request getVersionListRequest

	if err := c.BindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	if request.EnvId == "" {
		res.Error(c, errors.WithMessage(gateway.ErrParamsInvalid, "envId is required"))
		return
	}

	response, err := endpoint.UseUC(c).GetVersionList(c.Request.Context(), []byte(request.EnvId))
	if err != nil {
		res.Error(c, err)
		return
	}

	list := make([]versionResponse, len(response))
	for i, v := range response {
		list[i] = versionResponse{
			EnvId:            v.EnvId,
			AppVersion:       v.AppVersion,
			CompatAppVersion: v.CompatAppVersion,
			Changelog:        v.Changelog,
			PackageFileKey:   v.PackageFileKey,
			MustUpdate:       v.MustUpdate,
			PublishAt:        v.CreateTime.UnixNano() / int64(time.Millisecond),
		}
	}

	res.Success(c, list)
	return
}
