package main

import (
	"fmt"
	"os"

	"github.com/jim-nnamdi/jinx/pkg/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Alumni network server",
		Usage: "networking platform for alumnis",
		Commands: []*cli.Command{
			command.StartCommand(),
		},
		Version: "v0.1.5",
		Authors: []*cli.Author{
			{
				Name:  "csc group",
				Email: "csc432@lasu.com",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error running program:", err.Error())
		os.Exit(1)
	}
}
