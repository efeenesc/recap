package main

import (
	"fmt"
	"rcallport/internal/config"
	"rcallport/internal/db"
	"rcallport/internal/env"
	"rcallport/internal/llm"
	"rcallport/internal/schedule"
)

func main() {
	config.Initialize() // Get config. It can be accessed via config.Config...
	env.Initialize()    // Read environment variables into process

	fmt.Println("FUUUUUUUUUUUUUCK")
	llm.Initialize()     // Setup LLM (text, vision) connectors from config
	db.Initialize(false) // Initialize database, throw away the client afterwards

	schedule.Initialize() // Start the schedule in which timers are configured for automatic screenshots and vision processing

	createApp()
	fmt.Println("Program started")
}
