package app

import (
	"fmt"
	"rcallport/internal/db"

	"github.com/sqweek/dialog"
)

type TimerState struct {
	Scr bool `json:"scr"`
	Llm bool `json:"llm"`
}

func (a *AppMethods) EmitStartStopScrTimer(state bool) {
	if a.FunctionsGiven {
		a.CSetScrTimer(state)
	}
}

func (a *AppMethods) EmitStartStopLLMTimer(state bool) {
	if a.FunctionsGiven {
		a.CSetLLMTimer(state)
	}
}

func (a *AppMethods) CheckTimers() *TimerState {
	if a.FunctionsGiven {
		scr, llm := a.CCheckTimers()
		tmr := &TimerState{
			Scr: scr,
			Llm: llm,
		}

		return tmr
	}

	return &TimerState{}
}

func (a *AppMethods) GetScreenshots(limit int) []db.CaptureScreenshotImage {
	if a.FunctionsGiven {
		results, err := a.CGetScreenshots(limit)
		if err != nil {
			fmt.Println("FUUUCK")
			fmt.Println(err)
			return []db.CaptureScreenshotImage{}
		}

		return results
	}

	return []db.CaptureScreenshotImage{}
}

func (a *AppMethods) GetReports(limit int) []db.Report {
	if a.FunctionsGiven {
		results, err := a.CGetReports(limit)
		if err != nil {
			fmt.Println(err)
			return []db.Report{}
		}

		return results
	}

	return []db.Report{}
}

func (a *AppMethods) GetReportById(id int) *db.Report {
	if a.FunctionsGiven {
		result, err := a.CGetReportById(id)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return result
	}

	return nil
}

func (a *AppMethods) GetScreenshotById(id int) *db.CaptureScreenshotImage {
	if a.FunctionsGiven {
		result, err := a.CGetScreenshotById(id)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return result
	}

	return nil
}

func (a *AppMethods) GenerateReportFromScreenshotIds(ids []int) (*int64, error) {
	if a.FunctionsGiven {
		result, err := a.CGenerateReportWithSelectScr(ids)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}

		return result, nil
	}

	return nil, nil
}

func (a *AppMethods) GetConfig() *db.AppConfig {
	if a.FunctionsGiven {
		result, err := a.CGetConfig()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return result
	}

	return nil
}

func (a *AppMethods) GetDisplayValues() map[string]db.SettingDisplayProps {
	if a.FunctionsGiven {
		result := a.CGetDisplayValues()
		return result
	}

	return nil
}

func (a *AppMethods) SelectFolder() (*string, error) {
	var folder string
	var err error

	folder, err = dialog.Directory().Title("Select new screenshot folder").Browse()
	if err == dialog.ErrCancelled {
		folder = "Cancelled"
	}

	return &folder, err
}
