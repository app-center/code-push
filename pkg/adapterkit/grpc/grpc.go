package adapterkit_grpc

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/adapterkit"
	adapterkit_grpc_interceptor_error "github.com/funnyecho/code-push/pkg/adapterkit/grpc/interceptor/error"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strings"
)

type GrpcAdaptOption interface {
	apply(*grpcAdapter)
}

func GrpcAdapter(options ...GrpcAdaptOption) adapterkit.Adaptable {
	adapter := &grpcAdapter{}

	for _, fn := range options {
		fn.apply(adapter)
	}

	return adapter
}

func WithGrpcAdaptTarget(target string) GrpcAdaptOption {
	return newFuncGrpcAdaptOption(func(adapter *grpcAdapter) {
		if strings.HasPrefix(target, ":") {
			target = fmt.Sprintf("127.0.0.1%s", target)
		}

		adapter.target = target
	})
}

func WithGrpcAdaptName(name string) GrpcAdaptOption {
	return newFuncGrpcAdaptOption(func(adapter *grpcAdapter) {
		adapter.name = name
	})
}

func WithGrpcAdaptLogger(logger log.Logger) GrpcAdaptOption {
	return newFuncGrpcAdaptOption(func(adapter *grpcAdapter) {
		adapter.logger = logger
	})
}

func WithGrpcAdaptConnected(fn func(*grpc.ClientConn)) GrpcAdaptOption {
	return newFuncGrpcAdaptOption(func(adapter *grpcAdapter) {
		adapter.onConnected = fn
	})
}

func WithGrpcAdaptClosed(fn func(*grpc.ClientConn)) GrpcAdaptOption {
	return newFuncGrpcAdaptOption(func(adapter *grpcAdapter) {
		adapter.onClosed = fn
	})
}

type grpcAdapter struct {
	name   string
	target string

	logger log.Logger
	conn   *grpc.ClientConn

	onConnected func(*grpc.ClientConn)
	onClosed    func(*grpc.ClientConn)
}

func (a *grpcAdapter) Conn() error {
	conn, err := grpc.Dial(
		a.target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			adapterkit_grpc_interceptor_error.WithUnaryInterceptor(),
			//adapterkit_grpc_interceptor_opentracing.WithUnaryInterceptor(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			adapterkit_grpc_interceptor_error.WithStreamInterceptor(),
			//adapterkit_grpc_interceptor_opentracing.WithStreamInterceptor(),
		)),
	)
	if err != nil {
		return errors.Wrapf(err, "Dail to grpc server: %s failed", a.target)
	}

	a.conn = conn
	if a.onConnected != nil {
		a.onConnected(conn)
	}

	return nil
}

func (a *grpcAdapter) Close() error {
	if a.onClosed != nil {
		defer a.onClosed(a.conn)
	}

	if a.conn != nil {
		return a.conn.Close()
	}

	return nil
}

type grpcAdaptOption func(*grpcAdapter)

type funcGrpcAdaptOption struct {
	f grpcAdaptOption
}

func (fdo *funcGrpcAdaptOption) apply(do *grpcAdapter) {
	fdo.f(do)
}

func newFuncGrpcAdaptOption(f grpcAdaptOption) *funcGrpcAdaptOption {
	return &funcGrpcAdaptOption{
		f: f,
	}
}
