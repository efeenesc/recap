package schedule

import (
	"fmt"
	"log"
	"rcallport/internal/config"
	"rcallport/internal/db"
	"rcallport/internal/llm"
	"rcallport/internal/screenshot"
	"time"
)

var ScrTimer *time.Ticker
var LLMTimer *time.Ticker

func screenshotCallback() {
	fmt.Printf("Taking screenshot at %s\n", time.Now())
	fileName := screenshot.TakeScreenshot()
	cl, err := db.CreateConnection()
	if err != nil {
		log.Fatalf("Could not create database connection! %v\n", err.Error())
	}
	db.InsertCapture(cl, []string{fileName})
}

func startScreenshotTimer(minInterval int) {
	ScrTimer = time.NewTicker(time.Duration(minInterval) * time.Minute)
	go func() {
		for {
			<-ScrTimer.C
			screenshotCallback()
		}
	}()
}

func startLLMTimer(minInterval int) {
	llm.SendQueue()
	LLMTimer = time.NewTicker(time.Duration(minInterval) * time.Minute)
	go func() {
		for {
			<-LLMTimer.C
			fmt.Printf("Sending queued screenshots to LLM at %s\n", time.Now())
			llm.SendQueue()
		}
	}()
}

// Returns total minutes
func getFullIntervalTime(hrs *int, mins *int) int {
	var totalMins = 0

	if hrs != nil {
		totalMins += *hrs * 60
	}
	if mins != nil {
		totalMins += *mins
	}

	return totalMins
}

func Initialize() {
	ssGenObj := config.Config.LLM.Screenshot.DescriptionGenInterval
	ssTakeObj := config.Config.LLM.Screenshot.ScreenshotInterval

	totalMins := getFullIntervalTime(ssTakeObj.Hours, ssTakeObj.Minutes)
	if totalMins != 0 {
		startScreenshotTimer(totalMins)
	}

	if ssGenObj.Enabled {
		totalMins = getFullIntervalTime(ssGenObj.Hours, ssGenObj.Minutes)

		if totalMins != 0 {
			startLLMTimer(totalMins)
		}
	}
}
