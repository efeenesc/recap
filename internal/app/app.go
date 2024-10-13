package app

import (
	"context"
	"embed"
	"recap/internal/config"
	"recap/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var AppInstance App
var WailsContext *context.Context

type App struct {
	ctx context.Context
}

type AppMethods struct {
	FunctionsGiven               bool
	CSetScrTimer                 func(bool)
	CSetLLMTimer                 func(bool)
	CCheckTimers                 func() (bool, bool)
	CGetScreenshots              func(limit int) ([]db.CaptureScreenshotImage, error)
	CGetScreenshotById           func(id int) (*db.CaptureScreenshotImage, error)
	CGetScreenshotsNewerThan     func(timestamp int) ([]db.CaptureScreenshotImage, error)
	CGetScreenshotsOlderThan     func(timestamp int, limit int) ([]db.CaptureScreenshotImage, error)
	CDeleteScreenshotsById       func(ids []int) error
	CGenerateReportWithSelectScr func(ids []int) (*int64, error)
	CGetReports                  func(limit int) ([]db.Report, error)
	CGetReportById               func(id int) (*db.Report, error)
	CGetReportsNewerThan         func(id int) ([]db.Report, error)
	CGetReportsOlderThan         func(timestamp int, limit int) ([]db.Report, error)
	CDeleteReportsById           func(ids []int) error
	CGetConfig                   func() (*config.AppConfig, error)
	CGetDisplayValues            func() map[string]db.SettingDisplayProps
	CUpdateSettings              func(map[string]string) error
	CUpdateInfo                  func(map[string]string) error
	CWriteInfo                   func(key, value string) error
	CReadInfo                    func(key string) (*db.Info, error)
	CReadAllInfo                 func() (map[string]string, error)
}

func NewApp() *App {
	return &App{}
}

func NewAppMethods() *AppMethods {
	return &AppMethods{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) onDomReady(ctx context.Context) {
	WailsContext = &ctx
}

func LaunchAppInstance(assets embed.FS, methods *AppMethods, icon *[]byte) {
	AppInstance := NewApp()

	err := wails.Run(&options.App{
		Title:  "Recap",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		MinWidth:          600,
		MinHeight:         600,
		HideWindowOnClose: true,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
			Theme:                windows.SystemDefault,
		},
		Linux: &linux.Options{
			Icon:                *icon,
			WindowIsTranslucent: true,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "Recap",
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Recap",
				Message: "efeenesc",
				Icon:    *icon,
			},
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        AppInstance.startup,
		OnDomReady:       AppInstance.onDomReady,
		Bind: []interface{}{
			AppInstance,
			methods,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
