package main

import (
	"context"
	"fmt"
	sessiongrpc "github.com/funnyecho/code-push/daemon/session/interface/grpc"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/session/usecase"
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

	uc := usecase.New(
		log.New(gokitLog.With(logger, "component", "usecase")),
	)

	grpcServerLogger := log.New(gokitLog.With(logger, "component", "interfaces", "interface", "grpc"))
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
				)),
				grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
					grpcInterceptor.StreamServerMetricInterceptor(grpcServerLogger),
					grpcInterceptor.StreamServerErrorInterceptor(),
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
