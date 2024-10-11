package tray

import (
	"fmt"
	"recap/internal/app"
	"recap/internal/config"
	"recap/internal/llm"
	"recap/internal/schedule"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var iconBytes *[]byte
var ScrTrayBtn *systray.MenuItem
var LLMTrayBtn *systray.MenuItem

// Callback for when the systray is ready. It sets up the systray's title, tooltip, icon, menu items
// and a separator
func onReady() {
	systray.SetTitle("Recap")
	systray.SetTooltip("Open Recap")

	systray.SetIcon(*iconBytes)

	scrScheduleEnabled := config.Config.ScreenshotIntervalEnabled == 1
	llmScheduleEnabled := config.Config.DescGenIntervalEnabled == 1

	openBtn := systray.AddMenuItem("Open", "Show Recap")
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
				_, _ = llm.GenerateDailyReport()

			case <-exitBtn.ClickedCh:
				runtime.Quit(*app.WailsContext)
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exiting program")
}

// Initialize sets up the systray with the given icon and registers two callback functions.
// The first callback, onReady, is called when the systray is ready and sets up the systray's title, tooltip, and icon,
// and adds three menu items and a separator. The second callback, onExit, is called when the application exits and does
// nothing.
func Initialize(iconptr *[]byte) {
	iconBytes = iconptr
	systray.Register(onReady, onExit)
}
