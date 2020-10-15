package main

import (
	"context"
	adapter "github.com/funnyecho/code-push/daemon/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/interface/http"
	"github.com/funnyecho/code-push/gateway/usecase"
	http_kit "github.com/funnyecho/code-push/pkg/interfacekit/http"
	zap_log "github.com/funnyecho/code-push/pkg/log/zap"
	"github.com/funnyecho/code-push/pkg/svrkit"
	"github.com/funnyecho/code-push/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		"code-push.gateway",
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
		"sys.g",
		zap_log.New(logger.With("component", "opentracing")),
	)
	if openTracerErr == nil {
		opentracing.SetGlobalTracer(openTracer)
		defer openTracerCloser.Close()
	} else {
		logger.Infow("failed to init openTracer", "error", openTracerErr)
	}

	daemonAdapter := adapter.New(
		zap_log.New(logger.With("component", "adapters", "adapter", "code-push.d")),
		func(options *adapter.Options) {
			options.ServerAddr = serveCmdOptions.AddrDaemon
		},
	)

	daemonConnErr := daemonAdapter.Conn()
	if daemonConnErr != nil {
		return daemonConnErr
	}
	defer daemonAdapter.Close()
	daemonAdapter.Debug("connected to code-push.daemon", "addr", serveCmdOptions.AddrDaemon)

	uc := usecase.New(func(config *usecase.CtorConfig) {
		config.Logger = zap_log.New(logger.With("component", "usecase"))
		config.DaemonAdapter = daemonAdapter
		config.Options.RootUserName = serveCmdOptions.RootUserName
		config.Options.RootUserPwd = serveCmdOptions.RootUserPwd
	})

	return http_kit.ListenAndServe(
		http_kit.WithServePort(serveCmdOptions.Port),
		http_kit.WithServeHandler(http.New(func(options *http.Options) {
			options.Debug = serveCmdOptions.Debug
			options.UseCase = uc
			options.Logger = zap_log.New(logger.With("component", "interfaces", "interface", "http"))
			options.AppCenterPath = serveCmdOptions.AppCenterPath
		})),
	)
}
