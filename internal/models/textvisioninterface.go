package models

type ITextVisionModel interface {
	// Generate text with a text prompt. Returns the response, or an error if one is received
	GenerateText(prompt string) (string, error)

	// Describes screenshot, sending the screenshot file along with a text prompt.
	// The screenshot is loaded by combining ScrPath with fileName. Returns the response, or
	// an error if one is received
	DescribeScreenshot(fileName string, prompt string) (string, error)

	// Helper function to describe screenshots in bulk. Works almost the same as DescribeScreenshot under the hood
	DescribeBulkScreenshots(fileNames []string, prompt string) (string, error)
}
