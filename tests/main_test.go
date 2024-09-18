package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/user-manager/app"
	"github.com/user-manager/config"
	"github.com/user-manager/db"
	"github.com/user-manager/models"
)

var testApp *app.App

func TestMain(m *testing.M) {

	cfg, err := config.LoadTestConfig()

	if err != nil {
		log.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		log.Fatalf("couldn't initialize db: %v", err)
	}
	logChannel := make(chan models.RequestLog)

	testApp = app.BuildApp(cfg, postgresDB, logChannel)
	// Setup test database
	if err := setupTestDB(testApp); err != nil {
		log.Fatalf("Error setting up test database: %v", err)
	}
	go testRequestsLogger(testApp)
	defer teardownTestDB(testApp.PostgresDB().Conn)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setupTestDB(testApp *app.App) error {

	if err := resetTestDB(testApp.PostgresDB().Conn); err != nil {
		return err
	}
	if err := seedTestData(testApp.PostgresDB().Conn); err != nil {
		return err
	}

	return nil
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`
		DROP SCHEMA IF EXISTS public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		return fmt.Errorf("error dropping schema: %v", err)
	}
	if err := createTables(db); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func createTables(db *sql.DB) error {

	// Read SQL statements from file
	sqlFile, err := os.ReadFile("../db/db.sql")
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Create necessary tables
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func seedTestData(db *sql.DB) error {
	// Insert initial test data
	_, err := db.Exec(`
		INSERT INTO users (name, age,email) VALUES
			('Josh', 40,'josh@gmail.com'),
			('John', 35,'john@gmail.com');
	`)
	if err != nil {
		return fmt.Errorf("error seeding test data: %v", err)
	}
	return nil
}

func teardownTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// requestLogger listens on the logChannel and logs request durations
func testRequestsLogger(app *app.App) {
	for reqLog := range app.LogChannel() {
		log.Printf("Method: %s, Path: %s, Duration: %v\n",
			reqLog.Method, reqLog.Path, reqLog.Duration)
	}
}
