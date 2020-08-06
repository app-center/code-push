package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/code-push/interface/grpc"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
	"github.com/funnyecho/code-push/pkg/log"
	gokitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/peterbourgon/ff/v3/ffyaml"
	"google.golang.org/grpc"
	"net"
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
		Name:       fmt.Sprintf("Code-Push daemon, build at %s", BuildTime),
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
	serveCmdFS.StringVar(&(serveCmdOptions.ConfigFilePath), "config", "config/code-push.d/serve.yml", "alternative config file path")
	serveCmdFS.BoolVar(&(serveCmdOptions.Debug), "debug", false, "run in debug mode")
	serveCmdFS.IntVar(&(serveCmdOptions.Port), "port", 0, "port for grpc server listen to")
	serveCmdFS.StringVar(&(serveCmdOptions.BoltPath), "bolt-path", "storage/code-push.d/db", "path of bolt storage file")

	serveCmd = &ffcli.Command{
		Name:       "serve",
		ShortUsage: "serve grpc server",
		ShortHelp:  fmt.Sprintf("%s serve [arguments]", executableName),
		FlagSet:    serveCmdFS,
		Options: []ff.Option{
			ff.WithEnvVarPrefix("CODE_PUSH_D"),
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

	var g run.Group

	domainAdapter := bolt.NewClient()
	domainAdapter.Logger = log.New(gokitLog.With(logger, "component", "domain", "adapter", "bbolt"))
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		Logger:        log.New(gokitLog.With(logger, "component", "usecase")),
	})

	grpcServer := interfacegrpc.NewCodePushServer(
		endpoints,
		log.New(gokitLog.With(logger, "component", "interfaces", "interface", "grpc")),
	)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", serveCmdOptions.Port))
		if err != nil {
			return err
		}

		// Create gRPC server
		g.Add(func() (err error) {
			baseServer := grpc.NewServer()
			pb.RegisterBranchServer(baseServer, grpcServer)
			pb.RegisterEnvServer(baseServer, grpcServer)
			pb.RegisterVersionServer(baseServer, grpcServer)
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
