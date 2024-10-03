package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

/*
Logs the daily report by first inserting a row into 'dailyreports', then getting the last inserted row
and updating all rows in 'captures' which were used to generate the Daily Report.
*/
func LogDailyReport(db *sql.DB, reportText string, caps []CaptureDescription, genWithApi string, genWithModel string) (*int64, error) {
	if db == nil {
		db, err := CreateConnection()
		if err != nil {
			return nil, err
		}
		defer db.Close()
	}
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
		INSERT INTO dailyreports (timestamp, content, gen_with_api, gen_with_model)
		VALUES (?, ?, ?, ?)`,
		time.Now().UTC().Unix(), reportText, genWithApi, genWithModel)

	if err != nil {
		fmt.Printf("Error logging daily report: %v\n", err)
		return nil, err
	}

	// Get the last inserted ID (daily report ID)
	drId, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert ID: %v\n", err)
		return nil, err
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

	return &drId, nil
}

func GetReportsNewerThan(id int) ([]Report, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
		SELECT * FROM dailyreports
		WHERE 
			report_id > (?)
		ORDER BY 
			timestamp DESC
	`, id)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err) // Return error instead of log.Fatal
	}
	defer rows.Close() // Always close the rows after use

	var results []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(
			&r.ReportID,
			&r.Timestamp,
			&r.Content,
			&r.GenWithApi,
			&r.GenWithModel,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, r)
	}

	return results, nil
}

func GetReportsOlderThan(id int, limit int) ([]Report, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
		SELECT * FROM dailyreports
		WHERE 
			report_id < (?)
		ORDER BY 
			timestamp DESC
		LIMIT ?
	`, id, limit)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var results []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(
			&r.ReportID,
			&r.Timestamp,
			&r.Content,
			&r.GenWithApi,
			&r.GenWithModel,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, r)
	}

	return results, nil
}

func GetReports(limit int) ([]Report, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		return nil, err
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
	SELECT * FROM dailyreports
	ORDER BY 
		timestamp DESC
	LIMIT ?`, limit)

	if err != nil {
		return nil, err
	}

	var results []Report

	for rows.Next() {
		var r Report
		err := rows.Scan(
			&r.ReportID,
			&r.Timestamp,
			&r.Content,
			&r.GenWithApi,
			&r.GenWithModel,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		results = append(results, r)
	}

	fmt.Println(len(results))

	return results, nil
}

func GetReportById(id int) (*Report, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		return nil, err
	}
	defer dbCl.Close()

	rows, err := dbCl.Query(`
	SELECT * FROM dailyreports
	WHERE report_id = ?`, id)

	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err) // Return error instead of log.Fatal
	}
	defer rows.Close() // Always close the rows after use

	var rep Report

	for rows.Next() {
		err := rows.Scan(
			&rep.ReportID,
			&rep.Timestamp,
			&rep.Content,
			&rep.GenWithApi,
			&rep.GenWithModel,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
	}

	// Check for errors after row iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return &rep, nil
}

func DeleteReportsById(ids []int) error {
	dbCl, err := CreateConnection()
	if err != nil {
		return fmt.Errorf("could not connect to DB: %v", err)
	}
	defer dbCl.Close()

	questionMarks := generateNumOfQuestionMarks(len(ids))
	deleteQuery := fmt.Sprintf(`
		DELETE FROM dailyreports
		WHERE report_id IN (%s)
	`, questionMarks)

	// Create a slice for the query arguments
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = dbCl.Exec(deleteQuery, args...)
	if err != nil {
		return fmt.Errorf("error deleting reports: %v", err)
	}

	return nil
}
