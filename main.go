package main

import (
	_ "embed"
	"log"
	"recap/internal/db"
	"recap/internal/llm"
	_ "recap/internal/models/all" // Loads all models just to run their init()
	"recap/internal/schedule"
	"recap/internal/tray"
)

//go:embed assets/icon.ico
var iconBytes []byte

func addDbInitializers() *db.InitializerCallbacks {
	initializers := db.NewInitializers()
	initializers.InitSchedule = schedule.Initialize
	initializers.InitLLM = llm.Initialize
	initializers.FunctionsGiven = true
	return initializers
}

func main() {
	db.Initializers = *addDbInitializers()
	_, err := db.Initialize(false) // Initialize database, throw away the client afterwards
	if err != nil {
		log.Fatalf("Could not initialize database: %v\n", err.Error())
	}
	llm.Initialize()      // Setup LLM (text, vision) connectors from config
	schedule.Initialize() // Start the schedule in which timers are configured for automatic screenshots and vision processing
	tray.Initialize(&iconBytes)
	createApp()
}
