package db

type CaptureScreenshotImage struct {
	CaptureID    int     `json:"CaptureID"`
	Timestamp    int64   `json:"Timestamp"`
	ScreenshotID int     `json:"ScreenshotID"`
	Filename     string  `json:"Filename"`
	Description  *string `json:"Description"`
	Screenshot   string  `json:"Screenshot"`
}

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

type Report struct {
	ReportID  int    `json:"ReportID"`
	Timestamp int    `json:"Timestamp"`
	Content   string `json:"Content"`
}
