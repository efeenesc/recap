package main

import (
	"embed"
	"rcallport/internal/app"
	"rcallport/internal/db"
	"rcallport/internal/llm"
	"rcallport/internal/schedule"
)

//go:embed all:frontend/build
var assets embed.FS

func addBindings() *app.AppMethods {
	methods := app.NewAppMethods()
	methods.CCheckTimers = schedule.AreTimersRunning
	methods.CSetLLMTimer = schedule.SetLLMScheduleState
	methods.CSetScrTimer = schedule.SetScrScheduleState
	methods.CGetScreenshots = db.GetScreenshots
	methods.CGetReports = db.GetReports
	methods.CGetScreenshotById = db.GetScreenshotById
	methods.CGenerateReportWithSelectScr = llm.GenerateReportWithSelectScr
	methods.CGetReportById = db.GetReportById
	methods.CGetConfig = db.LoadConfig
	methods.CGetDisplayValues = db.GetDisplayValues
	methods.FunctionsGiven = true
	return methods
}

func createApp() {
	methods := addBindings()
	app.LaunchAppInstance(assets, methods)
}
