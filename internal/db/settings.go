package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type AppConfig struct {
	ScrPath                   string `json:"ScrPath"`
	DescGenAPI                string `json:"DescGenAPI"`
	DescGenModel              string `json:"DescGenModel"`
	DescGenPrompt             string `json:"DescGenPrompt"`
	DescGenIntervalMins       int    `json:"DescGenIntervalMins"`
	DescGenIntervalEnabled    int    `json:"DescGenIntervalEnabled"`
	ScreenshotIntervalMins    int    `json:"ScreenshotIntervalMins"`
	ScreenshotIntervalEnabled int    `json:"ScreenshotIntervalEnabled"`
	ReportAPI                 string `json:"ReportAPI"`
	ReportModel               string `json:"ReportModel"`
	ReportAutoEnabled         int    `json:"ReportAutoEnabled"`
	ReportAutoAt              string `json:"ReportAutoAt"`
	ReportPrompt              string `json:"ReportPrompt"`
	OllamaURL                 string `json:"OllamaURL"`
	GeminiAPIKey              string `json:"GeminiAPIKey"`
}

type Setting struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

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
	"ReportAutoAt":              "09:00", // Default auto report time
	"ReportPrompt":              "You are an AI assistant tasked with generating a daily activity report for a user based on a series of visual descriptions captured from their computer screen throughout the day. Your job is to summarize this data into brief items describing what the user worked on today.",
	"OllamaURL":                 "http://localhost:11434",
	"GeminiAPIKey":              "your-gemini-api-key",
}

func GetDisplayValues() map[string]SettingDisplayProps {
	var settingKeyDisplayVals = map[string]SettingDisplayProps{
		"ScrPath":                   {DisplayName: "Path", Description: "Where screenshots will be saved", Category: "Screenshots", InputType: "FolderPicker"},
		"DescGenAPI":                {DisplayName: "API", Description: "API to be used while generating screenshot descriptions", Category: "Vision", InputType: "APIPicker"},
		"DescGenModel":              {DisplayName: "Model", Description: "Vision model to be used while generating screenshot descriptions", Category: "Vision", InputType: "APIModelPicker"},
		"DescGenPrompt":             {DisplayName: "Prompt", Description: "Prompt to be used while generating screenshot descriptions", Category: "Vision", InputType: "ExtendedTextInput"},
		"DescGenIntervalMins":       {DisplayName: "Interval", Description: "Minutes interval of sending screenshots for description generation", Category: "Vision", InputType: "NumberInput"},
		"DescGenIntervalEnabled":    {DisplayName: "Schedule", Description: "Whether descriptions should be automatically sent for description generation", Category: "Vision", InputType: "Boolean"}, // 1 for enabled, 0 for disabled
		"ScreenshotIntervalMins":    {DisplayName: "Interval", Description: "Minutes interval of taking screenshots", Category: "Screenshots", InputType: "NumberInput"},                              // Default interval in minutes
		"ScreenshotIntervalEnabled": {DisplayName: "Schedule", Description: "Whether screenshots should be automatically taken", Category: "Screenshots", InputType: "Boolean"},                       // 1 for enabled, 0 for disabled
		"ReportAPI":                 {DisplayName: "API", Description: "API to be used while generating reports from screenshot descriptions", Category: "Reports", InputType: "APIPicker"},
		"ReportModel":               {DisplayName: "Model", Description: "Text generation model to be used while generating reports from screenshot descriptions", Category: "Reports", InputType: "APIModelPicker"},
		"ReportAutoEnabled":         {DisplayName: "Schedule", Description: "Whether reports should be automatically generated", Category: "Reports", InputType: "Boolean"},
		"ReportAutoAt":              {DisplayName: "Time", Description: "Time at which a report should be generated every day", Category: "Reports", InputType: "TimePicker"}, // Default auto report time
		"ReportPrompt":              {DisplayName: "Prompt", Description: "Prompt to be used when generating a report from screenshot descriptions", Category: "Reports", InputType: "ExtendedTextInput"},
		"OllamaURL":                 {DisplayName: "Ollama URL", Description: "URL of Ollama's endpoint (port included). Defaults to http://localhost:11434", Category: "Models", InputType: "URLInput"},
		"GeminiAPIKey":              {DisplayName: "Gemini API key", Description: "Your Gemini API key. You can get your own free API key from Google. Refer to the \"Setting up\" section of the tutorial for instructions", Category: "Models", InputType: "TextInput"},
	}

	return settingKeyDisplayVals
}

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

func LoadConfig() (*AppConfig, error) {
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

	config := &AppConfig{
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
			config.ScrPath = setting.Value
		case "DescGenAPI":
			config.DescGenAPI = setting.Value
		case "DescGenModel":
			config.DescGenModel = setting.Value
		case "DescGenPrompt":
			config.DescGenPrompt = setting.Value
		case "DescGenIntervalMins":
			config.DescGenIntervalMins, _ = strconv.Atoi(setting.Value)
		case "DescGenIntervalEnabled":
			config.DescGenIntervalEnabled, _ = strconv.Atoi(setting.Value)
		case "ScreenshotIntervalMins":
			config.ScreenshotIntervalMins, _ = strconv.Atoi(setting.Value)
		case "ScreenshotIntervalEnabled":
			config.ScreenshotIntervalEnabled, _ = strconv.Atoi(setting.Value)
		case "ReportAPI":
			config.ReportAPI = setting.Value
		case "ReportModel":
			config.ReportModel = setting.Value
		case "ReportAutoEnabled":
			config.ReportAutoEnabled, _ = strconv.Atoi(setting.Value)
		case "ReportAutoAt":
			config.ReportAutoAt = setting.Value
		case "ReportPrompt":
			config.ReportPrompt = setting.Value
		case "OllamaURL":
			config.OllamaURL = setting.Value
		case "GeminiAPIKey":
			config.GeminiAPIKey = setting.Value
		}
	}
	return config, nil
}

func validateSettings(settings map[string]Setting) {
	dbCl, err := CreateConnection()
	if err != nil {
		fmt.Printf("Could not create DB connection: %v\n", err.Error())
	}

	for key, setting := range settings {
		if !isValid(setting) {
			updateSetting(dbCl, key, defaultSettings[key])
		}
	}
}

func updateSetting(db *sql.DB, key, newValue string) error {
	_, err := db.Exec("UPDATE settings SET value = ? WHERE key = ?", newValue, key)
	return err
}

func isValid(setting Setting) bool {
	switch setting.Key {
	case "max_users":
		// Assume max_users should be a positive integer
		if _, err := strconv.Atoi(setting.Value); err != nil || setting.Value == "0" {
			return false
		}
	case "enable_feature":
		// Assume enable_feature should be "true" or "false"
		if setting.Value != "true" && setting.Value != "false" {
			return false
		}
	// Add more cases for other settings as needed
	default:
		return false // Unknown setting
	}
	return true
}
