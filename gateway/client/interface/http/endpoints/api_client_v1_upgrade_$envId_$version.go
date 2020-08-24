package endpoints

import (
	"github.com/funnyecho/code-push/gateway/client"
	"github.com/funnyecho/code-push/gateway/client/interface/http/middleware"
	res "github.com/funnyecho/code-push/pkg/ginResponse"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type versionUpgradeQueryRequest struct {
	EnvId      string `uri:"envId" binding:"required"`
	AppVersion string `uri:"version" binding:"required"`
}

type versionUpgradeQueryResponse struct {
	AppVersion          string                     `json:"app_version"`
	LatestAppVersion    string                     `json:"latest_app_version"`
	CanUpdateAppVersion string                     `json:"can_update_app_version"`
	MustUpdate          bool                       `json:"must_update"`
	PackageInfo         *versionUpgradePackageInfo `json:"package_info"`
}

type versionUpgradePackageInfo struct {
	PackageBlob string `json:"package_blob"`
	MD5         string `json:"md5"`
	Size        int64  `json:"size"`
}

func (e *Endpoints) VersionUpgradeQuery(c *gin.Context) {
	var request versionUpgradeQueryRequest

	if err := c.ShouldBindUri(&request); err != nil {
		res.Error(c, err)
		return
	}

	envId, authErr := middleware.AuthorizedWithReturns(e.uc, c)
	if authErr != nil {
		res.ErrorWithStatusCode(c, http.StatusUnauthorized, authErr)
		return
	}

	if string(envId) != request.EnvId {
		res.Error(c, errors.Wrapf(client.ErrInvalidEnv, "envId in url differ from token"))
		return
	}

	result, queryErr := e.uc.VersionStrictCompatQuery(c.Request.Context(), envId, []byte(request.AppVersion))
	if queryErr != nil {
		res.Error(c, queryErr)
		return
	}

	var queryResCode string

	queryData := &versionUpgradeQueryResponse{
		AppVersion:          result.AppVersion,
		LatestAppVersion:    result.LatestAppVersion,
		CanUpdateAppVersion: result.CanUpdateAppVersion,
		MustUpdate:          result.MustUpdate,
		PackageInfo:         nil,
	}

	if result.LatestAppVersion == "" {
		queryResCode = "S_LATEST"
	} else if result.CanUpdateAppVersion != "" {
		queryResCode = "S_UPGRADABLE"

		upgradeSource, upgradeSourceErr := e.uc.VersionPkgSource(c.Request.Context(), string(envId), result.CanUpdateAppVersion)
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
