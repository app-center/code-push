package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer/adapter/alioss"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/filer/interface/grpc"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
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
		Name:       fmt.Sprintf("Filer daemon, build at %s", BuildTime),
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
	serveCmdFS.StringVar(&(serveCmdOptions.ConfigFilePath), "config", "config/filer.d/serve.yml", "alternative config file path")
	serveCmdFS.BoolVar(&(serveCmdOptions.Debug), "debug", false, "run in debug mode")
	serveCmdFS.IntVar(&(serveCmdOptions.Port), "port", 7890, "port for grpc server listen to")
	serveCmdFS.StringVar(&(serveCmdOptions.BoltPath), "bolt-path", "storage/filer.d/db", "path of bolt file")
	serveCmdFS.StringVar(&(serveCmdOptions.AliOssEndpoint), "alioss-endpoint", "", "endpoint of ali-oss")
	serveCmdFS.StringVar(&(serveCmdOptions.AliOssAccessKeyId), "alioss-access-key-id", "", "access key id of ali-oss")
	serveCmdFS.StringVar(&(serveCmdOptions.AliOssAccessSecret), "alioss-access-secret", "", "access secret of ali-oss")

	serveCmd = &ffcli.Command{
		Name:       "serve",
		ShortUsage: "serve grpc server",
		ShortHelp:  fmt.Sprintf("%s serve [arguments]", executableName),
		FlagSet:    serveCmdFS,
		Options: []ff.Option{
			ff.WithEnvVarPrefix("FILER_D"),
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

	var g run.Group

	aliOssAdapter, aliOssAdapterErr := alioss.NewAliOssAdapter(
		serveCmdOptions.AliOssEndpoint,
		serveCmdOptions.AliOssAccessKeyId,
		serveCmdOptions.AliOssAccessSecret,
	)
	if aliOssAdapterErr != nil {
		return aliOssAdapterErr
	}

	domainAdapter := bolt.NewClient()
	domainAdapter.Path = serveCmdOptions.BoltPath
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		return domainAdapterOpenErr
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
		AliOssAdapter: aliOssAdapter,
	})

	grpcServer := interfacegrpc.NewFilerServer(endpoints)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serveCmdOptions.Port))
		if err != nil {
			return err
		}

		// Create gRPC server
		g.Add(func() error {
			baseServer := grpc.NewServer()
			pb.RegisterFileServer(baseServer, grpcServer)
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
