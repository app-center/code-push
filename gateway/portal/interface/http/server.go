package http

import (
	"fmt"
	"github.com/funnyecho/code-push/gateway/portal/interface/http/endpoints"
	"github.com/funnyecho/code-push/gateway/portal/interface/http/middleware"
	"github.com/funnyecho/code-push/gateway/portal/usecase"
	"github.com/gin-gonic/gin"
	stdHttp "net/http"
)

func New(uc usecase.UseCase, fns ...func(*Options)) *server {
	ctorOptions := &Options{}

	for _, fn := range fns {
		fn(ctorOptions)
	}

	svr := &server{
		uc:      uc,
		options: ctorOptions,
	}

	svr.initEndpoints()
	svr.initMiddleware()
	svr.initHttpHandler()

	return svr
}

type server struct {
	uc         usecase.UseCase
	options    *Options
	endpoints  *endpoints.Endpoints
	middleware *middleware.Middleware
	handler    stdHttp.Handler
}

func (s *server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.options.Port)
	server := &stdHttp.Server{
		Addr:           addr,
		Handler:        s,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
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

	apiGroup := r.Group("/api")
	apiGroup.POST("/auth", s.endpoints.Auth)
	apiGroup.POST("/env", s.middleware.Authorized, s.endpoints.CreateEnv)
	apiGroup.POST("/version", s.middleware.Authorized, s.endpoints.ReleaseVersion)
	apiGroup.POST("/upload/pkg", s.middleware.Authorized, s.endpoints.UploadPkg)

	s.handler = r
}

type Options struct {
	Port int
}
