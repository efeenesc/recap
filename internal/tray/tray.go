package tray

import (
	"fmt"
	"os"
	"recap/internal/app"
	"recap/internal/config"
	"recap/internal/llm"
	"recap/internal/schedule"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var ScrTrayBtn *systray.MenuItem
var LLMTrayBtn *systray.MenuItem

func getIcon() ([]byte, error) {
	return os.ReadFile("icon.ico")
}

func onReady() {
	systray.SetTitle("Recap")
	systray.SetTooltip("Open Recap")

	bytes, err := getIcon()
	if err == nil {
		systray.SetIcon(bytes)
	}

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
				llm.GenerateDailyReport()

			case <-exitBtn.ClickedCh:
				runtime.Quit(*app.WailsContext)
			}
		}
	}()
}

func onExit() {
	fmt.Println("Exiting program")
}

func Initialize() {
	systray.Register(onReady, onExit)
}
