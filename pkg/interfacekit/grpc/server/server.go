package grpckit_server

import (
	"fmt"
	grpckit_interceptor_error "github.com/funnyecho/code-push/pkg/interfacekit/grpc/interceptor/error"
	grpckit_interceptor_logger "github.com/funnyecho/code-push/pkg/interfacekit/grpc/interceptor/logger"
	grpckit_interceptor_metrics "github.com/funnyecho/code-push/pkg/interfacekit/grpc/interceptor/metrics"
	grpckit_interceptor_opentracing "github.com/funnyecho/code-push/pkg/interfacekit/grpc/interceptor/opentracing"
	grpckit_metrics "github.com/funnyecho/code-push/pkg/interfacekit/grpc/metrics"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"net"
)

func Actor(withOptions ...ServerOptions) (execute func() error, interrupt func(error)) {
	var listener net.Listener
	var server *grpc.Server
	var actErr error

	execute = func() error {
		actErr, listener, server = listenAndServe(withOptions...)

		return actErr
	}

	interrupt = func(_ error) {
		if server != nil {
			server.Stop()
		}

		if listener != nil {
			listener.Close()
		}
	}

	return
}

func ListenAndServe(withOptions ...ServerOptions) (error, net.Listener, *grpc.Server) {
	return listenAndServe(withOptions...)
}

type ServerOptions func(options *serverOptions)

func WithServePort(port int) ServerOptions {
	return func(options *serverOptions) {
		options.port = port
	}
}

func WithServeExecution(fn func(net.Listener, *grpc.Server) error) ServerOptions {
	return func(options *serverOptions) {
		options.onExecute = fn
	}
}

func WithDisableMetrics() ServerOptions {
	return func(options *serverOptions) {
		options.disableMetric = true
	}
}

func WithDisableOpentracing() ServerOptions {
	return func(options *serverOptions) {
		options.disableOpentracing = true
	}
}

func WithLogger(logger log.Logger) ServerOptions {
	return func(options *serverOptions) {
		options.logger = logger
	}
}

func listenAndServe(withOptions ...ServerOptions) (error, net.Listener, *grpc.Server) {
	options := &serverOptions{}
	for _, fn := range withOptions {
		fn(options)
	}

	listener, listenErr := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", options.port))
	if listenErr != nil {
		return listenErr, nil, nil
	}

	var unaryInterceptor []grpc.UnaryServerInterceptor
	var streamInterceptor []grpc.StreamServerInterceptor

	{
		unaryInterceptor = append(unaryInterceptor, grpckit_interceptor_error.UseUnaryInterceptor())
		streamInterceptor = append(streamInterceptor, grpckit_interceptor_error.UseStreamInterceptor())
	}

	{
		unaryInterceptor = append(unaryInterceptor, grpckit_interceptor_logger.UseUnaryInterceptor(options.logger))
		streamInterceptor = append(streamInterceptor, grpckit_interceptor_logger.UseStreamInterceptor(options.logger))
	}

	if !options.disableMetric {
		metrics := grpckit_metrics.NewMetrics()
		unaryInterceptor = append(unaryInterceptor, grpckit_interceptor_metrics.UseUnaryInterceptor(metrics))
		streamInterceptor = append(streamInterceptor, grpckit_interceptor_metrics.UseStreamInterceptor(metrics))
	}

	if !options.disableOpentracing {
		unaryInterceptor = append(unaryInterceptor, grpckit_interceptor_opentracing.UseUnaryInterceptor())
		streamInterceptor = append(streamInterceptor, grpckit_interceptor_opentracing.UseStreamInterceptor())
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptor...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptor...)),
	)

	executeErr := options.onExecute(listener, server)
	if executeErr != nil {
		return executeErr, listener, nil
	}

	serveErr := server.Serve(listener)
	if serveErr != nil {
		return serveErr, listener, nil
	}

	return nil, listener, server
}

type serverOptions struct {
	port               int
	logger             log.Logger
	disableOpentracing bool
	disableMetric      bool
	onExecute          func(net.Listener, *grpc.Server) error
}
