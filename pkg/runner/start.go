package runner

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	database "github.com/jim-nnamdi/bashfans/pkg/database/mysql"
	"github.com/jim-nnamdi/bashfans/pkg/handlers"
	"github.com/jim-nnamdi/bashfans/pkg/server"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type StartRunner struct {
	ListenAddr string

	LoggingProduction      bool
	LoggingOutputPath      string
	LogggingLevel          string
	ErrorLoggingOutputPath string

	MySQLDatabaseHost     string
	MySQLDatabasePort     string
	MySQLDatabaseUser     string
	MySQLDatabasePassword string
	MySQLDatabaseName     string
}

func (runner *StartRunner) Run(c *cli.Context) error {
	var (
		loggerConfig        = zap.NewDevelopmentConfig()
		logger              *zap.Logger
		err                 error
		mysqlDbInstance     *sql.DB
		mysqlDatabaseClient database.Database
	)
	if runner.LoggingProduction {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.OutputPaths = []string{runner.LoggingOutputPath}
		loggerConfig.ErrorOutputPaths = []string{runner.ErrorLoggingOutputPath}
	}

	if err = loggerConfig.Level.UnmarshalText([]byte(runner.LogggingLevel)); err != nil {
		return err
	}

	if logger, err = loggerConfig.Build(); err != nil {
		return err
	}

	logger.Sync()
	databaseConfig := &mysql.Config{
		User:                 runner.MySQLDatabaseUser,
		Passwd:               runner.MySQLDatabasePassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", runner.MySQLDatabaseHost, runner.MySQLDatabasePort),
		DBName:               runner.MySQLDatabaseName,
		AllowNativePasswords: true,
	}

	if mysqlDbInstance, err = sql.Open("mysql", databaseConfig.FormatDSN()); err != nil {
		return fmt.Errorf("unable to open connection to MySQL Server: %s", err.Error())
	}

	// logger.Debug("MYSQL Connection", zap.String("dsn", databaseConfig.FormatDSN()), zap.Any("mysqlDB", mysqlDB))
	if mysqlDatabaseClient, err = database.NewMySQLDatabase(mysqlDbInstance); err != nil {
		return fmt.Errorf("unable to create MySQL database client: %s", err.Error())
	}
	server := &server.GracefulShutdownServer{
		HTTPListenAddr:  runner.ListenAddr,
		RegisterHandler: handlers.NewRegisterHandler(logger, mysqlDatabaseClient),
		LoginHandler:    handlers.NewLoginHandler(logger, mysqlDatabaseClient),
		ProfileHandler:  handlers.NewProfileHandler(logger, mysqlDatabaseClient),
		HomeHandler:     handlers.NewHomeHandler(),
	}
	server.Start()
	return nil
}
