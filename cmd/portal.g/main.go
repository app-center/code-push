package main

import (
	"context"
	"fmt"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	code_push "github.com/funnyecho/code-push/gateway/portal/adapter/code-push"
	"github.com/funnyecho/code-push/gateway/portal/adapter/filer"
	"github.com/funnyecho/code-push/gateway/portal/interface/http"
	"github.com/funnyecho/code-push/gateway/portal/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/funnyecho/code-push/pkg/svrkit"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		svrkit.WithCmdName("portal.g"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("client.g", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("PORTAL_G"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeHttpPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCodePushPort(&(serveCmdOptions.PortCodePushD)),
			svrkit.WithServeFilerPort(&(serveCmdOptions.PortFilerD)),
			svrkit.WithServeSessionAddr(&(serveCmdOptions.AddrSessionD)),
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

	codePushAdapter := code_push.New(
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "code-push.d")),
		func(options *code_push.Options) {
			options.ServerAddr = fmt.Sprintf("127.0.0.1:%d", serveCmdOptions.PortCodePushD)
		},
	)

	codePushConnErr := codePushAdapter.Conn()
	if codePushConnErr != nil {
		return codePushConnErr
	}
	defer codePushAdapter.Close()
	codePushAdapter.Debug("connected to code-push.d", "addr", codePushAdapter.ServerAddr)

	sessionAdapter := sessionAdapter.New(
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "session.d")),
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

	filerAdapter := filer.New(
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "filer.d")),
		func(options *filer.Options) {
			options.ServerAddr = fmt.Sprintf("127.0.0.1:%d", serveCmdOptions.PortFilerD)
		},
	)
	filerConnErr := filerAdapter.Conn()
	if filerConnErr != nil {
		return filerConnErr
	}
	defer filerAdapter.Close()
	filerAdapter.Debug("connected to filer.d", "addr", filerAdapter.ServerAddr)

	uc := usecase.NewUseCase(
		usecase.CtorConfig{
			CodePushAdapter: codePushAdapter,
			SessionAdapter:  sessionAdapter,
			FilerAdapter:    filerAdapter,
			Logger:          log.New(gokitLog.With(logger, "component", "usecase")),
		},
		func(options *usecase.Options) {

		},
	)

	server := http.New(
		uc,
		log.New(gokitLog.With(logger, "component", "interfaces", "interface", "http")),
		func(options *http.Options) {
			options.Port = serveCmdOptions.Port
		},
	)

	httpServeErr := server.ListenAndServe()
	return httpServeErr
}
