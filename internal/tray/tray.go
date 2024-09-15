package tray

import (
	"fmt"
	"rcallport/internal/llm"

	"github.com/getlantern/systray"
)

func onReady() {
	// systray.SetIcon(icon.Data)
	systray.SetTitle("Rcallport")
	systray.SetTooltip("Rcallport")

	summaryBtn := systray.AddMenuItem("Generate report", "Get a report of today's descriptions")
	exitBtn := systray.AddMenuItem("Exit", "Close the application")

	go func() {
		for {
			select {
			case <-exitBtn.ClickedCh: // Wait for a click
				// Run your function here
				fmt.Println("Exit button clicked!")
				// Perform cleanup or exit operations
				systray.Quit()
			case <-summaryBtn.ClickedCh:
				fmt.Println("Summary button clicked!")
				llm.GenerateDailyReport()
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exitting program")
}

func Initialize() {
	systray.Run(onReady, onExit)
}
