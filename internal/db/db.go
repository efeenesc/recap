package db

import (
	"database/sql"
	"log"
	"path"
	"rcallport/internal/config"

	_ "github.com/mattn/go-sqlite3"
)

func createTable(db *sql.DB) error {
	capturesStmt := `
	create table captures (capture_id integer not null primary key, r_id integer, timestamp integer not null)
	`
	_, err := db.Exec(capturesStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, capturesStmt)
	}

	screenshotsStmt := `
	create table screenshots (screenshot_id integer not null primary key, capt_id integer not null, filename text, description text, FOREIGN KEY(capt_id) REFERENCES captures(capture_id))
	`

	_, err = db.Exec(screenshotsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, screenshotsStmt)
	}

	dailyReportsStmt := `
	create table dailyreports (report_id integer not null primary key, timestamp integer not null, content text)
	`

	_, err = db.Exec(dailyReportsStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, dailyReportsStmt)
	}

	settingsStmt := `
		CREATE TABLE settings (
			id INTEGER PRIMARY KEY,
			key TEXT UNIQUE NOT NULL,
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
	db, err := sql.Open("sqlite3", path.Join(proot, config.Config.Database.Path))

	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func Initialize(keepConnectionOpen bool) (*sql.DB, error) {
	db, _ := CreateConnection()
	createTable(db)

	if keepConnectionOpen {
		return db, nil
	}

	db.Close()
	return nil, nil
}
