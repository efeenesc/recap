package config

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var Config AppConfig

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

func createFolderIfNotExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0700)
		log.Fatalf("Could not create new folder: %v\n", err.Error())
	}
}

func Initialize() {
	proot, _ := GetProjectRoot()

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

		createFolderIfNotExists(path.Join(proot, Config.ScrPath))
		break
	}

	if !found {
		panic("config.yaml not found")
	}
}
