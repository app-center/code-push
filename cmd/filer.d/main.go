package main

import (
	"fmt"
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
		Use:     "Filer",
		Short:   "Filer daemon of Code-Push service",
		Long:    fmt.Sprintf("Filer daemon of Code-Push service. Build at %s", BuildTime),
		Version: Version,
		Run:     nil,
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
