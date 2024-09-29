package utils

import (
	"encoding/base64"
	"log"
	"os"
	"path"
	"rcallport/internal/config"
)

func ReadImageToBase64(fileName string) string {
	proot, _ := config.GetProjectRoot()
	fullPath := path.Join(proot, config.Config.System.ScreenshotPath, fileName)

	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Error reading image file: %v", err)
		return ""
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
}
