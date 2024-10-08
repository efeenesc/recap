package main

import (
	"embed"
	"recap/internal/app"
	"recap/internal/db"
	"recap/internal/llm"
	"recap/internal/schedule"
)

//go:embed all:frontend/build
var assets embed.FS

func addBindings() *app.AppMethods {
	methods := app.NewAppMethods()
	methods.CCheckTimers = schedule.AreTimersRunning
	methods.CSetLLMTimer = schedule.SetLLMScheduleState
	methods.CSetScrTimer = schedule.SetScrScheduleState
	methods.CGetScreenshots = db.GetScreenshots
	methods.CGetScreenshotById = db.GetScreenshotById
	methods.CGetScreenshotsNewerThan = db.GetScreenshotsNewerThan
	methods.CGetScreenshotsOlderThan = db.GetScreenshotsOlderThan
	methods.CDeleteScreenshotsById = db.DeleteScreenshotsById
	methods.CGenerateReportWithSelectScr = llm.GenerateReportWithSelectScr
	methods.CGetReports = db.GetReports
	methods.CGetReportById = db.GetReportById
	methods.CGetReportsNewerThan = db.GetReportsNewerThan
	methods.CGetReportsOlderThan = db.GetReportsOlderThan
	methods.CDeleteReportsById = db.DeleteReportsById
	methods.CGetConfig = db.LoadConfig
	methods.CGetDisplayValues = db.GetDisplayValues
	methods.CUpdateSettings = db.UpdateSettings
	methods.FunctionsGiven = true
	return methods
}

func createApp() {
	methods := addBindings()
	app.LaunchAppInstance(assets, methods)
}
