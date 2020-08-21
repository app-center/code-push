package svrkit

import (
	"context"
	"flag"
	"fmt"
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

type CmdOption func(kit *cmdKit)
type ServeCmdOption func(cmd *ffcli.Command)

func RunCmd(options ...CmdOption) {
	runner := &cmdKit{
		executableName: filepath.Base(os.Args[0]),
		name:           "RunCmd Kit",
	}

	for _, fn := range options {
		fn(runner)
	}

	var cmd, versionCmd, serveCmd *ffcli.Command

	versionCmd = &ffcli.Command{
		Name:      "version",
		ShortHelp: "Version of service",
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(fmt.Sprintf("%s %s %s", runner.name, Version, BuildPlatform))
			return nil
		},
	}

	serveCmd = &ffcli.Command{
		Name:        "serve",
		ShortHelp:   fmt.Sprintf("%s serve [arguments]", runner.executableName),
		FlagSet:     flag.NewFlagSet("serve", flag.ExitOnError),
		Options:     []ff.Option{},
		Subcommands: nil,
		Exec: func(ctx context.Context, args []string) error {
			return nil
		},
	}

	for _, fn := range runner.serveCmdOptions {
		fn(serveCmd)
	}

	cmd = &ffcli.Command{
		Name:       fmt.Sprintf("%s, build at %s", runner.name, BuildTime),
		ShortUsage: fmt.Sprintf("%s <command> [arguments]", runner.executableName),
		UsageFunc: func(c *ffcli.Command) string {
			return fmt.Sprintf("%s\n\n%s", c.Name, ffcli.DefaultUsageFunc(c))
		},
		FlagSet: nil,
		Options: nil,
		Subcommands: []*ffcli.Command{
			versionCmd,
			serveCmd,
		},
		Exec: func(ctx context.Context, args []string) error {
			return serveCmd.ParseAndRun(ctx, args)
		},
	}

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Printf("FF failed to parse and run: %s", err.Error())
		os.Exit(1)
	}
}

func WithCmdName(name string) CmdOption {
	return func(kit *cmdKit) {
		kit.name = name
	}
}

func WithServeCmd(options ...ServeCmdOption) CmdOption {
	return func(kit *cmdKit) {
		kit.serveCmdOptions = options
	}
}

func WithServeCmdFlagSet(fn func(set *flag.FlagSet)) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		fn(cmd.FlagSet)
	}
}

func WithServeCmdEnvPrefix(prefix string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.Options = append(cmd.Options, ff.WithEnvVarPrefix(prefix), ff.WithEnvVarSplit("_"))
	}
}

func WithServeCmdConfigurable(name string, path *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(path, "config", fmt.Sprintf("config/%s/serve.yml", name), "alternative config file path")
		cmd.Options = append(cmd.Options, ff.WithConfigFileFlag("config"), ff.WithAllowMissingConfigFile(true), ff.WithConfigFileParser(ffyaml.Parser))
	}
}

func WithServeCmdDebuggable(debug *bool) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.BoolVar(debug, "debug", false, "run in debug mode")
	}
}

func WithServeHttpPort(port *int) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.IntVar(port, "port", 0, "port for http server listen to")
	}
}

func WithServeGrpcPort(port *int) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.IntVar(port, "port", 0, "port for grpc server listen to")
	}
}

func WithServeCodePushAddr(addr *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr-code-push", "", "address of code-push.d")
	}
}

func WithServeFilerAddr(addr *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr-filer", "", "address of filer.d")
	}
}

func WithServeSessionAddr(addr *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr-session", "", "address of session.d")
	}
}

func WithServeMetricAddress(addr *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr-metric", "", "address of metric.g")
	}
}

func WithServeTracingReporterAddress(addr *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr-tracing-reporter", "", "address of opentracing reporter")
	}
}

func WithServeCmdBoltPath(name string, path *string) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(path, "bolt-path", fmt.Sprintf("storage/%s/db", name), "path of bolt storage file")
	}
}

func WithServeCmdRun(fn func(ctx context.Context, args []string) error) ServeCmdOption {
	return func(cmd *ffcli.Command) {
		onServe := cmd.Exec

		cmd.Exec = func(ctx context.Context, args []string) error {
			var err error
			err = onServe(ctx, args)
			if err != nil {
				return err
			}

			return fn(ctx, args)
		}
	}
}

func WithServeCmdConfigValidation(validator interface {
	Validate() error
}) ServeCmdOption {
	return WithServeCmdRun(func(ctx context.Context, args []string) error {
		if configErr := validator.Validate(); configErr != nil {
			return configErr
		}

		return nil
	})
}

type cmdKit struct {
	executableName  string
	name            string
	serveCmdOptions []ServeCmdOption
}
