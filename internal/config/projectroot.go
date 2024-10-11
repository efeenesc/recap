package config

import (
	"os"
	"path"
)

// Returns the current working directory.
func GetProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return currentDir
}

func RelativeToAbsPath(relativePath string) string {
	return path.Join(GetProjectRoot(), relativePath)
}
