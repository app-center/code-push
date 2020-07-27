package endpoints

import (
	"github.com/funnyecho/code-push/gateway/client/interface/http/middleware"
	res "github.com/funnyecho/code-push/pkg/gin_res"
	"github.com/gin-gonic/gin"
	"net/http"
)

type versionUpgradeQueryRequest struct {
	AppVersion string `uri:"version" binding:"required"`
}

type versionUpgradeQueryResponse struct {
	AppVersion          []byte `json:"app_version"`
	LatestAppVersion    []byte `json:"latest_app_version"`
	CanUpdateAppVersion []byte `json:"can_update_app_version"`
	MustUpdate          bool   `json:"must_update"`
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

	result, queryErr := e.uc.VersionStrictCompatQuery(envId, []byte(request.AppVersion))
	if queryErr != nil {
		res.Error(c, queryErr)
		return
	}

	res.Success(c, &versionUpgradeQueryResponse{
		AppVersion:          result.AppVersion,
		LatestAppVersion:    result.LatestAppVersion,
		CanUpdateAppVersion: result.CanUpdateAppVersion,
		MustUpdate:          result.MustUpdate,
	})
}
