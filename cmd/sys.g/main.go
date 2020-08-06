package main

import (
	"context"
	"flag"
	"fmt"
	code_push "github.com/funnyecho/code-push/gateway/sys/adapter/code-push"
	"github.com/funnyecho/code-push/gateway/sys/adapter/session"
	"github.com/funnyecho/code-push/gateway/sys/interface/http"
	"github.com/funnyecho/code-push/gateway/sys/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
	"os"
	"path/filepath"
)

var (
	Version       string
	BuildTime     string
	BuildPlatform string
)

var (
	executableName string
)

var cmd, versionCmd, serveCmd *ffcli.Command
var serveCmdOptions serveConfig

func init() {
	executableName = filepath.Base(os.Args[0])
}

func main() {
	initServeCmd()

	versionCmd = &ffcli.Command{
		Name:      "version",
		ShortHelp: "Version of service",
		Exec:      onVersion,
	}

	cmd = &ffcli.Command{
		Name:       fmt.Sprintf("Sys gateway, build at %s", BuildTime),
		ShortUsage: fmt.Sprintf("%s <command> [arguments]", executableName),
		UsageFunc: func(c *ffcli.Command) string {
			return fmt.Sprintf("%s\n\n%s", c.Name, ffcli.DefaultUsageFunc(c))
		},
		FlagSet: nil,
		Options: nil,
		Subcommands: []*ffcli.Command{
			versionCmd,
			serveCmd,
		},
		Exec: onRoot,
	}

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Printf("FF failed to parse and run: %s", err.Error())
		os.Exit(1)
	}
}

func initServeCmd() {
	serveCmdFS := flag.NewFlagSet("serve", flag.ExitOnError)
	serveCmdFS.StringVar(&(serveCmdOptions.ConfigFilePath), "config", "config/sys.g/serve.yml", "alternative config file path")
	serveCmdFS.BoolVar(&(serveCmdOptions.Debug), "debug", false, "run in debug mode")
	serveCmdFS.IntVar(&(serveCmdOptions.Port), "port", 0, "port for grpc server listen to")
	serveCmdFS.IntVar(&(serveCmdOptions.PortCodePushD), "port-code-push", 0, "port of code-push.d")
	serveCmdFS.IntVar(&(serveCmdOptions.PortSessionD), "port-session", 0, "port of session.d")

	serveCmd = &ffcli.Command{
		Name:       "serve",
		ShortUsage: "serve grpc server",
		ShortHelp:  fmt.Sprintf("%s serve [arguments]", executableName),
		FlagSet:    serveCmdFS,
		Options: []ff.Option{
			ff.WithEnvVarPrefix("SYS_G"),
			ff.WithEnvVarSplit("_"),
			ff.WithConfigFileFlag("config"),
			ff.WithAllowMissingConfigFile(true),
			ff.WithConfigFileParser(ffyaml.Parser),
		},
		Subcommands: nil,
		Exec:        onServe,
	}
}

func onRoot(ctx context.Context, args []string) error {
	return serveCmd.ParseAndRun(ctx, args)
}

func onVersion(ctx context.Context, args []string) error {
	fmt.Println(fmt.Sprintf("Version of Xiner-Web %s %s%", Version, BuildPlatform))
	return nil
}

func onServe(ctx context.Context, args []string) error {
	if configErr := serveCmdOptions.validate(); configErr != nil {
		return configErr
	}

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

	sessionAdapter := session.New(
		log.New(gokitLog.With(logger, "component", "adapters", "adapter", "session.d")),
		func(options *session.Options) {
			options.ServerAddr = fmt.Sprintf("127.0.0.1:%d", serveCmdOptions.PortSessionD)
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
