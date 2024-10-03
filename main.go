package main

import (
	"rcallport/internal/db"
	"rcallport/internal/llm"
	"rcallport/internal/schedule"
	"rcallport/internal/tray"
)

func addDbInitializers() *db.InitializerCallbacks {
	initializers := db.NewInitializers()
	initializers.InitSchedule = schedule.Initialize
	initializers.InitLLM = llm.Initialize
	initializers.FunctionsGiven = true
	return initializers
}

func main() {
	db.Initializers = *addDbInitializers()
	db.Initialize(false) // Initialize database, throw away the client afterwards

	llm.Initialize() // Setup LLM (text, vision) connectors from config

	schedule.Initialize() // Start the schedule in which timers are configured for automatic screenshots and vision processing

	go func() {
		tray.Initialize()
	}()

	createApp()
}
