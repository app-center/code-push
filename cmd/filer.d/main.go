package main

import (
	"github.com/funnyecho/code-push/cmd/filer.d/cmd_serve"
	"github.com/funnyecho/code-push/cmd/filer.d/internal"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var (
	Version   = "development"
	BuildTime = time.Now()
)

func main() {
	app := &cli.App{
		Name:                   "Filer",
		Usage:                  "Filer for uploading file and file management",
		Version:                Version,
		Commands:               nil,
		Flags:                  nil,
		EnableBashCompletion:   true,
		Action:                 nil,
		OnUsageError:           nil,
		Compiled:               BuildTime,
		Authors:                nil,
		Copyright:              "",
		Writer:                 nil,
		ErrWriter:              nil,
		ExitErrHandler:         nil,
		Metadata:               nil,
		ExtraInfo:              nil,
		CustomAppHelpTemplate:  "",
		UseShortOptionHandling: false,
	}

	internal.UseCommands(app, cmd_serve.UseCommand)

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
