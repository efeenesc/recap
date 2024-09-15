package main

import (
	"fmt"
	"rcallport/internal/config"
	"rcallport/internal/db"
	"rcallport/internal/env"
	"rcallport/internal/llm"
	"rcallport/internal/schedule"
	"rcallport/internal/tray"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		config.Initialize() // Get config. It can be accessed via config.Config...
	}()
	go func() {
		defer wg.Done()
		env.Initialize() // Read environment variables into process
	}()

	wg.Wait()
	wg.Add(2)
	go func() {
		defer wg.Done()
		llm.Initialize() // Setup LLM (text, vision) connectors from config
	}()
	go func() {
		defer wg.Done()
		db.Initialize(false) // Initialize database, throw away the client afterwards
	}()

	wg.Wait()

	go func() {
		schedule.Initialize() // Start the schedule in which timers are configured for automatic screenshots and vision processing
	}()

	fmt.Println("Program started")
	tray.Initialize() // Blocking call
}
