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
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/funnyecho/code-push/pkg/svrkit"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"net"
	"os"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		svrkit.WithCmdName("filer.d"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("filer.d", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("CODE_PUSH_D"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeGrpcPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCmdBoltPath("filer.d", &(serveCmdOptions.BoltPath)),
			svrkit.WithServeCmdFlagSet(func(set *flag.FlagSet) {
				set.StringVar(&(serveCmdOptions.AliOssEndpoint), "alioss-endpoint", "", "endpoint of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssBucket), "alioss-bucket", "", "bucket of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssAccessKeyId), "alioss-access-key-id", "", "access key id of ali-oss")
				set.StringVar(&(serveCmdOptions.AliOssAccessSecret), "alioss-access-secret", "", "access secret of ali-oss")
			}),
			svrkit.WithServeCmdConfigValidation(&serveCmdOptions),
			svrkit.WithServeCmdRun(onServe),
		),
	)
}

func onServe(ctx context.Context, args []string) error {
	var logger gokitLog.Logger
	{
		logger = gokitLog.NewLogfmtLogger(os.Stdout)
		logger = gokitLog.With(logger, "ts", gokitLog.DefaultTimestampUTC)
		logger = gokitLog.With(logger, "caller", gokitLog.DefaultCaller)

		if serveCmdOptions.Debug {
			logger = level.NewFilter(logger, level.AllowDebug())
		} else {
			logger = level.NewFilter(logger, level.AllowInfo())
		}
	}

	var g run.Group

	aliOssAdapter, aliOssAdapterErr := alioss.NewAliOssAdapter(
		serveCmdOptions.AliOssEndpoint,
		serveCmdOptions.AliOssBucket,
		serveCmdOptions.AliOssAccessKeyId,
		serveCmdOptions.AliOssAccessSecret,
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "ali-oss")),
	)
	if aliOssAdapterErr != nil {
		return aliOssAdapterErr
	}

	domainAdapter := bolt.NewClient()
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapter.Logger = log.New(gokitLog.With(logger, "component", "domain", "adapter", "bbolt"))
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		AliOssAdapter: aliOssAdapter,
		Logger:        log.New(gokitLog.With(logger, "component", "usecase")),
	})

	grpcServerLogger := log.New(gokitLog.With(logger, "component", "interfaces", "interface", "grpc"))
	grpcServer := interfacegrpc.NewFilerServer(
		endpoints,
		grpcServerLogger,
	)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveCmdOptions.Port))
		if err != nil {
			return err
		}

		// Create gRPC server
		g.Add(func() error {
			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
					grpcInterceptor.UnaryServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.UnaryServerErrorInterceptor(),
				)),
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpcInterceptor.StreamServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.StreamServerErrorInterceptor(),
				)),
			)
			pb.RegisterFileServer(baseServer, grpcServer)
			pb.RegisterUploadServer(baseServer, grpcServer)
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
