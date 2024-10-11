package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Logs a daily report into the database and updates the associated captures with the report ID.
//
// Parameters:
//   - db: A pointer to the sql.DB connection. If nil, a new connection will be created.
//   - reportText: The content of the daily report.
//   - caps: A slice of CaptureDescription containing the captures to be associated with the report.
//   - genWithApi: A string indicating the API used to generate the report.
//   - genWithModel: A string indicating the model used to generate the report.
//
// Note: The function uses dynamic SQL placeholders for the IN clause to update the captures.
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

// Retrieves all reports from the database with a report_id greater than the specified id.
// The results are ordered by timestamp in descending order.
//
// Parameters:
//   - id: The report_id threshold
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

// Retrieves a list of reports from the database that have a report_id less than the specified id.
// The results are ordered by timestamp in descending order and limited to the specified number of reports.
//
// Parameters:
//   - id: The report_id threshold
//   - limit: The maximum number of reports to retrieve
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

// Retrieves a list of reports from the database, limited by the specified number.
//
// Parameters:
//   - limit: The maximum number of reports to retrieve
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

	return results, nil
}

// GetReportById retrieves a report from the database by its ID.
//
// Parameters:
//   - id: The ID of the report to retrieve
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

// DeleteReportsById deletes reports from the database based on the provided list of report IDs.
//
// Parameters:
//   - ids []int - A slice of report IDs to be deleted
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
