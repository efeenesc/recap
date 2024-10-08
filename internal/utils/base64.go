package utils

import (
	"encoding/base64"
	"log"
	"os"
	"path"
	"recap/internal/config"
	"strings"
)

func ReadImageToBase64(fileName string) string {
	fullPath := path.Join(config.Config.ScrPath, fileName)

	bytes, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Error reading image file: %v", err)
		return ""
	}

	if strings.HasSuffix(fileName, ".png") {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
	} else {
		return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(bytes)
	}
}

func readImageBytes(filepath string) (*[]byte, error) {
	bytes, err := os.ReadFile(filepath)

	if err != nil {
		log.Printf("Error reading image file: %v", err)
		return nil, err
	}

	return &bytes, err
}

func formatResponse(bytes *[]byte, filename string) string {
	if strings.HasSuffix(filename, ".png") {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(*bytes)
	} else {
		return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(*bytes)
	}
}

// Read the thumbnail. If it doesn't exist, try loading full image instead
func ReadImageToBase64PreferThumb(fullName string, thumbName *string) string {
	if thumbName != nil {
		thumbPath := path.Join(config.Config.ScrPath, *thumbName)
		// Attempt to load thumbnail image and return
		bytes, err := readImageBytes(thumbPath)
		if err == nil {
			return formatResponse(bytes, thumbPath)
		}
	}

	fullPath := path.Join(config.Config.ScrPath, fullName)
	bytes, err := readImageBytes(fullPath)
	if err == nil {
		return formatResponse(bytes, fullPath)
	}

	return ""
}

// Read the full image. If it doesn't exist, try loading thumbnail instead
func ReadImageToBase64PreferFull(fullName string, thumbName *string) string {
	fullPath := path.Join(config.Config.ScrPath, fullName)

	// Attempt to load full image and return
	bytes, err := readImageBytes(fullPath)
	if err == nil {
		return formatResponse(bytes, fullPath)
	}

	if thumbName != nil {
		thumbPath := path.Join(config.Config.ScrPath, *thumbName)

		bytes, err = readImageBytes(thumbPath)
		if err == nil {
			return formatResponse(bytes, thumbPath)
		}
	}

	return ""
}
