package db

import (
	"database/sql"
	"fmt"
	"log"
	"rcallport/internal/utils"
	"time"
)

func InsertCapture(db *sql.DB, ssFilenames []string) int64 {
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

	return capt_id
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

func getLastScreenshots(db *sql.DB, limit int) ([]CaptureScreenshotImage, error) {
	// Perform the query using placeholders and pass 'limit' as a parameter
	rows, err := db.Query(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename
		FROM 
			captures c
		INNER JOIN 
			screenshots s ON c.capture_id = s.capt_id
		ORDER BY 
			c.timestamp DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err) // Return error instead of log.Fatal
	}
	defer rows.Close() // Always close the rows after use

	var results []CaptureScreenshotImage

	for rows.Next() {
		var cs CaptureScreenshotImage
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.Description,
			&cs.Filename, // Make sure 'cs.Screenshot' matches the correct field name/type
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, cs)
	}

	// Check for errors after row iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return results, nil
}

func GetScreenshots(limit int) ([]CaptureScreenshotImage, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("could not connect to DB: %v", err)
	}
	defer dbCl.Close()

	screenshots, err := getLastScreenshots(dbCl, limit)
	if err != nil {
		return nil, fmt.Errorf("could not get screenshots: %v", err)
	}

	for i := range screenshots {
		screenshots[i].Screenshot = utils.ReadImageToBase64(screenshots[i].Filename)
	}

	fmt.Printf("Number of screenshots: %d\n", len(screenshots))

	return screenshots, nil
}

func GetScreenshotById(id int) (*CaptureScreenshotImage, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("Could not connect to DB: %v", err)
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename
		FROM 
			captures c
		INNER JOIN 
			screenshots s ON c.capture_id = s.capt_id
		WHERE c.capture_id = ?
		ORDER BY 
			c.timestamp DESC
		LIMIT ?
	`, id, 1)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err) // Return error instead of log.Fatal
	}
	defer rows.Close() // Always close the rows after use

	var cs CaptureScreenshotImage

	for rows.Next() {
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.Description,
			&cs.Filename, // Make sure 'cs.Screenshot' matches the correct field name/type
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}

	// Check for errors after row iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &cs, nil
}

func GetScreenshotByIds(db *sql.DB, ids []int) ([]CaptureScreenshot, error) {
	if db == nil {
		db, err := CreateConnection()
		if err != nil {
			return nil, fmt.Errorf("could not connect to DB: %v", err)
		}
		defer db.Close()
	}

	questionMarks := generateNumOfQuestionMarks(len(ids))
	preparedQuery := fmt.Sprintf(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename
		FROM 
			captures c
		INNER JOIN 
			screenshots s ON c.capture_id = s.capt_id
		WHERE c.capture_id IN (%s)
		ORDER BY 
			c.timestamp DESC
	`, questionMarks)

	// Create a slice for the query arguments
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.Query(preparedQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []CaptureScreenshot

	for rows.Next() {
		var cs CaptureScreenshot
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.Description,
			&cs.Filename, // Make sure 'cs.Screenshot' matches the correct field name/type
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, cs)
	}

	// Check for errors after row iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return results, nil
}
