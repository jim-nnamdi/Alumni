package command

import (
	"github.com/jim-nnamdi/jinx/pkg/runner"
	"github.com/urfave/cli/v2"
)

func StartCommand() *cli.Command {
	var (
		startRunner = &runner.StartRunner{}
	)

	cmd := &cli.Command{
		Name:  "start",
		Usage: "starts the server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "listen-addr",
				EnvVars:     []string{"LISTEN_ADDR"},
				Usage:       "the address that the server will listen for request on",
				Destination: &startRunner.ListenAddr,
				Value:       ":8080", // TODO: check that this is correct port to serve on
			},
			&cli.StringFlag{
				Name:        "mysql-database-name",
				EnvVars:     []string{"AC_DBNAME"},
				Usage:       "Sample database name",
				Destination: &startRunner.MySQLDatabaseName,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-password",
				EnvVars:     []string{"AC_PASSWORD"},
				Usage:       "Sample database password",
				Destination: &startRunner.MySQLDatabasePassword,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-User",
				EnvVars:     []string{"AC_USER"},
				Usage:       "Sample database user",
				Destination: &startRunner.MySQLDatabaseUser,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Host",
				EnvVars:     []string{"AC_HOST"},
				Usage:       "Sample database host",
				Destination: &startRunner.MySQLDatabaseHost,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Port",
				EnvVars:     []string{"AC_PORT"},
				Usage:       "Sample database port",
				Destination: &startRunner.MySQLDatabasePort,
				Value:       "",
			},
		},
		Action: startRunner.Run,
	}
	return cmd
}
