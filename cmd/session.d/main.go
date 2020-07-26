package main

import (
	"fmt"
	sessiongrpc "github.com/funnyecho/code-push/daemon/session/interface/grpc"
	"github.com/funnyecho/code-push/daemon/session/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/session/usecase"
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
		Use:     "Filer",
		Short:   "Daemon of Session service",
		Long:    fmt.Sprintf("Daemon of Session service. Build at %s", BuildTime),
		Version: Version,
		Run:     onCmdAction,
	}
)

func init() {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 7982, "grpc server port")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onCmdAction(cmd *cobra.Command, args []string) {
	var g run.Group

	uc := usecase.New()

	grpcServer := sessiongrpc.New(uc)
	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			// FIXME: record log
			os.Exit(1)
		}

		g.Add(func() error {
			baseServer := grpc.NewServer()
			pb.RegisterAccessTokenServer(baseServer, grpcServer)
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
