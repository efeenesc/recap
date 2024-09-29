package llm

import (
	"database/sql"
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

func GenerateDailyReport() (*int64, error) {
	dbCl, err := db.CreateConnection()
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return nil, err
	}
	defer dbCl.Close()

	SendQueue() // Make sure all screenshots are described by AI

	todayCaps, err := db.GetCapturesToday(dbCl)
	if err != nil {
		fmt.Println("Error getting unprocessed captures:", err)
		return nil, err
	}

	if len(todayCaps) == 0 {
		fmt.Println("Could not find any captures from today")
		return nil, nil
	}

	finalPrompt := preprocessContext(todayCaps)

	res, err := textModel.GenerateText(finalPrompt)
	if err != nil {
		log.Fatalln(err.Error())
		// return nil, err
	}

	return db.LogDailyReport(dbCl, res, todayCaps)
}

func GenerateReportWithSelectScr(ids []int) (*int64, error) {
	log.Println("Starting report generation")
	dbCl, err := db.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("error creating database connection: %w", err)
	}
	defer dbCl.Close()

	caps, err := db.GetScreenshotByIds(dbCl, ids)
	if err != nil {
		return nil, fmt.Errorf("error getting unprocessed captures: %w", err)
	}

	var toProcess []db.CaptureScreenshot
	var ready []db.CaptureDescription
	for _, cap := range caps {
		if cap.Description != nil {
			ready = append(ready, db.CaptureDescription{
				CaptureID:   cap.CaptureID,
				Timestamp:   cap.Timestamp,
				Description: *cap.Description,
			})
		} else {
			toProcess = append(toProcess, cap)
		}
	}

	newDescs, err := SendQueueFromObject(dbCl, toProcess)
	if err != nil {
		return nil, fmt.Errorf("error sending queue: %w", err)
	}

	descs := append(ready, newDescs...)

	if len(descs) == 0 {
		return nil, fmt.Errorf("no screenshot descriptions found")
	}

	finalPrompt := preprocessContext(descs)

	res, err := textModel.GenerateText(finalPrompt)
	if err != nil {
		return nil, fmt.Errorf("error generating text: %w", err)
	}

	log.Println("Logging report")
	return db.LogDailyReport(dbCl, res, descs)
}

func SendQueueFromObject(dbCl *sql.DB, scrs []db.CaptureScreenshot) ([]db.CaptureDescription, error) {
	if dbCl == nil {
		newCl, err := db.CreateConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to create database connection: %w", err)
		}
		defer newCl.Close()
		dbCl = newCl
	}

	if scrs == nil {
		return nil, fmt.Errorf("input slice of CaptureScreenshot is nil")
	}

	var returnQ []db.CaptureDescription
	batchSize := 15

	for i := 0; i < len(scrs); i += batchSize {
		end := min(i+batchSize, len(scrs))
		processingQueue := scrs[i:end]

		fmt.Println("Starting batch processing")

		for _, cap := range processingQueue {
			if cap.Filename == "" {
				log.Printf("Warning: Empty filename for capture ID %d, skipping", cap.CaptureID)
				continue
			}

			res, err := visionModel.DescribeScreenshot(cap.Filename, config.Config.LLM.Screenshot.Prompt)
			if err != nil {
				log.Printf("Error processing file %s: %v", cap.Filename, err)
				continue
			}

			newDescObj := db.CaptureDescription{
				CaptureID:   cap.CaptureID,
				Timestamp:   cap.Timestamp,
				Description: res,
			}
			returnQ = append(returnQ, newDescObj)

			fmt.Printf("Processed capture ID %d: %s\n", cap.CaptureID, truncateString(newDescObj.Description, 50))

			if _, err := db.UpdateScreenshotDescription(dbCl, cap.CaptureID, res); err != nil {
				log.Printf("Error updating description for capture %d: %v", cap.CaptureID, err)
			}
		}

		if end < len(scrs) {
			log.Println("Waiting for 1 minute before processing next batch...")
			time.Sleep(time.Minute)
		}
	}

	log.Printf("Queue processing completed. Processed %d items.", len(returnQ))
	return returnQ, nil
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SendQueue() {
	dbCl, err := db.CreateConnection()
	if err != nil {
		fmt.Println("Error creating database connection:", err)
		return
	}
	defer dbCl.Close()
	// defer instance.AppInstance.SendLLMRanMessage()

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

	fmt.Println("INITIALIZING THIS SHIT")

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

	fmt.Printf("Using %s %s for vision and %s %s for text\n", visionModelName, visionAPI.Model, textModelName, textAPI.Model)
}
