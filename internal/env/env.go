package env

import (
	"fmt"
	"os"
	"path/filepath"
	"rcallport/internal/config"
	"strings"
)

func findEnvFiles(root string) ([]string, error) {
	var envFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == ".env" {
			envFiles = append(envFiles, path)
		}
		return nil
	})
	return envFiles, err
}

func loadEnvFile(path string) error {
	filebinary, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load .env file %s: %w", path, err)
	}

	if len(filebinary) == 0 {
		return nil
	}

	filestring := string(filebinary)
	lines := strings.Split(filestring, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])
		os.Setenv(key, value)
	}
	return nil
}

func Initialize() (status int) {
	projectRoot, _ := config.GetProjectRoot()

	envFiles, err := findEnvFiles(projectRoot)
	if err != nil {
		fmt.Println("Error searching for .env files:", err)
		return 0
	}

	if len(envFiles) == 0 {
		fmt.Println("Could not find any .env files in the project directory.")
		return 0
	}

	for _, envFile := range envFiles {
		err := loadEnvFile(envFile)
		if err != nil {
			fmt.Printf("Error loading .env file %s: %v\n", envFile, err)
			return 0
		}
		fmt.Printf("Loaded .env file: %s\n", envFile)
	}

	return 1
}
