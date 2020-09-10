package httpkit

import (
	"fmt"
	"net/http"
)

type ServeOptions func(options *serverOptions)
type ServeMuxOptions func(*http.ServeMux)

func Actor(withOptions ...ServeOptions) (execute func() error, interrupt func(error)) {
	var server *http.Server
	var serverErr error

	execute = func() error {
		serverErr, server = createServer(withOptions...)

		if serverErr != nil {
			return serverErr
		}

		serverErr = server.ListenAndServe()
		return serverErr
	}

	interrupt = func(_ error) {
		if server != nil {
			server.Close()
		}
	}

	return
}

func ListenAndServe(withOptions ...ServeOptions) error {
	err, server := createServer(withOptions...)

	if err != nil {
		return err
	}

	return server.ListenAndServe()
}

func createServer(withOptions ...ServeOptions) (error, *http.Server) {
	option := &serverOptions{
		port:    0,
		handler: nil,
	}

	for _, fn := range withOptions {
		fn(option)
	}

	addr := fmt.Sprintf(":%d", option.port)
	server := &http.Server{
		Addr:    addr,
		Handler: option.handler,
	}

	return nil, server
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
	port    int
	handler http.Handler
}
