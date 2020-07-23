package main

import (
	"fmt"
	code_push "github.com/funnyecho/code-push/gateway/sys/adapter/code-push"
	"github.com/funnyecho/code-push/gateway/sys/interface/http"
	"github.com/funnyecho/code-push/gateway/sys/usecase"
	"github.com/spf13/cobra"
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
		Use:     "System Gateway",
		Short:   "Gateway of System service",
		Long:    fmt.Sprintf("Gateway of System service. Build at %s", BuildTime),
		Version: Version,
		Run:     onCmdAction,
	}
)

func init() {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 7892, "http server port")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onCmdAction(cmd *cobra.Command, args []string) {
	codePushAdapter := code_push.New(func(options *code_push.Options) {
		options.ServerAddr = ":7890"
	})

	codePushConnErr := codePushAdapter.Conn()
	if codePushConnErr != nil {
		os.Exit(1)
		return
	}
	defer codePushAdapter.Close()

	useCase := usecase.NewUseCase(usecase.CtorConfig{CodePushAdapter: codePushAdapter}, func(options *usecase.Options) {

	})

	server := http.New(useCase, func(options *http.Options) {
		options.Port = port
	})

	httpServeErr := server.ListenAndServe()
	if httpServeErr != nil {
		os.Exit(1)
	}
}
