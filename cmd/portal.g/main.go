package main

import (
	"fmt"
	code_push "github.com/funnyecho/code-push/gateway/portal/adapter/code-push"
	"github.com/funnyecho/code-push/gateway/portal/adapter/filer"
	"github.com/funnyecho/code-push/gateway/portal/adapter/session"
	"github.com/funnyecho/code-push/gateway/portal/interface/http"
	"github.com/funnyecho/code-push/gateway/portal/usecase"
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
		Use:     "Portal Gateway",
		Short:   "Gateway of Portal service",
		Long:    fmt.Sprintf("Gateway of Portal service. Build at %s", BuildTime),
		Version: Version,
		Run:     onCmdAction,
	}
)

func init() {
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 7882, "http server port")
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

	sessionAdapter := session.New(func(options *session.Options) {
		options.ServerAddr = ":7892"
	})
	sessionConnErr := sessionAdapter.Conn()
	if sessionConnErr != nil {
		os.Exit(1)
		return
	}
	defer sessionAdapter.Close()

	filerAdapter := filer.New(func(options *filer.Options) {
		options.ServerAddr = ":7891"
	})
	filerConnErr := filerAdapter.Conn()
	if filerConnErr != nil {
		os.Exit(1)
		return
	}
	defer filerAdapter.Close()

	uc := usecase.NewUseCase(
		usecase.CtorConfig{
			CodePushAdapter: codePushAdapter,
			SessionAdapter:  sessionAdapter,
			FilerAdapter:    filerAdapter,
		},
		func(options *usecase.Options) {

		},
	)

	server := http.New(uc, func(options *http.Options) {
		options.Port = port
	})

	httpServeErr := server.ListenAndServe()
	if httpServeErr != nil {
		os.Exit(1)
	}
}