package schedule

import (
	"fmt"
	"log"
	"rcallport/internal/app"
	"rcallport/internal/config"
	"rcallport/internal/db"
	"rcallport/internal/llm"
	"rcallport/internal/screenshot"
	"time"
)

type Timer struct {
	ticker  *time.Ticker
	running bool
}

var (
	screenshotTimer = &Timer{}
	llmTimer        = &Timer{}
)

// Initializes and starts the timer with the specified interval and callback function.
// It stops any existing timer before starting a new one
func (t *Timer) start(interval time.Duration, callback func()) {
	t.stop()
	t.ticker = time.NewTicker(interval)
	t.running = true
	go func() {
		for {
			<-t.ticker.C
			if !t.running {
				return
			}
			callback()
		}
	}()
}

// Halts the timer if it is currently running and stops the ticker
func (t *Timer) stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.running = false
	}
}

// Callback function; captures a screenshot, saves it, writes it to the database
func screenshotCallback() {
	fmt.Printf("Taking screenshot at %s\n", time.Now())
	fullFilename, thumbFilename := screenshot.TakeScreenshot()
	cl, err := db.CreateConnection()
	if err != nil {
		log.Fatalf("Could not create database connection! %v\n", err.Error())
	}
	lastId := db.InsertCapture(cl, []db.FullThumbScrPair{{Full: fullFilename, Thumb: thumbFilename}})
	app.AppInstance.SendScreenshotRanMessage(lastId)
}

// Sends unprocessed screenshots to the vision model and inserts the descriptions
// in the database
func llmCallback() {
	fmt.Printf("Sending queued screenshots to LLM at %s\n", time.Now())
	llm.SendQueue()
}

// Initiates the screenshot capturing process at the specified interval.
// It stops any currently running screenshot timer before starting a new one
func StartScreenshotSchedule(interval time.Duration) {
	fmt.Println("Starting screenshot schedule")
	if screenshotTimer.running {
		screenshotTimer.ticker.Stop()
	}
	screenshotTimer.start(interval, screenshotCallback)
}

// Initiates the LLM timer process at the specified interval.
// It stops any currently running LLM timer before starting a new one
func StartLLMTimer(interval time.Duration) {
	if llmTimer.running {
		llmTimer.ticker.Stop()
	}
	llmTimer.start(interval, llmCallback)
}

// Enables or disables the LLM schedule based on the provided state.
// If enabled, it starts the LLM timer based on configuration; stops the timer otherwise
func SetLLMScheduleState(state bool) {
	if state {
		StartLLMTimer(time.Minute * time.Duration(config.Config.DescGenIntervalMins))
	} else {
		llmTimer.stop()
	}

	app.AppInstance.SendLLMStateMessage(state)
}

// Enables or disables the screenshot schedule based on the provided state.
// If enabled, it starts the screenshot timer based on configuration; stops the timer otherwise
func SetScrScheduleState(state bool) {
	if state {
		StartScreenshotSchedule(time.Minute * time.Duration(config.Config.ScreenshotIntervalMins))
	} else {
		screenshotTimer.stop()
	}

	app.AppInstance.SendScreenshotStateMessage(state)
}

// Returns the running state of the screenshot and LLM timers. Used by frontend
func AreTimersRunning() (bool, bool) {
	return screenshotTimer.running, llmTimer.running
}

// Sets up the timers based on configuration settings for screenshot capturing
// and LLM generation. It starts the timers if enabled in the config.
func Initialize() {
	ssTakeEnabled := config.Config.ScreenshotIntervalEnabled
	descGenEnabled := config.Config.DescGenIntervalEnabled
	ssTakeInterval := config.Config.ScreenshotIntervalMins
	descGenInterval := config.Config.DescGenIntervalMins

	if ssTakeEnabled == 1 && ssTakeInterval > 0 {
		StartScreenshotSchedule(time.Duration(ssTakeInterval) * time.Minute)
	}

	if descGenEnabled == 1 && descGenInterval > 0 {
		StartLLMTimer(time.Duration(descGenInterval) * time.Minute)
	}
}
