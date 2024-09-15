package db

type CaptureScreenshot struct {
	CaptureID    int
	Timestamp    int64
	ScreenshotID int
	Filename     string
	Description  *string
}

type CaptureDescription struct {
	CaptureID   int
	Timestamp   int64
	Description string
}
