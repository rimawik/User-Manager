package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/user-manager/app"
	"github.com/user-manager/config"
	"github.com/user-manager/db"
	"github.com/user-manager/models"
	"github.com/user-manager/routes"
)

// @title User API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8080
// @BasePath /
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		log.Fatalf("couldn't initialize db: %v", err)
	}
	logChannel := make(chan models.RequestLog)

	app := app.BuildApp(cfg, postgresDB, logChannel)
	go requestLogger(app)

	r := routes.NewRouter(app)

	log.Println("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

}

// requestLogger listens on the logChannel and logs request durations
func requestLogger(app *app.App) {
	for reqLog := range app.LogChannel() {
		log.Printf("Method: %s, Path: %s, Duration: %v\n",
			reqLog.Method, reqLog.Path, reqLog.Duration)
	}
}
