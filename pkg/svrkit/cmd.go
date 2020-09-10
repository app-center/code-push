package svrkit

import (
	"context"
	"flag"
	"fmt"
	"github.com/funnyecho/code-push/pkg/flagkit"
	"github.com/funnyecho/code-push/pkg/metrics"
	promfactory "github.com/funnyecho/code-push/pkg/metrics/prometheus"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
	"os"
	"path/filepath"
	"strings"
)

var (
	Version       string
	BuildTime     string
	BuildPlatform string
)

type CmdOption func(kit *CmdKit)
type ServeCmdOption func(kit *CmdKit, cmd *ffcli.Command)

func RunCmd(cmdName string, options ...CmdOption) {
	flagPrefix := strings.Replace(cmdName, ".", "_", -1)
	flagPrefix = strings.Replace(flagPrefix, "-", "_", -1)
	runner := &CmdKit{
		executableName: filepath.Base(os.Args[0]),
		name:           cmdName,
		flagPrefix:     flagPrefix,
	}

	for _, fn := range options {
		fn(runner)
	}

	var cmd, versionCmd, serveCmd *ffcli.Command

	versionCmd = &ffcli.Command{
		Name:      "version",
		ShortHelp: "Version of service",
		Options:   []ff.Option{ff.WithIgnoreUndefined(true)},
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(fmt.Sprintf("%s %s %s", runner.name, Version, BuildPlatform))
			return nil
		},
	}

	serveCmd = &ffcli.Command{
		Name:        "serve",
		ShortHelp:   fmt.Sprintf("%s serve [arguments]", runner.executableName),
		FlagSet:     flag.NewFlagSet("serve", flag.ExitOnError),
		Options:     []ff.Option{ff.WithIgnoreUndefined(true)},
		Subcommands: nil,
		Exec: func(ctx context.Context, args []string) error {
			return nil
		},
	}

	serveCmd.Options = append(serveCmd.Options, ff.WithEnvVarNoPrefix())

	for _, fn := range runner.serveCmdOptions {
		fn(runner, serveCmd)
	}

	cmd = &ffcli.Command{
		Name:       fmt.Sprintf("%s, build at %s", runner.name, BuildTime),
		ShortUsage: fmt.Sprintf("%s <command> [arguments]", runner.executableName),
		UsageFunc: func(c *ffcli.Command) string {
			return fmt.Sprintf("%s\n\n%s", c.Name, ffcli.DefaultUsageFunc(c))
		},
		FlagSet: nil,
		Options: []ff.Option{ff.WithIgnoreUndefined(true)},
		Subcommands: []*ffcli.Command{
			versionCmd,
			serveCmd,
		},
		Exec: func(ctx context.Context, args []string) error {
			return serveCmd.ParseAndRun(ctx, args)
		},
	}

	if err := cmd.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		os.Exit(1)
	}
}

func WithCmdName(name string) CmdOption {
	return func(kit *CmdKit) {
		kit.name = name
	}
}

func WithServeCmd(options ...ServeCmdOption) CmdOption {
	return func(kit *CmdKit) {
		kit.serveCmdOptions = options
	}
}

func WithServeCmdFlagSet(fn func(kit *CmdKit, set *flag.FlagSet)) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		fn(kit, cmd.FlagSet)
	}
}

func WithServeCmdBindFlag(flags interface{}) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		flagkit.MustBind(flags, cmd.FlagSet)
	}
}

func WithServeCmdConfigurable() ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.Options = append(cmd.Options, ff.WithConfigFileFlag("config"), ff.WithAllowMissingConfigFile(true), ff.WithConfigFileParser(ffyaml.Parser))
	}
}

func WithServeCmdDebuggable(debug *bool) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.BoolVar(debug, kit.FlagNameWithPrefix("debug"), false, "run in debug mode")
	}
}

func WithServeHttpPort(port *int) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.IntVar(port, kit.FlagNameWithPrefix("port_http"), 0, "port for http server listen to")
	}
}

func WithServeGrpcPort(port *int) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.IntVar(port, kit.FlagNameWithPrefix("port_grpc"), 0, "port for grpc server listen to")
	}
}

func WithServeCodePushAddr(addr *string) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr_code_push_d", "", "address of code-push.d")
	}
}

func WithServeFilerAddr(addr *string) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr_filer_d", "", "address of filer.d")
	}
}

func WithServeSessionAddr(addr *string) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(addr, "addr_session_d", "", "address of session.d")
	}
}

func WithServeCmdBBoltPath(path *string) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
		cmd.FlagSet.StringVar(path, kit.FlagNameWithPrefix("bbolt_path"), fmt.Sprintf("storage/%s/bbolt.db", kit.name), "path of bbolt storage file")
	}
}

func WithServeCmdRun(fn func(ctx context.Context, args []string) error) ServeCmdOption {
	return func(kit *CmdKit, cmd *ffcli.Command) {
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

func WithServeCmdPromFactorySetup() ServeCmdOption {
	return WithServeCmdRun(func(ctx context.Context, args []string) error {
		metrics.SetDefaultFactory(promfactory.New())
		return nil
	})
}

type CmdKit struct {
	executableName  string
	name            string
	flagPrefix      string
	serveCmdOptions []ServeCmdOption
}

func (k *CmdKit) FlagNameWithPrefix(flagName string) string {
	return fmt.Sprintf("%s.%s", k.flagPrefix, flagName)
}
