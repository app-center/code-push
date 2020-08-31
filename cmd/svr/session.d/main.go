package main

import (
	"context"
	sessiongrpc "github.com/funnyecho/code-push/daemon/session/interface/grpc"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/session/usecase"
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
		"session.d",
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
		defer logger.Sync()

		logger = zapLogger.Sugar()
	}

	openTracer, openTracerCloser, openTracerErr := tracing.InitTracer(
		"session.d",
		zap_log.New(logger.With("component", "opentracing")),
	)
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

	{
		grpcServerLogger := zap_log.New(logger.With("component", "interfaces", "interface", "grpc"))
		grpcServer := sessiongrpc.New(
			uc,
			grpcServerLogger,
		)

		g.Add(grpckit_server.Actor(
			grpckit_server.WithServePort(serveCmdOptions.PortGrpc),
			grpckit_server.WithLogger(grpcServerLogger),
			grpckit_server.WithServeExecution(func(_ net.Listener, server *grpc.Server) error {
				pb.RegisterAccessTokenServer(server, grpcServer)
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
