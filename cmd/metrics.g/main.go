package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc"
	"github.com/funnyecho/code-push/gateway/metric/interface/grpc/pb"
	"github.com/funnyecho/code-push/gateway/metric/interface/http"
	"github.com/funnyecho/code-push/gateway/metric/usecase"
	"github.com/funnyecho/code-push/pkg/svrkit"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	stdgrpc "google.golang.org/grpc"
	"net"
	"os"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		svrkit.WithCmdName("metrics.d"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("metrics.d", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("METRICS_G"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeHttpPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCmdFlagSet(func(set *flag.FlagSet) {
				set.IntVar(&(serveCmdOptions.PortMetric), "port-metric", 0, "port for grpc server listen to")
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

	httpServer := http.New(
		nil,
		func(options *http.Options) {
			options.Port = serveCmdOptions.Port
		},
	)

	{
		g.Add(func() error {
			return httpServer.ListenAndServe()
		}, func(err error) {

		})
	}

	uc := usecase.New()
	grpcServer := grpc.New(uc)
	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveCmdOptions.PortMetric))
		if err != nil {
			return err
		}

		g.Add(func() error {
			baseServer := stdgrpc.NewServer()
			pb.RegisterRequestDurationServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			grpcListener.Close()
		})
	}

	return g.Run()
}
