package config

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var Config AppConfig

type AppConfig struct {
	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	System struct {
		ReportPath     string `yaml:"report_path"`
		ScreenshotPath string `yaml:"screenshot_path"`
	} `yaml:"system"`
	LLM struct {
		Screenshot struct {
			API struct {
				Connector string `yaml:"connector"`
				Model     string `yaml:"model"`
			} `yaml:"api"`
			DescriptionGenInterval struct {
				Enabled bool `yaml:"enabled"`
				Hours   *int `yaml:"hours"`
				Minutes *int `yaml:"minutes"`
			} `yaml:"description_gen_interval"`
			ScreenshotInterval struct {
				Hours   *int `yaml:"hours"`
				Minutes *int `yaml:"minutes"`
			} `yaml:"screenshot_interval"`
			Prompt string `yaml:"prompt"`
		} `yaml:"screenshot"`
		Report struct {
			API struct {
				Connector string `yaml:"connector"`
				Model     string `yaml:"model"`
			} `yaml:"api"`
			Prompt string `yaml:"prompt"`
		} `yaml:"report"`
	} `yaml:"llm"`
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
		break
	}

	if !found {
		panic("config.yaml not found")
	}
}
