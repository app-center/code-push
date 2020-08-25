package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer/adapter/alioss"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/filer/interface/grpc"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
	"github.com/funnyecho/code-push/pkg/grpcInterceptor"
	http_kit "github.com/funnyecho/code-push/pkg/interfacekit/http"
	zap_log "github.com/funnyecho/code-push/pkg/log/zap"
	prometheus_http "github.com/funnyecho/code-push/pkg/promEndpoint/http"
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
		"filer.d",
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable(&(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeGrpcPort(&(serveCmdOptions.PortGrpc)),
			svrkit.WithServeHttpPort(&(serveCmdOptions.PortHttp)),
			svrkit.WithServeCmdBBoltPath(&(serveCmdOptions.BoltPath)),
			svrkit.WithServeCmdFlagSet(func(kit *svrkit.CmdKit, set *flag.FlagSet) {
				set.StringVar(&(serveCmdOptions.AliOssEndpoint), kit.FlagNameWithPrefix("alioss_endpoint"), "", "endpoint of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssBucket), kit.FlagNameWithPrefix("alioss_bucket"), "", "bucket of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssAccessKeyId), kit.FlagNameWithPrefix("alioss_access_key_id"), "", "access key id of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssAccessSecret), kit.FlagNameWithPrefix("alioss_access_secret"), "", "access secret of ali-oss")
			}),
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

	grpcServerLogger := zap_log.New(logger.With("component", "interfaces", "interface", "grpc"))
	grpcServer := interfacegrpc.NewFilerServer(
		endpoints,
		grpcServerLogger,
	)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveCmdOptions.PortGrpc))
		if err != nil {
			return err
		}

		// Create gRPC server
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
			pb.RegisterFileServer(baseServer, grpcServer)
			pb.RegisterUploadServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			grpcListener.Close()
		})
	}

	{
		g.Add(func() error {
			return http_kit.ListenAndServe(
				http_kit.WithServePort(serveCmdOptions.PortHttp),
				http_kit.WithDefaultServeMuxHandler(
					prometheus_http.Handle,
				),
			)
		}, func(err error) {

		})
	}

	err := g.Run()
	if err != nil {
		return err
	}

	return nil
}
