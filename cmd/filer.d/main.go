package main

import (
	"fmt"
	"github.com/funnyecho/code-push/daemon/filer/domain/bolt"
	interfacegrpc "github.com/funnyecho/code-push/daemon/filer/interface/grpc"
	"github.com/funnyecho/code-push/daemon/filer/interface/grpc/pb"
	"github.com/funnyecho/code-push/daemon/filer/usecase"
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
		Short:   "Filer daemon of Code-Push service",
		Long:    fmt.Sprintf("Filer daemon of Code-Push service. Build at %s", BuildTime),
		Version: Version,
		Run:     onCmdAction,
	}
)

func init() {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 7981, "grpc server port")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onCmdAction(cmd *cobra.Command, args []string) {
	var g run.Group

	var (
		client = bolt.NewClient()
	)

	client.Path = ""
	domainOpenErr := client.Open()

	if domainOpenErr != nil {
		// FIXME: record log
		os.Exit(1)
	}

	var (

		fileUseCase = usecase.NewFileUseCase(usecase.FileUseCaseConfig{
			SchemeService: client.SchemeService(),
			FileService:   client.FileService(),
		})

		schemeUseCase = usecase.NewSchemeUseCase(usecase.SchemeUseCaseConfig{
			SchemeService: client.SchemeService(),
		})

		uploadUseCase = usecase.NewUploadUseCase(usecase.UploadUseCaseConfig{
			SchemeService: client.SchemeService(),
		})
	)

	{
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			// FIXME: record log
			os.Exit(1)
		}

		// Create gRPC server
		g.Add(func() error {
			baseServer := grpc.NewServer()
			grpcServer := interfacegrpc.NewFilerServer(interfacegrpc.FilerServerConfig{
				FileUseCase:   fileUseCase,
				SchemeUseCase: schemeUseCase,
				UploadUseCase: uploadUseCase,
			})
			pb.RegisterFileServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			grpcListener.Close()
		})
	}

	g.Run()
}
