package main

import (
	"context"
	"fmt"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/code-push/interface/grpc"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
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

	domainAdapter := bolt.NewClient()
	domainAdapter.Logger = log.New(gokitLog.With(logger, "component", "domain", "adapter", "bbolt"))
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		Logger:        log.New(gokitLog.With(logger, "component", "usecase")),
	})

	grpcServerLogger := log.New(gokitLog.With(logger, "component", "interfaces", "interface", "grpc"))
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
				)),
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpcInterceptor.StreamServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.StreamServerErrorInterceptor(),
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
