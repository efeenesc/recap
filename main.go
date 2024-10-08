package main

import (
	"recap/internal/db"
	"recap/internal/llm"
	"recap/internal/schedule"
	"recap/internal/tray"
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
	db.Initialize(false)  // Initialize database, throw away the client afterwards
	llm.Initialize()      // Setup LLM (text, vision) connectors from config
	schedule.Initialize() // Start the schedule in which timers are configured for automatic screenshots and vision processing
	tray.Initialize()
	createApp()
}
