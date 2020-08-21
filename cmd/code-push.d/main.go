package main

import (
	"context"
	"fmt"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/code-push/interface/grpc"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
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
		svrkit.WithCmdName("code-push.d"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("code-push.d", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("CODE_PUSH_D"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeGrpcPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCmdBoltPath("code-push.d", &(serveCmdOptions.BoltPath)),
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

	openTracer, openTracerCloser, openTracerErr := tracing.InitTracer("code-push.d")
	if openTracerErr == nil {
		opentracing.SetGlobalTracer(openTracer)
		defer openTracerCloser.Close()
	} else {
		logger.Infow("failed to init openTracer", "error", openTracerErr)
	}

	var g run.Group

	domainAdapter := bolt.NewClient()
	domainAdapter.Logger = zap_log.New(logger.With("component", "adapters", "adapter", "bbolt"))
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		Logger:        zap_log.New(logger.With("component", "usecase")),
	})

	grpcServerLogger := zap_log.New(logger.With("component", "interfaces", "interface", "grpc"))
	grpcServer := interfacegrpc.NewCodePushServer(
		endpoints,
		grpcServerLogger,
	)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", serveCmdOptions.Port))
		if err != nil {
			return err
		}

		// Create gRPC server
		g.Add(func() (err error) {
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
			pb.RegisterBranchServer(baseServer, grpcServer)
			pb.RegisterEnvServer(baseServer, grpcServer)
			pb.RegisterVersionServer(baseServer, grpcServer)
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
