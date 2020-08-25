package http

import (
	"github.com/funnyecho/code-push/gateway/client/interface/http/endpoints"
	"github.com/funnyecho/code-push/gateway/client/interface/http/middleware"
	"github.com/funnyecho/code-push/gateway/client/usecase"
	"github.com/funnyecho/code-push/pkg/ginMiddleware/opentracing"
	"github.com/funnyecho/code-push/pkg/log"
	prometheus_gin "github.com/funnyecho/code-push/pkg/promEndpoint/gin"
	"github.com/gin-gonic/gin"
	stdHttp "net/http"
)

func New(uc usecase.UseCase, logger log.Logger, fns ...func(*Options)) *server {
	ctorOptions := &Options{}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	svr := &server{
		uc:      uc,
		Logger:  logger,
		options: ctorOptions,
	}

	svr.initEndpoints()
	svr.initMiddleware()
	svr.initHttpHandler()

	return svr
}

type server struct {
	uc usecase.UseCase
	log.Logger
	options    *Options
	endpoints  *endpoints.Endpoints
	middleware *middleware.Middleware
	handler    stdHttp.Handler
}

func (s *server) ServeHTTP(writer stdHttp.ResponseWriter, request *stdHttp.Request) {
	s.handler.ServeHTTP(writer, request)
}

func (s *server) initEndpoints() {
	s.endpoints = endpoints.New(s.uc)
}

func (s *server) initMiddleware() {
	s.middleware = middleware.New(s.uc)
}

func (s *server) initHttpHandler() {
	r := gin.New()

	prometheus_gin.Init(r)

	r.GET("/client/download/pkg/:fileId", opentracing.StartTracing(), s.endpoints.DownloadFile)

	apiGroup := r.Group("/api/client", opentracing.StartTracing())
	apiGroup.POST("/auth/ddder", s.endpoints.Auth)
	apiGroup.GET("/v1/upgrade/:envId/:version", s.endpoints.VersionUpgradeQuery)

	s.handler = r
}

type Options struct {
}
