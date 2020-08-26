package main

import (
	"context"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc_adapter"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/sys/interface/http"
	"github.com/funnyecho/code-push/gateway/sys/usecase"
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
		"sys.g",
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
		"sys.g",
		zap_log.New(logger.With("component", "opentracing")),
	)
	if openTracerErr == nil {
		opentracing.SetGlobalTracer(openTracer)
		defer openTracerCloser.Close()
	} else {
		logger.Infow("failed to init openTracer", "error", openTracerErr)
	}

	codePushAdapter := codePushAdapter.New(
		zap_log.New(logger.With("component", "adapters", "adapter", "code-push.d")),
		func(options *codePushAdapter.Options) {
			options.ServerAddr = serveCmdOptions.AddrCodePushD
		},
	)

	codePushConnErr := codePushAdapter.Conn()
	if codePushConnErr != nil {
		return codePushConnErr
	}
	defer codePushAdapter.Close()
	codePushAdapter.Debug("connected to code-push.d", "addr", codePushAdapter.ServerAddr)

	sessionAdapter := sessionAdapter.New(
		zap_log.New(logger.With("component", "adapters", "adapter", "session.d")),
		func(options *sessionAdapter.Options) {
			options.ServerAddr = serveCmdOptions.AddrSessionD
		},
	)
	sessionConnErr := sessionAdapter.Conn()
	if sessionConnErr != nil {
		return sessionConnErr
	}
	defer sessionAdapter.Close()
	sessionAdapter.Debug("connected to session.d", "addr", sessionAdapter.ServerAddr)

	uc := usecase.NewUseCase(
		usecase.CtorConfig{
			CodePushAdapter: codePushAdapter,
			SessionAdapter:  sessionAdapter,
			Logger:          zap_log.New(logger.With("component", "usecase")),
		},
		func(options *usecase.Options) {
			options.RootUserName = serveCmdOptions.RootUserName
			options.RootUserPwd = serveCmdOptions.RootUserPwd
		},
	)

	return http_kit.ListenAndServe(
		http_kit.WithServePort(serveCmdOptions.Port),
		http_kit.WithServeHandler(http.New(
			&http.CtorConfig{
				UseCase: uc,
				Logger:  zap_log.New(logger.With("component", "interfaces", "interface", "http")),
			},
			func(options *http.Options) {
				options.Debug = serveCmdOptions.Debug
			},
		)),
	)
}
