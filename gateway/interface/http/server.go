package http

import (
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint/client"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint/portal"
	"github.com/funnyecho/code-push/gateway/interface/http/endpoint/sys"
	"github.com/funnyecho/code-push/gateway/usecase"
	ginkit_server "github.com/funnyecho/code-push/pkg/ginkit/server"
	"github.com/funnyecho/code-push/pkg/log"
	stdHttp "net/http"
)

func New(configFn func(*Options)) stdHttp.Handler {
	config := &Options{}

	configFn(config)

	r := ginkit_server.New(
		ginkit_server.WithDebugMode(config.Debug),
		ginkit_server.WithLogger(config.Logger),
	)

	r.Use(
		endpoint.WithLogger(config.Logger),
		endpoint.WithUseCase(config.UseCase),
	)



	gSys := r.Group("/sys")
	gPortal := r.Group("/portal")
	gClient := r.Group("/client")

	gSysApi := gSys.Group("/api/sys")
	gSysApi.POST("/auth", sys.Auth)

	gSysApiV1 := gSysApi.Group("/v1")
	gSysApiV1.POST("/branch", sys.MidAuthorized, sys.CreateBranch)

	gPortalApi := gPortal.Group("/api/portal")
	gPortalApi.POST("/auth", portal.Auth)

	gPortalApiV1 := gPortalApi.Group("/v1")
	gPortalApiV1.POST("/env", portal.MidAuthorized, portal.CreateEnv)
	gPortalApiV1.POST("/version", portal.MidAuthorized, portal.ReleaseVersion)
	gPortalApiV1.POST("/upload/pkg", portal.MidAuthorized, portal.UploadPkg)

	gClient.GET("/download/pkg/:fileId", client.DownloadFile)

	gClientApi := gSys.Group("/api/client")
	gClientApi.POST("/auth", client.Auth)
	gClientApiV1 := gClientApi.Group("/v1")
	gClientApiV1.GET("/upgrade/:envId/:version", client.MidAuthorized, client.VersionUpgradeQuery)

	return r
}

type Options struct {
	usecase.UseCase
	log.Logger
	Debug bool
}