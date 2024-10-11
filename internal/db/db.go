package db

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"recap/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

var Initializers InitializerCallbacks

// A list of callbacks that are used to initialize parts of the application,
// called directly during the database initialization process. This avoids
// circular imports
type InitializerCallbacks struct {
	FunctionsGiven bool
	InitSchedule   func()
	InitLLM        func()
}

func NewInitializers() *InitializerCallbacks {
	return &InitializerCallbacks{}
}

func createTable(db *sql.DB) {
	capturesStmt := `
	CREATE TABLE IF NOT EXISTS captures (
		capture_id INTEGER NOT NULL PRIMARY KEY, 
		r_id INTEGER, 
		timestamp INTEGER NOT NULL
	);
	`
	_, err := db.Exec(capturesStmt)
	if err != nil {
		log.Printf("Error executing query: %q: %s\n", err, capturesStmt)
	}

	screenshotsStmt := `
	CREATE TABLE IF NOT EXISTS screenshots (
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
		log.Printf("Error executing query: %q: %s\n", err, screenshotsStmt)
	}

	dailyReportsStmt := `
	CREATE TABLE IF NOT EXISTS dailyreports (
		report_id INTEGER NOT NULL PRIMARY KEY,
		timestamp INTEGER NOT NULL,
		content TEXT,
		gen_with_api TEXT,
		gen_with_model TEXT
	);
	`
	_, err = db.Exec(dailyReportsStmt)
	if err != nil {
		log.Printf("Error executing query: %q: %s\n", err, dailyReportsStmt)
	}

	settingsStmt := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY UNIQUE NOT NULL,
		value TEXT NOT NULL
	);
	`
	_, err = db.Exec(settingsStmt)
	if err != nil {
		log.Printf("Error executing query: %q: %s\n", err, settingsStmt)
	}

	infoStmt := `
	CREATE TABLE IF NOT EXISTS info (
		key TEXT PRIMARY KEY UNIQUE NOT NULL,
		value TEXT NOT NULL
	);
	`
	_, err = db.Exec(infoStmt)
	if err != nil {
		log.Printf("Error executing query: %q: %s\n", err, infoStmt)
	}
}

// Establishes a connection to the SQLite database located
// at the project root directory and returns the database handle.
//
// Returns:
//   - *sql.DB: A handle to the SQLite database
//   - error: An error object if there was an issue opening the database
func CreateConnection() (*sql.DB, error) {
	proot := config.GetProjectRoot()
	return sql.Open("sqlite3", path.Join(proot, "recap.db"))
}

// Sets up the database connection, creates necessary tables,
// initializes settings with default values, and loads the configuration.
// If keepConnectionOpen is true, the database connection remains open and
// is returned; otherwise, the connection is closed before returning.
//
// Parameters:
//   - keepConnectionOpen: Whether to keep the database
//     connection open
func Initialize(keepConnectionOpen bool) (*sql.DB, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		log.Fatalf("Error creating DB connection: %v\n", err.Error())
	}

	createTable(dbCl)

	err = initializeSettings(dbCl, defaultSettings)
	if err != nil {
		fmt.Printf("Error when inserting setting defaults: %v\n", err.Error())
	}

	err = InitializeInfo(dbCl)
	if err != nil {
		fmt.Printf("Error when inserting info defaults: %v\n", err.Error())
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
