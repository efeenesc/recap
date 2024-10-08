package db

import (
	"database/sql"
	"fmt"
	"log"
	"recap/internal/config"
	"reflect"
	"strconv"
)

type Setting struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// Display properties for settings
type SettingDisplayProps struct {
	DisplayName string `json:"DisplayName"`
	Description string `json:"Description"`
	Category    string `json:"Category"`
	InputType   string `json:"InputType"`
}

var defaultSettings = map[string]string{
	"ScrPath":                   "./screenshots",
	"DescGenAPI":                "gemini",
	"DescGenModel":              "gemini-1.5-flash",
	"DescGenPrompt":             "This image was captured on a user's computer. Describe what the user was working on. Do not expose passwords, other people's names, emails, and other private and secure information.",
	"DescGenIntervalMins":       "120", // Default interval in minutes
	"DescGenIntervalEnabled":    "1",   // 1 for enabled, 0 for disabled
	"ScreenshotIntervalMins":    "10",  // Default interval in minutes
	"ScreenshotIntervalEnabled": "1",   // 1 for enabled, 0 for disabled
	"ReportAPI":                 "gemini",
	"ReportModel":               "gemini-1.5-flash",
	"ReportAutoEnabled":         "0",
	"ReportAutoAt":              "17:00", // Default auto report time
	"ReportPrompt":              "You are an AI assistant tasked with generating a daily activity report for a user based on a series of visual descriptions captured from their computer screen throughout the day. Your job is to summarize this data into brief items describing what the user worked on today.",
	"OllamaURL":                 "http://localhost:11434",
	"GeminiAPIKey":              "your-gemini-api-key",
}

func GetDisplayValues() map[string]SettingDisplayProps {
	var settingKeyDisplayVals = map[string]SettingDisplayProps{
		"ScrPath":                   {DisplayName: "Path", Description: "Specify the directory where screenshots will be saved on your device", Category: "Screenshots", InputType: "FolderPicker"},
		"DescGenAPI":                {DisplayName: "API", Description: "Select the AI service to use for generating descriptions of your screenshots", Category: "Vision", InputType: "APIPicker"},
		"DescGenModel":              {DisplayName: "Model", Description: "Choose the specific AI model for analyzing screenshots and generating descriptions", Category: "Vision", InputType: "APIModelPicker"},
		"DescGenPrompt":             {DisplayName: "Prompt", Description: "Customize the instructions given to the AI when generating screenshot descriptions", Category: "Vision", InputType: "ExtendedTextInput"},
		"DescGenIntervalEnabled":    {DisplayName: "Schedule", Description: "Toggle automatic description generation after Recap starts", Category: "Vision", InputType: "Boolean"},
		"DescGenIntervalMins":       {DisplayName: "Interval", Description: "Set how often (in minutes) screenshots should be automatically sent for description generation", Category: "Vision", InputType: "NumberInput"},
		"ScreenshotIntervalEnabled": {DisplayName: "Schedule", Description: "Toggle automatic screenshot capturing at regular intervals after Recap starts", Category: "Screenshots", InputType: "Boolean"},
		"ScreenshotIntervalMins":    {DisplayName: "Interval", Description: "Define how frequently (in minutes) automatic screenshots should be taken", Category: "Screenshots", InputType: "NumberInput"},
		"ReportAPI":                 {DisplayName: "API", Description: "Choose the AI service for generating reports based on your screenshot descriptions", Category: "Reports", InputType: "APIPicker"},
		"ReportModel":               {DisplayName: "Model", Description: "Select the specific AI model for creating reports from your screenshot descriptions", Category: "Reports", InputType: "APIModelPicker"},
		"ReportAutoEnabled":         {DisplayName: "Schedule", Description: "Enable or disable automatic daily report generation", Category: "Reports", InputType: "Boolean"},
		"ReportAutoAt":              {DisplayName: "Time", Description: "Set the specific time each day when an automatic report should be generated", Category: "Reports", InputType: "TimePicker"},
		"ReportPrompt":              {DisplayName: "Prompt", Description: "Customize the instructions given to the AI when generating reports from your screenshot descriptions", Category: "Reports", InputType: "ExtendedTextInput"},
		"OllamaURL":                 {DisplayName: "Ollama URL", Description: "Enter the URL (including port) for your Ollama instance. The default is http://localhost:11434", Category: "Models", InputType: "URLInput"},
		"GeminiAPIKey":              {DisplayName: "Gemini API key", Description: "Enter your Gemini API key. You can obtain a free API key from Google. For instructions, please refer to the 'Setting up' section in the tutorial", Category: "Models", InputType: "TextInput"},
	}

	return settingKeyDisplayVals
}

// Retrieves all settings from the database and returns them as a map where each key is the setting name
// and each value is a Setting struct. Returns an error if the query fails.
func LoadSettings(db *sql.DB) (map[string]Setting, error) {
	rows, err := db.Query("SELECT * FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]Setting)
	for rows.Next() {
		var s Setting
		if err := rows.Scan(&s.Key, &s.Value); err != nil {
			return nil, err
		}
		settings[s.Key] = s
	}
	return settings, nil
}

// Initializes the application configuration by loading settings from the database and merging them with default values.
// It updates the config object with the final values and returns the populated AppConfig struct or an error if the operation fails.
func LoadConfig() (*config.AppConfig, error) {
	dbCl, err := CreateConnection()
	if err != nil {
		return nil, err
	}
	defer dbCl.Close()

	settingsMap, err := LoadSettings(dbCl) // Load settings from database
	if err != nil {
		return nil, err
	}

	defaultDescIntervalMins, _ := strconv.Atoi(defaultSettings["DescGenIntervalMins"])
	defaultDescIntervalEnabled, _ := strconv.Atoi(defaultSettings["DescGenIntervalEnabled"])
	defaultScrIntervalMins, _ := strconv.Atoi(defaultSettings["ScreenshotIntervalMins"])
	defaultScrIntervalEnabled, _ := strconv.Atoi(defaultSettings["ScreenshotIntervalEnabled"])
	defaultReportAutoEnabled, _ := strconv.Atoi(defaultSettings["ReportAutoEnabled"])

	loadedConf := &config.AppConfig{
		ScrPath:                   defaultSettings["ScrPath"],
		DescGenAPI:                defaultSettings["DescGenAPI"],
		DescGenModel:              defaultSettings["DescGenModel"],
		DescGenPrompt:             defaultSettings["DescGenPrompt"],
		DescGenIntervalMins:       defaultDescIntervalMins,
		DescGenIntervalEnabled:    defaultDescIntervalEnabled,
		ScreenshotIntervalMins:    defaultScrIntervalMins,
		ScreenshotIntervalEnabled: defaultScrIntervalEnabled,
		ReportAPI:                 defaultSettings["ReportAPI"],
		ReportModel:               defaultSettings["ReportModel"],
		ReportAutoEnabled:         defaultReportAutoEnabled,
		ReportAutoAt:              defaultSettings["ReportAutoAt"],
		ReportPrompt:              defaultSettings["ReportPrompt"],
		OllamaURL:                 defaultSettings["OllamaURL"],
		GeminiAPIKey:              defaultSettings["GeminiAPIKey"],
	}

	// Update config with values from the database, validating them
	for key, setting := range settingsMap {
		switch key {
		case "ScrPath":
			loadedConf.ScrPath = setting.Value
		case "DescGenAPI":
			loadedConf.DescGenAPI = setting.Value
		case "DescGenModel":
			loadedConf.DescGenModel = setting.Value
		case "DescGenPrompt":
			loadedConf.DescGenPrompt = setting.Value
		case "DescGenIntervalMins":
			loadedConf.DescGenIntervalMins, _ = strconv.Atoi(setting.Value)
		case "DescGenIntervalEnabled":
			loadedConf.DescGenIntervalEnabled, _ = strconv.Atoi(setting.Value)
		case "ScreenshotIntervalMins":
			loadedConf.ScreenshotIntervalMins, _ = strconv.Atoi(setting.Value)
		case "ScreenshotIntervalEnabled":
			loadedConf.ScreenshotIntervalEnabled, _ = strconv.Atoi(setting.Value)
		case "ReportAPI":
			loadedConf.ReportAPI = setting.Value
		case "ReportModel":
			loadedConf.ReportModel = setting.Value
		case "ReportAutoEnabled":
			loadedConf.ReportAutoEnabled, _ = strconv.Atoi(setting.Value)
		case "ReportAutoAt":
			loadedConf.ReportAutoAt = setting.Value
		case "ReportPrompt":
			loadedConf.ReportPrompt = setting.Value
		case "OllamaURL":
			loadedConf.OllamaURL = setting.Value
		case "GeminiAPIKey":
			loadedConf.GeminiAPIKey = setting.Value
		}
	}

	config.Config = *loadedConf
	return loadedConf, nil
}

// Triggers re-initialization of specific app components (like scheduling or LLM) based on updated settings.
// It checks which settings have changed and only reinitializes components if required by those changes.
func RefreshInit(newSettings map[string]string) {
	_, ok1 := newSettings["ScreenshotIntervalMins"]
	_, ok2 := newSettings["DescGenIntervalMins"]

	if (ok1 || ok2) && Initializers.FunctionsGiven {
		Initializers.InitSchedule()
	}

	_, ok1 = newSettings["DescGenAPI"]
	_, ok2 = newSettings["ReportAPI"]

	if (ok1 || ok2) && Initializers.FunctionsGiven {
		Initializers.InitLLM()
	}
}

// Updates settings in both the database and the in-memory configuration (config.Config) using reflection.
// It ensures that changes to settings are saved persistently and reflected immediately in the running application.
func UpdateSettings(newSettings map[string]string) error {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
		return err
	}
	defer dbCl.Close()

	for key, val := range newSettings {
		err = updateSetting(dbCl, key, val)
		if err != nil {
			fmt.Printf("Error when updating %s with %s: %v\n", key, val, err.Error())
			return err
		}

		// Update the config.Config struct via reflection. Find field with 'key', then update the field with 'val'
		r := reflect.ValueOf(&config.Config).Elem()
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

// Updates a specific setting in the database using the provided key and new value.
// It performs a SQL UPDATE operation and returns an error if the update fails.
func updateSetting(db *sql.DB, key, newValue string) error {
	_, err := db.Exec("UPDATE settings SET value = ? WHERE key = ?", newValue, key)
	return err
}

// Goes through a map of keys and default values,
// checks if they exist in the settings table, and creates them if they don't exist.
func initializeSettings(db *sql.DB, defaultSettings map[string]string) error {
	// Prepare the SELECT statement
	selectStmt, err := db.Prepare("SELECT value FROM settings WHERE key = ?")
	if err != nil {
		return fmt.Errorf("error preparing SELECT statement: %w", err)
	}
	defer selectStmt.Close()

	// Prepare the INSERT statement
	insertStmt, err := db.Prepare("INSERT INTO settings (key, value) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing INSERT statement: %w", err)
	}
	defer insertStmt.Close()

	// Iterate through the default settings
	for key, defaultValue := range defaultSettings {
		// Check if the key exists
		var value string
		err := selectStmt.QueryRow(key).Scan(&value)

		if err == sql.ErrNoRows {
			// Key doesn't exist, insert the default value
			_, err = insertStmt.Exec(key, defaultValue)
			if err != nil {
				return fmt.Errorf("error inserting default value for key %s: %w", key, err)
			}
			log.Printf("Inserted default value for key: %s", key)
		} else if err != nil {
			// An error occurred during the SELECT
			return fmt.Errorf("error checking existence of key %s: %w", key, err)
		} else {
			// Key exists, do nothing
			log.Printf("Key already exists: %s", key)
		}
	}

	return nil
}

// Checks the settings in the provided map to ensure they are valid according to specific criteria.
// If a setting is found to be invalid, it updates that setting with the default value in the database.
//! There is currently no setting validation in place.
// func validateSettings(settings map[string]Setting) {
// 	dbCl, err := CreateConnection()
// 	if err != nil {
// 		fmt.Printf("Could not create DB connection: %v\n", err.Error())
// 	}

// 	for key, setting := range settings {
// 		if !isValid(setting) {
// 			err = updateSetting(dbCl, key, defaultSettings[key])
// 			if err != nil {
// 				fmt.Printf("Could not update setting: %v\n", err.Error())
// 			}
// 		}
// 	}
// }

// Checks if a specific setting meets the defined validation criteria for that setting.
// Returns true if the setting is valid, otherwise returns false.
// ! There is currently no setting validation.
// func isValid(setting Setting) bool {
// 	switch setting.Key {
// 	case "max_users":
// 		// Assume max_users should be a positive integer
// 		if _, err := strconv.Atoi(setting.Value); err != nil || setting.Value == "0" {
// 			return false
// 		}
// 	case "enable_feature":
// 		// Assume enable_feature should be "true" or "false"
// 		if setting.Value != "true" && setting.Value != "false" {
// 			return false
// 		}
// 	// Add more cases for other settings as needed
// 	default:
// 		return false // Unknown setting
// 	}
// 	return true
// }
