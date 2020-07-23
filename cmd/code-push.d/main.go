package main

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/code-push/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/code-push/interface/grpc"
	"github.com/funnyecho/code-push/daemon/code-push/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/code-push/usecase"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"net"
	"os"
)

var (
	Version   string
	BuildTime string
)

var (
	port int
)

var (
	cmd = &cobra.Command{
		Use:     "Code Push",
		Short:   "Daemon of Code-Push service",
		Long:    fmt.Sprintf("Daemon of Code-Push service. Build at %s", BuildTime),
		Version: Version,
		Run:     onCmdAction,
	}
)

func init() {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 7980, "grpc server port")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onCmdAction(cmd *cobra.Command, args []string) {
	var g run.Group

	domainAdapter := bolt.NewClient()
	domainAdapter.Path = ""
	domainAdapterOpenErr := domainAdapter.Open()
	if domainAdapterOpenErr != nil {
		os.Exit(1)
		return
	}
	defer domainAdapter.Close()

	endpoints := usecase.NewUseCase(usecase.CtorConfig{
		DomainAdapter: domainAdapter.DomainService(),
	})

	grpcServer := interfacegrpc.NewCodePushServer(endpoints)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			// FIXME: record log
			os.Exit(1)
		}

		// Create gRPC server
		g.Add(func() error {
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
		os.Exit(1)
		return
	}
}
