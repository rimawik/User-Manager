package app

import (
	"github.com/user-manager/config"
	"github.com/user-manager/db"
	"github.com/user-manager/models"
)

type App struct {
	config     *config.Config
	postgresDB *db.DB
	logChannel chan models.RequestLog
}

func (app *App) Conf() *config.Config {
	return app.config
}

func (app *App) PostgresDB() *db.DB {
	return app.postgresDB
}
func (app *App) LogChannel() chan models.RequestLog {
	return app.logChannel
}


func BuildApp(cfg *config.Config,
	postgres *db.DB, logChannel chan models.RequestLog) *App {
	return &App{config: cfg,
		postgresDB: postgres,
		logChannel:logChannel,
	}
}
