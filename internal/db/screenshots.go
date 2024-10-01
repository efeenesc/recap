package db

import (
	"database/sql"
	"fmt"
	"log"
	"rcallport/internal/utils"
	"time"
)

type FullThumbScrPair struct {
	Full, Thumb string
}

// Inserts a new capture record into the database and associates it with the provided screenshot filenames.
// It returns the ID of the newly created capture or logs a fatal error if an operation fails.
// Currently each capture equates to just one screenshot file, but this might change in the future
func InsertCapture(db *sql.DB, scrFullThumbPairs []FullThumbScrPair) int64 {
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
	insertScreenshots(tx, scrFullThumbPairs, capt_id) // Pass the capture row's ID as well as the transaction
	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}

	return capt_id
}

// Updates the description of a specific screenshot identified by its ID.
// Returns the result of the update operation or an error if the operation fails
func UpdateScreenshotDescription(db *sql.DB, screenshot_id int, description string) (sql.Result, error) {
	return db.Exec(`
	UPDATE screenshots
	SET description = ? 
	WHERE screenshot_id = ?`, description, screenshot_id)
}

// Inserts one or multiple screenshot records into the database within a transaction,
// associating them with the given capture ID.
func insertScreenshots(tx *sql.Tx, scrs []FullThumbScrPair, capt_id int64) {
	stmt, err := tx.Prepare(`
	INSERT INTO 
		screenshots(
			filename,
			thumbname,
			capt_id,
			description
		)
	VALUES (
		?, ?, ?, ?
	)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, el := range scrs {
		_, err = stmt.Exec(el.Full, el.Thumb, capt_id, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Retrieves all captures from today, including their associated descriptions.
// It returns a list of CaptureDescription objects or an error if the operation fails.
// It does so by first getting the UNIX second timestamp equivalent of 12AM today, then filtering
// rows' timestamp values to be higher than today's 12AM timestamp
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

// Retrieves all screenshots that have not been processed by description generation via a vision model yet.
// It returns a list of CaptureScreenshot objects or an error if the operation fails
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

func GetScreenshotsNewerThan(id int) ([]CaptureScreenshotImage, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename,
			s.thumbname
		FROM 
			captures c
		INNER JOIN 
			screenshots s ON c.capture_id = s.capt_id
		WHERE 
			c.capture_id > (?)
		ORDER BY 
			c.timestamp DESC
	`, id)

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
			&cs.Filename,
			&cs.Thumbname,
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

	for i := range results {
		results[i].Screenshot = utils.ReadImageToBase64PreferThumb(results[i].Filename, results[i].Thumbname)
	}

	return results, nil
}

func GetScreenshotsOlderThan(id int, limit int) ([]CaptureScreenshotImage, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename,
			s.thumbname
		FROM 
			captures c
		INNER JOIN 
			screenshots s ON c.capture_id = s.capt_id
		WHERE 
			c.capture_id < (?)
		ORDER BY 
			c.timestamp DESC
		LIMIT ?
	`, id, limit)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []CaptureScreenshotImage

	for rows.Next() {
		var cs CaptureScreenshotImage
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.Description,
			&cs.Filename, // Make sure 'cs.Screenshot' matches the correct field name/type
			&cs.Thumbname,
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

	for i := range results {
		results[i].Screenshot = utils.ReadImageToBase64PreferThumb(results[i].Filename, results[i].Thumbname)
	}

	return results, nil
}

// getLastScreenshots retrieves a limited number of the most recent screenshots from the database.
// It returns a slice of CaptureScreenshotImage objects or an error if the operation fails.
func getLastScreenshots(db *sql.DB, limit int) ([]CaptureScreenshotImage, error) {
	// Perform the query using placeholders and pass 'limit' as a parameter
	rows, err := db.Query(`
		SELECT 
			c.capture_id,
			c.timestamp, 
			s.description,
			s.filename,
			s.thumbname
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
			&cs.Thumbname,
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

// Reads images' thumbnails
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
		screenshots[i].Screenshot = utils.ReadImageToBase64PreferThumb(screenshots[i].Filename, screenshots[i].Thumbname)
	}

	fmt.Printf("Number of screenshots: %d\n", len(screenshots))

	return screenshots, nil
}

// Reads the full image
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
			s.filename,
			s.thumbname
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
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var cs CaptureScreenshotImage

	for rows.Next() {
		err := rows.Scan(
			&cs.CaptureID,
			&cs.Timestamp,
			&cs.Description,
			&cs.Filename,
			&cs.Thumbname,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}

	cs.Screenshot = utils.ReadImageToBase64PreferFull(cs.Filename, cs.Thumbname)

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
			s.filename,
			s.thumbname
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
			&cs.Filename,
			&cs.Thumbname,
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
