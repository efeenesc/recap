package screenshot

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"

	"recap/internal/config"

	"github.com/google/uuid"
	"github.com/kbinani/screenshot"
)

// Saves the provided RGBA image as a JPEG file in the specified directory. Called by TakeScreenshot.
// Returns error if the operation fails
func saveScreenshotJPEG(img *image.RGBA, filename string, quality int) error {
	file, err := os.Create(path.Join(config.Config.ScrPath, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return err
	}

	return nil
}

// Saves the provided RGBA image as a PNG file in the specified directory. Called by TakeScreenshot.
// Returns error if the operation fails
func saveScreenshotPNG(img *image.RGBA, filename string) error {
	file, err := os.Create(path.Join(config.Config.ScrPath, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	enc := png.Encoder{
		CompressionLevel: png.BestCompression,
	}

	err = enc.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// Captures the entire screen based on the bounds of active displays. The end result is one image showing all screens.
// It saves the captured image using the saveScreenshot function once as PNG with best compression, and once as JPEG with 40% quality for thumbnails.
// Returns PNG filename first, the JPEG thumbnail filename second.
func TakeScreenshot() (string, string) {
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

	scrUuid := uuid.New()
	fullFilename := fmt.Sprintf("%s.png", scrUuid)
	thumbFilename := fmt.Sprintf("%s_thumb.jpg", scrUuid)

	err = saveScreenshotPNG(img, fullFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = saveScreenshotJPEG(img, thumbFilename, 40)

	return fullFilename, thumbFilename
}
