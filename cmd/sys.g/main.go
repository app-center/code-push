package main

import (
	"context"
	"flag"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc_adapter"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc_adapter"
	"github.com/funnyecho/code-push/gateway/sys/interface/http"
	"github.com/funnyecho/code-push/gateway/sys/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	"github.com/funnyecho/code-push/pkg/svrkit"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

var serveCmdOptions serveConfig

func main() {
	svrkit.RunCmd(
		svrkit.WithCmdName("sys.g"),
		svrkit.WithServeCmd(
			svrkit.WithServeCmdConfigurable("sys.g", &(serveCmdOptions.ConfigFilePath)),
			svrkit.WithServeCmdEnvPrefix("SYS_G"),
			svrkit.WithServeCmdDebuggable(&(serveCmdOptions.Debug)),
			svrkit.WithServeHttpPort(&(serveCmdOptions.Port)),
			svrkit.WithServeCodePushAddr(&(serveCmdOptions.AddrCodePushD)),
			svrkit.WithServeSessionAddr(&(serveCmdOptions.AddrSessionD)),
			svrkit.WithServeCmdFlagSet(func(set *flag.FlagSet) {
				set.StringVar(&(serveCmdOptions.RootUserName), "root-user-name", "", "root user name")
				set.StringVar(&(serveCmdOptions.RootUserPwd), "root-user-pwd", "", "root user password")
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

	codePushAdapter := codePushAdapter.New(
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "code-push.d")),
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

	useCase := usecase.NewUseCase(
		usecase.CtorConfig{
			CodePushAdapter: codePushAdapter,
			SessionAdapter:  sessionAdapter,
			Logger:          log.New(gokitLog.With(logger, "component", "usecase")),
		},
		func(options *usecase.Options) {
			options.RootUserName = serveCmdOptions.RootUserName
			options.RootUserPwd = serveCmdOptions.RootUserPwd
		},
	)

	server := http.New(
		useCase,
		log.New(gokitLog.With(logger, "component", "interfaces", "interface", "http")),
		func(options *http.Options) {
			options.Port = serveCmdOptions.Port
		},
	)

	httpServeErr := server.ListenAndServe()
	if httpServeErr != nil {
		return httpServeErr
	}

	return nil
}
