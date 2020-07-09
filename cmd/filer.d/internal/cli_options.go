package internal

import "github.com/urfave/cli/v2"

func UseCommands(app *cli.App, cmds ...ICommandOptionFunc) {
	for _, cmd := range cmds {
		app.Commands = cmd(app.Commands)
	}
}

func UseAppFlags(app *cli.App, flags ...IFlagOptionFunc) {
	for _, flag := range flags {
		app.Flags = flag(app.Flags)
	}
}

func UseCommandFlags(cmd *cli.Command, flags ...IFlagOptionFunc) {
	for _, flag := range flags {
		cmd.Flags = flag(cmd.Flags)
	}
}

type ICommandOptionFunc = func(commands cli.Commands) cli.Commands
type IFlagOptionFunc = func(flags []cli.Flag) []cli.Flag
