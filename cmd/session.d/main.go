package main

import (
	"context"
	"fmt"
	sessiongrpc "github.com/funnyecho/code-push/daemon/session/interface/grpc"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/session/usecase"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	zap_log "github.com/funnyecho/code-push/pkg/log/zap"
	"github.com/funnyecho/code-push/pkg/svrkit"
	"github.com/funnyecho/code-push/pkg/tracing"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		svrkit.WithCmdName("session.d"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("session.d", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("SESSION_D"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeGrpcPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCmdConfigValidation(&serveCmdOptions),
			svrkit.WithServeCmdRun(onServe),
		),
	)
}

func onServe(ctx context.Context, args []string) error {
	var logger *zap.SugaredLogger
	{
		var zapLogger *zap.Logger
		if serveCmdOptions.Debug {
			zapLogger, _ = zap.NewDevelopment()
		} else {
			zapLogger, _ = zap.NewProduction()
		}
		defer logger.Sync()

		logger = zapLogger.Sugar()
	}

	openTracer, openTracerCloser, openTracerErr := tracing.InitTracer("session.d")
	if openTracerErr == nil {
		opentracing.SetGlobalTracer(openTracer)
		defer openTracerCloser.Close()
	} else {
		logger.Infow("failed to init openTracer", "error", openTracerErr)
	}

	var g run.Group

	uc := usecase.New(
		zap_log.New(logger.With("component", "usecase")),
	)

	grpcServerLogger := zap_log.New(logger.With("component", "interfaces", "interface", "grpc"))
	grpcServer := sessiongrpc.New(
		uc,
		grpcServerLogger,
	)
	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveCmdOptions.Port))
		if err != nil {
			return err
		}

		g.Add(func() error {
			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
					grpcInterceptor.UnaryServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.UnaryServerErrorInterceptor(),
					grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
				)),
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpcInterceptor.StreamServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.StreamServerErrorInterceptor(),
					grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
				)),
			)
			pb.RegisterAccessTokenServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			grpcListener.Close()
		})
	}

	err := g.Run()
	if err != nil {
		return err
	}

	return nil
}
