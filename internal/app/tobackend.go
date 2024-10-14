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
	if a.CSetScrTimer != nil {
		a.CSetScrTimer(state)
	}
}

func (a *AppMethods) EmitStartStopLLMTimer(state bool) {
	if a.CSetLLMTimer != nil {
		a.CSetLLMTimer(state)
	}
}

func (a *AppMethods) CheckTimers() *TimerState {
	if a.CCheckTimers != nil {
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
	if a.CGetScreenshots != nil {
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
	if a.CGetScreenshotsOlderThan != nil {
		results, err := a.CGetScreenshotsOlderThan(id, limit)
		if err != nil {
			fmt.Printf("Received error from GetScreenshotsOlderThan: %v\n", err)
			return []db.CaptureScreenshotImage{}, err
		}

		return results, nil
	}

	return []db.CaptureScreenshotImage{}, fmt.Errorf("callback functions were not passed to AppMethods")
}

func (a *AppMethods) GetScreenshotsNewerThan(id int) ([]db.CaptureScreenshotImage, error) {
	if a.CGetScreenshotsNewerThan != nil {
		results, err := a.CGetScreenshotsNewerThan(id)
		if err != nil {
			fmt.Printf("Received error from GetScreenshotsNewerThan: %v\n", err)
			return []db.CaptureScreenshotImage{}, err
		}

		return results, nil
	}

	return []db.CaptureScreenshotImage{}, fmt.Errorf("callback functions were not passed to AppMethods")
}

func (a *AppMethods) GetReports(limit int) []db.Report {
	if a.CGetReports != nil {
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
	if a.CGetReportById != nil {
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
	if a.CGetReportsNewerThan != nil {
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
	if a.CGetReportsOlderThan != nil {
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
	if a.CGetScreenshotById != nil {
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
	if a.CDeleteScreenshotsById != nil {
		err := a.CDeleteScreenshotsById(ids)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	return fmt.Errorf("callback functions not passed")
}

func (a *AppMethods) DeleteReportsById(ids []int) error {
	if a.CDeleteReportsById != nil {
		err := a.CDeleteReportsById(ids)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}

	return fmt.Errorf("callback functions not passed")
}

func (a *AppMethods) GenerateReportFromScreenshotIds(ids []int) (*int64, error) {
	if a.CGenerateReportWithSelectScr != nil {
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
	if a.CGetConfig != nil {
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
	if a.CGetDisplayValues != nil {
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
	if a.CUpdateSettings != nil {
		err := a.CUpdateSettings(newSettings)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("missing function UpdateSettings")
}

func (a *AppMethods) UpdateInfo(newInfo map[string]string) error {
	if a.CUpdateInfo != nil {
		err := a.CUpdateInfo(newInfo)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("missing function UpdateInfo")
}

func (a *AppMethods) WriteInfo(key, value string) error {
	if a.CWriteInfo != nil {
		err := a.CWriteInfo(key, value)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("missing function WriteInfo")
}

func (a *AppMethods) ReadInfo(key string) (*db.Info, error) {
	if a.CReadInfo != nil {
		val, err := a.CReadInfo(key)
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return nil, fmt.Errorf("missing function WriteInfo")
}

func (a *AppMethods) ReadAllInfo() (map[string]string, error) {
	if a.CReadAllInfo != nil {
		val, err := a.CReadAllInfo()
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return nil, fmt.Errorf("missing function WriteInfo")
}
