package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) SendMessage(msg *string) {
	if msg == nil {
		defaultMsg := ""
		msg = &defaultMsg
	}
	runtime.EventsEmit(*WailsContext, *msg)
}

func (a *App) SendScreenshotRanMessage(lastId int64) {
	runtime.EventsEmit(*WailsContext, "rcv:screenshotran", lastId)
}

func (a *App) SendLLMRanMessage(lastId int64) {
	runtime.EventsEmit(*WailsContext, "rcv:llmran", lastId)
}

func (a *App) SendScreenshotStateMessage(newState bool) {
	runtime.EventsEmit(*WailsContext, "rcv:screenshotstate", newState)
}

func (a *App) SendLLMStateMessage(newState bool) {
	runtime.EventsEmit(*WailsContext, "rcv:llmstate", newState)
}
