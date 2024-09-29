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

func (t *Timer) stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.running = false
	}
}

func screenshotCallback() {
	fmt.Printf("Taking screenshot at %s\n", time.Now())
	fileName := screenshot.TakeScreenshot()
	cl, err := db.CreateConnection()
	if err != nil {
		log.Fatalf("Could not create database connection! %v\n", err.Error())
	}
	lastId := db.InsertCapture(cl, []string{fileName})
	app.AppInstance.SendScreenshotRanMessage(lastId)
}

func llmCallback() {
	fmt.Printf("Sending queued screenshots to LLM at %s\n", time.Now())
	llm.SendQueue()
}

func StartScreenshotSchedule(interval time.Duration) {
	fmt.Println("Starting screenshot schedule")
	screenshotTimer.start(interval, screenshotCallback)
}

func StartLLMTimer(interval time.Duration) {
	llmTimer.start(interval, llmCallback)
}

func SetLLMScheduleState(state bool) {
	if state {
		StartLLMTimer(getDefaultInterval(120))
	} else {
		llmTimer.stop()
	}
}

func SetScrScheduleState(state bool) {
	if state {
		StartScreenshotSchedule(getDefaultInterval(5))
	} else {
		screenshotTimer.stop()
	}
}

func AreTimersRunning() (bool, bool) {
	return screenshotTimer.running, llmTimer.running
}

func Initialize() {
	ssTakeObj := config.Config.LLM.Screenshot.ScreenshotInterval
	ssGenObj := config.Config.LLM.Screenshot.DescriptionGenInterval

	if interval := getDuration(ssTakeObj.Hours, ssTakeObj.Minutes); interval > 0 {
		StartScreenshotSchedule(interval)
	}

	if ssGenObj.Enabled {
		if interval := getDuration(ssGenObj.Hours, ssGenObj.Minutes); interval > 0 {
			StartLLMTimer(interval)
		}
	}
}

func getDuration(hours, minutes *int) time.Duration {
	var totalMinutes int
	if hours != nil {
		totalMinutes += *hours * 60
	}
	if minutes != nil {
		totalMinutes += *minutes
	}
	return time.Duration(totalMinutes) * time.Minute
}

func getDefaultInterval(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}
