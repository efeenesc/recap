package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"recap/internal/config"
	"reflect"
	"strconv"
)

type Info struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

var defaultInfo = map[string]string{
	"Version":                "0.0.1",
	"FirstTimeTutorialShown": "0",
}

// Inserts a key-value pair into the info table of the provided database.
// ! This is unused
// Parameters:
//   - value: A string representing the value to be inserted
func WriteInfo(key, value string) error {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
		return err
	}
	defer dbCl.Close()

	_, err = dbCl.Exec("INSERT INTO info (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		return fmt.Errorf("error inserting info: %v", err)
	}
	return nil
}

// Updates a specific info row in the database using the provided key and new value.
// It performs a SQL UPDATE operation and returns an error if the update fails.
func updateInfo(db *sql.DB, key, newValue string) error {
	_, err := db.Exec("UPDATE info SET value = ? WHERE key = ?", newValue, key)
	return err
}

// Updates info in both the database and the in-memory configuration (config.Info) using reflection.
// It ensures that changes to info are saved persistently and reflected immediately in the running application.
func UpdateInfo(newInfo map[string]string) error {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
		return err
	}
	defer dbCl.Close()

	for key, val := range newInfo {
		err = updateInfo(dbCl, key, val)
		if err != nil {
			fmt.Printf("Error when updating %s with %s: %v\n", key, val, err.Error())
			return err
		}

		// Update the config.Config struct via reflection. Find field with 'key', then update the field with 'val'
		r := reflect.ValueOf(&config.Info).Elem()
		field := r.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Int:
				intVal, err := strconv.Atoi(val)
				if err != nil {
					log.Panic("Could not parse int to string during setting reflection")
				}
				field.SetInt(int64(intVal))

			default:
				field.SetString(val)
			}

		} else {
			fmt.Printf("Field %s not found or not settable\n", key)
		}
	}

	return nil
}

// Retrieves an Info record from the database based on the provided key.
//
// Parameters:
//   - key - A string representing the key to search for in the info table
func ReadInfo(key string) (*Info, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
		return nil, err
	}
	defer dbCl.Close()

	row := dbCl.QueryRow("SELECT key, value FROM info WHERE key = ?", key)

	var info Info
	if err := row.Scan(&info.Key, &info.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No rows found
		}
		return nil, fmt.Errorf("error reading info: %v", err)
	}
	return &info, nil
}

func ReadInfoToStruct(infoStruct *config.AppInfo, infoMap map[string]string) *config.AppInfo {
	for key, val := range infoMap {
		switch key {
		case "Version":
			infoStruct.Version = val
		case "FirstTimeTutorialShown":
			infoStruct.FirstTimeTutorialShown = val
		}
	}

	return infoStruct
}

// Retrieves all records from the "info" table in the database.
func ReadAllInfo() (map[string]string, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
		return nil, err
	}
	defer dbCl.Close()

	rows, err := dbCl.Query("SELECT * FROM info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := make(map[string]string)
	for rows.Next() {
		var s Info
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		info[s.Key] = s.Value
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	ReadInfoToStruct(&config.Info, info)

	return info, nil
}

// Fills the info table in the database with default values.
// It starts a new transaction, inserts the default key-value pairs into the info table,
// and commits the transaction. If any error occurs during the process, the transaction
// is rolled back and an error is returned.
func InitializeInfo(dbCl *sql.DB) error {
	tx, err := dbCl.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start database transaction: %v", err)
	}

	for k, v := range defaultInfo {
		if _, err := tx.Exec("INSERT INTO info (key, value) VALUES (?, ?)", k, v); err != nil {
			tx.Rollback() // nolint: all
			return fmt.Errorf("error inserting info (%s, %s): %v", k, v, err)
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback() // nolint: all
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
