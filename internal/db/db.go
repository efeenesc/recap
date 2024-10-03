package db

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"rcallport/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

var Initializers InitializerCallbacks

type InitializerCallbacks struct {
	FunctionsGiven bool
	InitSchedule   func()
	InitLLM        func()
}

func NewInitializers() *InitializerCallbacks {
	return &InitializerCallbacks{}
}

func createTable(db *sql.DB) error {
	capturesStmt := `
	CREATE TABLE captures (
		capture_id INTEGER NOT NULL PRIMARY KEY, 
		r_id INTEGER, 
		timestamp INTEGER NOT NULL
	);
	`
	_, err := db.Exec(capturesStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, capturesStmt)
	}

	screenshotsStmt := `
	CREATE TABLE screenshots (
		screenshot_id INTEGER NOT NULL PRIMARY KEY, 
		capt_id INTEGER NOT NULL,
		filename TEXT, 
		thumbname TEXT, 
		description TEXT, 
		gen_with_api TEXT,
		gen_with_model TEXT,
		FOREIGN KEY(capt_id) REFERENCES captures(capture_id)
	);
	`

	_, err = db.Exec(screenshotsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, screenshotsStmt)
	}

	dailyReportsStmt := `
	CREATE TABLE dailyreports (
		report_id INTEGER NOT NULL PRIMARY KEY,
		timestamp INTEGER NOT NULL,
		content TEXT,
		gen_with_api TEXT,
		gen_with_model TEXT,
	);
	`

	_, err = db.Exec(dailyReportsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, dailyReportsStmt)
	}

	settingsStmt := `
		CREATE TABLE settings (
			key TEXT PRIMARY KEY UNIQUE NOT NULL,
			value TEXT NOT NULL
	);
	`

	_, err = db.Exec(settingsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, settingsStmt)
	}

	return nil
}

func CreateConnection() (*sql.DB, error) {
	proot, _ := config.GetProjectRoot()
	db, err := sql.Open("sqlite3", path.Join(proot, "rcallport.db"))

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func Initialize(keepConnectionOpen bool) (*sql.DB, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Error creating DB connection: %v\n", err.Error())
	}
	createTable(dbCl)

	err = initializeSettings(dbCl, defaultSettings)
	if err != nil {
		fmt.Printf("Encountered the following error when trying to insert defaults: %v\n", err.Error())
		dbCl.Close()
		return nil, err
	}

	_, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err.Error())
	}

	if keepConnectionOpen {
		return dbCl, nil
	}

	dbCl.Close()
	return nil, nil
}
