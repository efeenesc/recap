package screenshot

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path"

	"rcallport/internal/config"

	"github.com/google/uuid"
	"github.com/kbinani/screenshot"
)

func saveScreenshot(img *image.RGBA) (string, error) {
	fileName := fmt.Sprintf("%s.png", uuid.New())

	proot, _ := config.GetProjectRoot()
	file, err := os.Create(path.Join(proot, config.Config.System.ScreenshotPath, fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func TakeScreenshot() string {
	var displayBounds image.Rectangle
	displayBounds.Min.X = 0
	displayBounds.Min.Y = 0

	for idx := range screenshot.NumActiveDisplays() {
		bounds := screenshot.GetDisplayBounds(idx)

		displayBounds.Max.X = max(bounds.Max.X, displayBounds.Max.X)
		displayBounds.Max.Y = max(bounds.Max.Y, displayBounds.Max.Y)
	}

	// fmt.Printf("%d : %d\n", displayBounds.Dx(), displayBounds.Dy())
	img, err := screenshot.CaptureRect(displayBounds)
	if err != nil {
		panic(err)
	}

	fileName, err := saveScreenshot(img)
	if err != nil {
		log.Fatal(err)
	}
	return fileName
}
