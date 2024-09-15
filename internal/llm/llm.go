package llm

import (
	"fmt"
	"log"
	"rcallport/internal/config"
	"rcallport/internal/db"
	"rcallport/internal/models"
	"rcallport/internal/models/gemini"
	"rcallport/internal/models/ollama"
	"time"
)

var visionModel models.ITextVisionModel
var textModel models.ITextVisionModel

func preprocessContext(caps []db.CaptureDescription) string {
	prompt := config.Config.LLM.Report.Prompt

	for _, cap := range caps {
		prompt += "BEGIN DESCRIPTION\n"
		prompt += cap.Description
		prompt += "END DESCRIPTION\n"
	}

	return prompt
}

func GenerateDailyReport() {
	dbCl, err := db.CreateConnection()
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}
	defer dbCl.Close()

	SendQueue() // Make sure all screenshots are described by AI

	todayCaps, err := db.GetCapturesToday(dbCl)
	if err != nil {
		fmt.Println("Error getting unprocessed captures:", err)
		return
	}

	if len(todayCaps) == 0 {
		fmt.Println("Could not find any captures from today")
		return
	}

	finalPrompt := preprocessContext(todayCaps)

	res, err := textModel.GenerateText(finalPrompt)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db.LogDailyReport(dbCl, res, todayCaps)
}

func SendQueue() {
	dbCl, err := db.CreateConnection()
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}
	defer dbCl.Close()

	fullQueue, err := db.GetUnprocessedCaptures(dbCl)
	if err != nil {
		fmt.Println("Error getting unprocessed captures:", err)
		return
	}

	fmt.Println(len(fullQueue))

	for i := 0; i < len(fullQueue); i += 15 {
		end := i + 15
		if end > len(fullQueue) {
			end = len(fullQueue)
		}
		processingQueue := fullQueue[i:end]

		for _, cap := range processingQueue {
			res, err := visionModel.DescribeScreenshot(cap.Filename, config.Config.LLM.Screenshot.Prompt)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", cap.Filename, err)
				continue
			}

			// Update the database with the Gemini response
			_, err = db.UpdateScreenshotDescription(dbCl, cap.CaptureID, res)
			if err != nil {
				fmt.Printf("Error updating description for capture %d: %v\n", cap.CaptureID, err)
			}
		}

		// If there are more items to process, wait for one minute
		if end < len(fullQueue) {
			fmt.Println("Waiting for 1 minute before processing next batch...")
			time.Sleep(time.Minute)
		}
	}

	fmt.Println("Queue processing completed")
}

func Initialize() {
	visionAPI := config.Config.LLM.Screenshot.API
	visionModelName := visionAPI.Connector

	textAPI := config.Config.LLM.Report.API
	textModelName := textAPI.Connector

	switch visionModelName {
	case "gemini":
		visionModel = gemini.CreateAPIClient(visionAPI.Model)

	case "ollama":
		visionModel = ollama.CreateAPIClient(visionAPI.Model)

	default:
		log.Fatalf("Unsupported screenshot/vision API: %s\n", visionAPI.Connector)
	}

	if visionModelName == textModelName {
		textModel = visionModel
	} else {
		switch textModelName {
		case "gemini":
			textModel = gemini.CreateAPIClient(textAPI.Model)

		case "ollama":
			textModel = ollama.CreateAPIClient(textAPI.Model)

		default:
			log.Fatalf("Unsupported report/text API: %s\n", visionAPI.Connector)
		}
	}

	fmt.Printf("Using %s %s for vision and %s %s for text", visionModelName, visionAPI.Model, textModelName, textAPI.Model)
}
