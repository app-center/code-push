package cmd_serve

import (
	"fmt"
	"github.com/funnyecho/code-push/cmd/filer.d/internal"
	"github.com/urfave/cli/v2"
)

var UseCommand = func(commands cli.Commands) cli.Commands {
	cmd := &cli.Command{
		Name:         "serve",
		Usage:        "start grpc server of filer",
		Action:       onAction,
		OnUsageError: onUsageError,
	}

	internal.UseCommandFlags(cmd, useFlagPort)

	return append(commands, cmd)
}

func onAction(c *cli.Context) error {
	fmt.Printf("kadjflkajsfd %v", flagPort)
	return nil
}

func onUsageError(context *cli.Context, err error, isSubcommand bool) error {
	return err
}
