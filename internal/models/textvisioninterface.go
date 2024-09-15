package models

type ITextVisionModel interface {
	GenerateText(prompt string) (string, error)
	DescribeScreenshot(fileName string, prompt string) (string, error)
	DescribeBulkScreenshots(fileNames []string, prompt string) (string, error)
}
