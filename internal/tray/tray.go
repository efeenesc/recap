package tray

import (
	"fmt"
	"rcallport/internal/app"
	"rcallport/internal/config"
	"rcallport/internal/llm"
	"rcallport/internal/schedule"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var ScrTrayBtn *systray.MenuItem
var LLMTrayBtn *systray.MenuItem

func onReady() {
	// systray.SetIcon(icon.Data)
	systray.SetTitle("Rcallport")
	systray.SetTooltip("Open Rcallport")

	scrScheduleEnabled := config.Config.ScreenshotIntervalEnabled == 1
	llmScheduleEnabled := config.Config.DescGenIntervalEnabled == 1

	openBtn := systray.AddMenuItem("Open", "Show Rcallport's window")
	systray.AddSeparator()
	ScrTrayBtn := systray.AddMenuItemCheckbox("Automatic screenshots", "Turn automatic screenshots on or off", scrScheduleEnabled)
	LLMTrayBtn := systray.AddMenuItemCheckbox("Description generation", "Turn automatic screenshot description generation on or off", llmScheduleEnabled)
	summaryBtn := systray.AddMenuItem("Generate report", "Get a report of today's descriptions")
	systray.AddSeparator()
	exitBtn := systray.AddMenuItem("Exit", "Close the application")

	go func() {
		for {
			select {
			case <-openBtn.ClickedCh:
				runtime.Show(*app.WailsContext)

			case <-ScrTrayBtn.ClickedCh:
				if ScrTrayBtn.Checked() {
					ScrTrayBtn.Uncheck()
					schedule.SetScrScheduleState(false)
				} else {
					ScrTrayBtn.Check()
					schedule.SetScrScheduleState(true)
				}

			case <-LLMTrayBtn.ClickedCh:
				if LLMTrayBtn.Checked() {
					LLMTrayBtn.Uncheck()
					schedule.SetLLMScheduleState(false)
				} else {
					LLMTrayBtn.Check()
					schedule.SetLLMScheduleState(true)
				}

			case <-summaryBtn.ClickedCh:
				fmt.Println("Summary button clicked!")
				llm.GenerateDailyReport()

			case <-exitBtn.ClickedCh:
				fmt.Println("Exit button clicked!")
				runtime.Quit(*app.WailsContext)
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
