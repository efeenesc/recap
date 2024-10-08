package llm

import (
	"database/sql"
	"fmt"
	"log"
	"recap/internal/config"
	"recap/internal/db"
	"recap/internal/models"
	"recap/internal/models/gemini"
	"recap/internal/models/ollama"
	"time"
)

var visionModel models.ITextVisionModel
var textModel models.ITextVisionModel

// Processes descriptions from screenshots into formatted context to be passed to the LLM
func preprocessContext(caps []db.CaptureDescription) string {
	prompt := config.Config.ReportPrompt

	for _, cap := range caps {
		prompt += "BEGIN DESCRIPTION\n"
		prompt += cap.Description
		prompt += "END DESCRIPTION\n"
	}

	return prompt
}

// Generates a daily report based on captures from today.
// It retrieves today's captures, processes them through AI for descriptions,
// and logs the resulting report. Returns the ID of the logged report or an error.
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

	return db.LogDailyReport(dbCl, res, todayCaps, config.Config.ReportAPI, config.Config.ReportModel)
}

// Generates a report using a selected list of screenshot IDs.
// It retrieves the specified captures, processes those that lack descriptions,
// and generates a report based on the combined descriptions. Returns the ID of the logged report or an error.
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
	return db.LogDailyReport(dbCl, res, descs, config.Config.ReportAPI, config.Config.ReportModel)
}

// Processes a batch of screenshot captures to generate descriptions.
// It describes each screenshot using the vision model and updates the descriptions in the database.
// Returns a slice of CaptureDescription with the generated descriptions or an error.
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
		fmt.Printf("input slice of CaptureScreenshot is nil, returning nil. Everything good")
		return nil, nil
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

			res, err := visionModel.DescribeScreenshot(cap.Filename, config.Config.DescGenPrompt)
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

			if _, err := db.UpdateScreenshotDescription(dbCl, cap.CaptureID, res, config.Config.DescGenAPI, config.Config.DescGenModel); err != nil {
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

// Truncates a string to a specified maximum length and appends ellipsis if truncated.
// Returns the truncated string.
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// Returns the smaller of two integers a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Processes unprocessed captures from the database in batches,
// generates descriptions using the vision model, and updates the database.
// Waits for one minute between batches if necessary.
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

	for i := 0; i < len(fullQueue); i += 15 {
		end := i + 15
		if end > len(fullQueue) {
			end = len(fullQueue)
		}
		processingQueue := fullQueue[i:end]

		for _, cap := range processingQueue {
			res, err := visionModel.DescribeScreenshot(cap.Filename, config.Config.DescGenPrompt)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", cap.Filename, err)
				continue
			}

			// Update the database with the Gemini response
			_, err = db.UpdateScreenshotDescription(dbCl, cap.CaptureID, res, config.Config.DescGenAPI, config.Config.DescGenModel)
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

// Sets up the vision and text models based on configuration settings.
// It selects the appropriate API clients for image description and report generation.
func Initialize() {
	visionAPI := config.Config.DescGenAPI
	visionModelName := config.Config.DescGenModel

	textAPI := config.Config.ReportAPI
	textModelName := config.Config.ReportModel

	switch visionAPI {
	case "gemini":
		visionModel = gemini.CreateAPIClient(visionModelName)

	case "ollama":
		visionModel = ollama.CreateAPIClient(visionModelName)

	default:
		log.Fatalf("Unsupported screenshot/vision API: %s\n", visionAPI)
	}

	if visionModelName == textModelName {
		textModel = visionModel
	} else {
		switch textAPI {
		case "gemini":
			textModel = gemini.CreateAPIClient(textModelName)

		case "ollama":
			textModel = ollama.CreateAPIClient(textModelName)

		default:
			log.Fatalf("Unsupported report/text API: %s\n", textAPI)
		}
	}

	fmt.Printf("Using %s %s for vision and %s %s for text\n", visionAPI, visionModelName, textAPI, textModelName)
}
