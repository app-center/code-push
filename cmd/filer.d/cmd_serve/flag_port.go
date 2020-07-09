package cmd_serve

import "github.com/urfave/cli/v2"

const flagPortName = "port"
const flagPortDefault = 7891

var flagPort int

var useFlagPort = func(flags []cli.Flag) []cli.Flag {
	return append(flags, &cli.IntFlag{
		Name:        flagPortName,
		Aliases:     []string{"p"},
		Usage:       "grpc server port",
		EnvVars:     []string{"CODE_PUSH__FILER__PORT"},
		Required:    false,
		Value:       flagPortDefault,
		Destination: &flagPort,
	})
}
