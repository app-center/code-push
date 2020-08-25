package http_kit

import (
	"fmt"
	"net/http"
)

type ServeOptions func(options *serverOptions)
type ServeMuxOptions func(*http.ServeMux)

func ListenAndServe(options ...ServeOptions) error {
	option := &serverOptions{
		port:    0,
		handler: nil,
	}

	for _, fn := range options {
		fn(option)
	}

	addr := fmt.Sprintf(":%d", option.port)
	server := &http.Server{
		Addr:              addr,
		Handler:           option.handler,
	}

	return server.ListenAndServe()
}

func WithServePort(port int) ServeOptions {
	return func(options *serverOptions) {
		options.port = port
	}
}

func WithServeHandler(handler http.Handler) ServeOptions {
	return func(options *serverOptions) {
		options.handler = handler
	}
}

func WithDefaultServeMuxHandler(serveMuxOptions ...ServeMuxOptions) ServeOptions {
	return func(options *serverOptions) {
		for _, fn := range serveMuxOptions {
			fn(http.DefaultServeMux)
		}

		options.handler = http.DefaultServeMux
	}
}

func WithServeMuxPatternHandler(pattern string, handler http.Handler) ServeMuxOptions {
	return func(mux *http.ServeMux) {
		mux.Handle(pattern, handler)
	}
}

type serverOptions struct {
	port int
	handler http.Handler
}
