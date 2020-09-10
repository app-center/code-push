package main

import (
	"context"
	"github.com/funnyecho/code-push/daemon/filer/adapter/alioss"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/filer/interface/grpc"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
	grpckit_server "github.com/funnyecho/code-push/pkg/interfacekit/grpc/server"
	"github.com/funnyecho/code-push/pkg/interfacekit/http"
	zap_log "github.com/funnyecho/code-push/pkg/log/zap"
	prometheus_http "github.com/funnyecho/code-push/pkg/prom-endpoint/http"
	"github.com/funnyecho/code-push/pkg/svrkit"
	"github.com/funnyecho/code-push/pkg/tracing"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"syscall"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		"filer.d",
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable(),
			svrkit.WithServeCmdBindFlag(&serveCmdOptions),
			svrkit.WithServeCmdConfigValidation(&serveCmdOptions),
			svrkit.WithServeCmdPromFactorySetup(),
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

		logger = zapLogger.Sugar()

		defer logger.Sync()
	}

	openTracer, openTracerCloser, openTracerErr := tracing.InitTracer(
		"filer.d",
		zap_log.New(logger.With("component", "opentracing")),
	)
	if openTracerErr == nil {
		opentracing.SetGlobalTracer(openTracer)
		defer openTracerCloser.Close()
	} else {
		logger.Infow("failed to init openTracer", "error", openTracerErr)
	}

	var g run.Group

	aliOssAdapter, aliOssAdapterErr := alioss.NewAliOssAdapter(
		serveCmdOptions.AliOssEndpoint,
		serveCmdOptions.AliOssBucket,
		serveCmdOptions.AliOssAccessKeyId,
		serveCmdOptions.AliOssAccessSecret,
		zap_log.New(logger.With("component", "adapters", "adapter", "ali-oss")),
	)
	if aliOssAdapterErr != nil {
		return aliOssAdapterErr
	}

	domainAdapter := bolt.NewClient()
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapter.Logger = zap_log.New(logger.With("component", "domain", "adapter", "bbolt"))
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		AliOssAdapter: aliOssAdapter,
		Logger:        zap_log.New(logger.With("component", "usecase")),
	})

	{
		grpcServerLogger := zap_log.New(logger.With("component", "interfaces", "interface", "grpc"))
		grpcServer := interfacegrpc.NewFilerServer(
			endpoints,
			grpcServerLogger,
		)

		g.Add(grpckit_server.Actor(
			grpckit_server.WithServePort(serveCmdOptions.PortGrpc),
			grpckit_server.WithLogger(grpcServerLogger),
			grpckit_server.WithServeExecution(func(_ net.Listener, server *grpc.Server) error {
				pb.RegisterFileServer(server, grpcServer)
				pb.RegisterUploadServer(server, grpcServer)
				return nil
			}),
		))
	}

	{
		g.Add(httpkit.Actor(
			httpkit.WithServePort(serveCmdOptions.PortHttp),
			httpkit.WithDefaultServeMuxHandler(prometheus_http.Handle),
		))
	}

	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))

	err := g.Run()
	if err != nil {
		return err
	}

	return nil
}
