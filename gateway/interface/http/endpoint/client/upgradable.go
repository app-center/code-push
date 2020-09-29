package client

import (
	"github.com/funnyecho/code-push/gateway"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	res "github.com/funnyecho/code-push/pkg/ginkit/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type versionUpgradeQueryRequest struct {
	EnvId      string `uri:"envId" binding:"required"`
	AppVersion string `uri:"version" binding:"required"`
}

type versionUpgradeQueryResponse struct {
	AppVersion          string                     `json:"appVersion"`
	LatestAppVersion    string                     `json:"latestAppVersion"`
	CanUpdateAppVersion string                     `json:"canUpdateAppVersion"`
	MustUpdate          bool                       `json:"mustUpdate"`
	PackageInfo         *versionUpgradePackageInfo `json:"packageInfo"`
}

type versionUpgradePackageInfo struct {
	PackageBlob string `json:"packageBlob"`
	MD5         string `json:"md5"`
	Size        int64  `json:"size"`
}

func VersionUpgradeQuery(c *gin.Context) {
	var request versionUpgradeQueryRequest

	if err := c.ShouldBindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	envId := UseEnvId(c)
	if envId != request.EnvId {
		res.Error(c, errors.Wrapf(gateway.ErrInvalidEnv, "envId in url differ from token"))
		return
	}

	result, queryErr := endpoint.UseUC(c).VersionStrictCompatQuery(c.Request.Context(), []byte(envId), []byte(request.AppVersion))
	if queryErr != nil {
		res.Error(c, queryErr)
		return
	}

	var queryResCode string

	queryData := &versionUpgradeQueryResponse{
		MustUpdate:  result.MustUpdate,
		PackageInfo: nil,
	}

	if result.AppVersion != nil {
		queryData.AppVersion = string(result.AppVersion)
	}

	if result.CanUpdateAppVersion != nil {
		queryData.CanUpdateAppVersion = string(result.CanUpdateAppVersion)
	}

	if result.LatestAppVersion == nil {
		queryResCode = "S_LATEST"
	} else if result.CanUpdateAppVersion != nil {
		queryResCode = "S_UPGRADABLE"

		upgradeSource, upgradeSourceErr := endpoint.UseUC(c).VersionPkgSource(c.Request.Context(), envId, string(result.CanUpdateAppVersion))
		if upgradeSourceErr != nil {
			res.Error(c, errors.Wrapf(upgradeSourceErr, "failed to get upgradable version source"))
			return
		}

		queryData.PackageInfo = &versionUpgradePackageInfo{
			PackageBlob: upgradeSource.Key,
			MD5:         upgradeSource.FileMD5,
			Size:        upgradeSource.FileSize,
		}
	} else {
		queryResCode = "S_LIMITED_LATEST"
	}

	res.SuccessWithCode(c, queryResCode, queryData)
}
