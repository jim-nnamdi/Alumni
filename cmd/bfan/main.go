package main

import (
	"fmt"
	"os"

	"github.com/jim-nnamdi/bashfans/pkg/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bfans-server",
		Usage: "start a server for bfans",
		Commands: []*cli.Command{
			command.StartCommand(),
		},
		Version: "v0.1.5",
		Authors: []*cli.Author{
			{
				Name:  "Jim Samuel Nnamdi",
				Email: "jsamuel@bfans.net",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error running program:", err.Error())
		os.Exit(1)
	}
}
