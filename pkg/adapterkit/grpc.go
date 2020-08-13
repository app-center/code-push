package adapterkit

import (
	"fmt"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	"github.com/funnyecho/code-push/pkg/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strings"
)

type GrpcAdaptOption interface {
	apply(*grpcAdapter)
}

func GrpcAdapter(options ...GrpcAdaptOption) Adaptable {
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
			grpcInterceptor.UnaryClientMetricInterceptor(a.logger),
			grpcInterceptor.UnaryClientErrorInterceptor(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpcInterceptor.StreamClientMetricInterceptor(a.logger),
			grpcInterceptor.StreamClientErrorInterceptor(),
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
