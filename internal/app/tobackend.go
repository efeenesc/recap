package app

import (
	"fmt"
	"recap/internal/config"
	"recap/internal/db"

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
			fmt.Println(err)
			return []db.CaptureScreenshotImage{}
		}

		return results
	}

	return []db.CaptureScreenshotImage{}
}

func (a *AppMethods) GetScreenshotsOlderThan(id int, limit int) ([]db.CaptureScreenshotImage, error) {
	if a.FunctionsGiven {
		results, err := a.CGetScreenshotsOlderThan(id, limit)
		if err != nil {
			fmt.Printf("Received error from GetScreenshotsOlderThan: %v\n", err)
			return []db.CaptureScreenshotImage{}, err
		}

		return results, nil
	}

	return []db.CaptureScreenshotImage{}, fmt.Errorf("Callback functions were not passed to AppMethods\n")
}

func (a *AppMethods) GetScreenshotsNewerThan(id int) ([]db.CaptureScreenshotImage, error) {
	if a.FunctionsGiven {
		results, err := a.CGetScreenshotsNewerThan(id)
		if err != nil {
			fmt.Printf("Received error from GetScreenshotsNewerThan: %v\n", err)
			return []db.CaptureScreenshotImage{}, err
		}

		return results, nil
	}

	return []db.CaptureScreenshotImage{}, fmt.Errorf("Callback functions were not passed to AppMethods\n")
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

func (a *AppMethods) GetReportsNewerThan(id int) []db.Report {
	if a.FunctionsGiven {
		result, err := a.CGetReportsNewerThan(id)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		return result
	}

	return nil
}

func (a *AppMethods) GetReportsOlderThan(id int, limit int) []db.Report {
	if a.FunctionsGiven {
		result, err := a.CGetReportsOlderThan(id, limit)
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

func (a *AppMethods) DeleteScreenshotsById(ids []int) error {
	if a.FunctionsGiven {
		err := a.CDeleteScreenshotsById(ids)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	return fmt.Errorf("Callback functions not passed")
}

func (a *AppMethods) DeleteReportsById(ids []int) error {
	if a.FunctionsGiven {
		err := a.CDeleteReportsById(ids)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	return fmt.Errorf("Callback functions not passed")
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

func (a *AppMethods) GetConfig() *config.AppConfig {
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

func (a *AppMethods) UpdateSettings(newSettings map[string]string) error {
	if a.FunctionsGiven {
		err := a.CUpdateSettings(newSettings)
		if err != nil {
			return err
		}

		return nil
	} else {
		return fmt.Errorf("Missing function when calling UpdateSettings\n")
	}
}
