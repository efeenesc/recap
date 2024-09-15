package db

import (
	"database/sql"
	"fmt"
	"log"
	"path"
	"rcallport/internal/config"
	"strings"
	"time"

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

	return nil
}

/*
Logs the daily report by first inserting a row into 'dailyreports', then getting the last inserted row
and updating all rows in 'captures' which were used to generate the Daily Report.
*/
func LogDailyReport(db *sql.DB, reportText string, caps []CaptureDescription) {
	// Extract capture IDs into a slice
	capIds := make([]int, len(caps))
	for i, cap := range caps {
		capIds[i] = cap.CaptureID
	}

	// Create placeholders for the IN clause
	questionMarks := strings.Repeat("?,", len(capIds))
	questionMarks = strings.TrimSuffix(questionMarks, ",")

	// Insert the daily report
	res, err := db.Exec(`
		INSERT INTO dailyreports (timestamp, content)
		VALUES (?, ?)`,
		time.Now().UTC().Unix(), reportText)

	if err != nil {
		fmt.Printf("Error logging daily report: %v\n", err)
		return
	}

	// Get the last inserted ID (daily report ID)
	drId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert ID: %v\n", err)
		return
	}

	// Update the captures with the report ID
	query := fmt.Sprintf(`
		UPDATE captures
		SET r_id = ?
		WHERE capture_id IN (%s)`, questionMarks)

	// Prepare arguments for db.Exec (drId + capIds...)
	// Passing drId, capIds... would have been better but this will do
	args := make([]interface{}, len(capIds)+1)
	args[0] = drId
	for i, id := range capIds {
		args[i+1] = id
	}

	// Execute the update query with dynamic placeholders
	_, err = db.Exec(query, args...)
	if err != nil {
		fmt.Printf("Error updating captures: %v\n", err)
	}
}

func InsertCapture(db *sql.DB, ssFilenames []string) {
	stmt, err := db.Prepare(`
	INSERT INTO captures(timestamp)
	VALUES (?)`)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(time.Now().UTC().Unix())
	if err != nil {
		log.Fatal(err)
	}
	capt_id, err := res.LastInsertId() // Get the last inserted ID of capture row
	if err != nil {
		log.Fatal(err)
	}
	stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	insertScreenshots(tx, ssFilenames, capt_id) // Pass the capture row's ID as well as the transaction
	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateScreenshotDescription(db *sql.DB, screenshot_id int, description string) (sql.Result, error) {
	return db.Exec(`
	UPDATE screenshots
	SET description = ? 
	WHERE screenshot_id = ?`, description, screenshot_id)
}

func insertScreenshots(tx *sql.Tx, ss_filenames []string, capt_id int64) {
	stmt, err := tx.Prepare(`
	INSERT INTO 
		screenshots(
			filename,
			capt_id,
			description
		)
	VALUES (
		?, ?, ?
	)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, el := range ss_filenames {
		_, err = stmt.Exec(el, capt_id, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetCapturesToday(db *sql.DB) ([]CaptureDescription, error) {
	now := time.Now().UTC()
	y, m, d := now.Date()
	startOfDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()

	rows, err := db.Query(`
	SELECT 
		c.capture_id,
		c.timestamp, 
		s.description
	FROM 
		captures c
	INNER JOIN 
		screenshots s ON c.capture_id = s.capt_id
	WHERE 
		c.r_id IS NULL & c.timestamp > (?)
	ORDER BY 
		c.timestamp ASC
	`, startOfDay)
	if err != nil {
		log.Fatal(err)
	}

	var results []CaptureDescription

	for rows.Next() {
		var cd CaptureDescription
		err := rows.Scan(
			&cd.CaptureID,
			&cd.Timestamp,
			&cd.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, cd)
	}

	return results, nil
}

func GetUnprocessedCaptures(db *sql.DB) ([]CaptureScreenshot, error) {
	rows, err := db.Query(`
	SELECT 
		c.capture_id, 
		c.timestamp, 
		s.screenshot_id, 
		s.filename, 
		s.description
	FROM 
		captures c
	INNER JOIN 
		screenshots s ON c.capture_id = s.capt_id
	WHERE 
		s.description IS NULL
	ORDER BY 
		c.timestamp DESC
	`)
	if err != nil {
		log.Fatal(err)
	}

	var results []CaptureScreenshot

	for rows.Next() {
		var cs CaptureScreenshot
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.ScreenshotID,
			&cs.Filename,
			&cs.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, cs)
	}

	return results, nil
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
