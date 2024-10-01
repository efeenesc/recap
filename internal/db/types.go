package db

// Contains the full screenshot image in Base64 along with its description and other props
type CaptureScreenshotImage struct {
	CaptureID    int     `json:"CaptureID"`
	Timestamp    int64   `json:"Timestamp"`
	ScreenshotID int     `json:"ScreenshotID"`
	Filename     string  `json:"Filename"`
	Thumbname    *string `json:"Thumbname"`
	Description  *string `json:"Description"`
	Screenshot   string  `json:"Screenshot"`
}

// Contains description of screen capture along with other properties. Thumbname contains the thumbnail's filename
type CaptureScreenshot struct {
	CaptureID    int
	Timestamp    int64
	ScreenshotID int
	Filename     string
	Thumbname    *string
	Description  *string
}

// Basic properties of a screen capture
type CaptureDescription struct {
	CaptureID   int
	Timestamp   int64
	Description string
}

// Contains the report's content, ID, and UNIX second timestamp
type Report struct {
	ReportID  int    `json:"ReportID"`
	Timestamp int    `json:"Timestamp"`
	Content   string `json:"Content"`
}
