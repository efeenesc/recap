package config

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var Config AppConfig
var Info AppInfo

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

type AppInfo struct {
	Version                string `json:"Version"`
	FirstTimeTutorialShown string `json:"FirstTimeTutorialShown"`
}

// CreateFolderIfNotExists checks if a folder exists at the given path, and if not, creates it with permissions set to 0700.
// If the folder creation fails, it logs a fatal error and terminates the program.
//
// Parameters:
//   - path: The file path where the folder should be checked and potentially created.
func CreateFolderIfNotExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0700)
		log.Fatalf("Could not create new folder: %v\n", err.Error())
	}
}

// Initialize attempts to load the configuration from predefined paths and
// unmarshals it into the Config struct. It searches for the configuration
// file in the following locations relative to the project root:
// - ./configs/config.yaml
// - ./config.yaml
//
// If the configuration file is found, it reads the file and unmarshals its
// contents into the Config struct. Additionally, it ensures that the
// directory specified by Config.ScrPath exists, creating it if necessary.
//
// If no configuration file is found, the function panics.
//
// NOTE: This function is currently unused as configurations are now stored
// in the database instead of YAML files.
func Initialize() {
	proot := GetProjectRoot()

	configPaths := []string{path.Join(proot, "./configs/config.yaml"), path.Join(proot, "./config.yaml")}
	found := false

	for _, p := range configPaths {
		_, err := os.Stat(p)
		if err != nil {
			continue
		}
		found = true

		data, err := os.ReadFile(p)
		if err != nil {
			log.Fatalf("Failed to load config.yaml: %v", err.Error())
		}

		err = yaml.Unmarshal([]byte(data), &Config)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		CreateFolderIfNotExists(path.Join(proot, Config.ScrPath))
		break
	}

	if !found {
		panic("config.yaml not found")
	}
}
