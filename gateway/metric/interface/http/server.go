package http

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func New(config *Config, fns ...func(*Options)) interface {
	http.Handler
	ListenAndServe() error
} {
	svrOptions := &Options{}

	for _, fn := range fns {
		fn(svrOptions)
	}

	svr := &server{
		options: svrOptions,
	}

	svr.initHttpHandler()

	return svr
}

type server struct {
	options *Options
}

func (s *server) ListenAndServe() error {
	addr := fmt.Sprintf(":%d", s.options.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        s,
		MaxHeaderBytes: 1 << 20,
	}

	return server.ListenAndServe()
}

func (s *server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	http.DefaultServeMux.ServeHTTP(writer, request)
}

func (s *server) initHttpHandler() {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
}

type Config struct {
}

type Options struct {
	Port int
}
