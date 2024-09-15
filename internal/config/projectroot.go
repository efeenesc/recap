package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Navigate up to the project root
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "configs")); err == nil {
			return currentDir, nil
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("project root not found")
		}
		currentDir = parent
	}
}
