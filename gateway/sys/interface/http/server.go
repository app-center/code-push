package http

import (
	"github.com/funnyecho/code-push/gateway/sys/interface/http/endpoints"
	"github.com/funnyecho/code-push/gateway/sys/interface/http/middleware"
	"github.com/funnyecho/code-push/gateway/sys/usecase"
	ginkit_server "github.com/funnyecho/code-push/pkg/ginkit/server"
	"github.com/funnyecho/code-push/pkg/log"
	stdHttp "net/http"
)

func New(config *CtorConfig, fns ...func(*Options)) *server {
	ctorOptions := &Options{}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	svr := &server{
		uc:      config.UseCase,
		Logger:  config.Logger,
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
	r := ginkit_server.New(
		ginkit_server.WithDebugMode(s.options.Debug),
		ginkit_server.WithLogger(s.Logger),
	)

	apiGroup := r.Group("/api")
	apiGroup.POST("/auth", s.endpoints.Auth)

	apiGroup.POST("/v1/branch", s.middleware.Authorized, s.endpoints.CreateBranch)

	s.handler = r
}

type CtorConfig struct {
	usecase.UseCase
	log.Logger
}

type Options struct {
	Debug bool
}
