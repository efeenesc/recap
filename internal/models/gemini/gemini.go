package gemini

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"recap/internal/config"
	"strings"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIModel struct {
	client       *genai.Client
	clientTicker *time.Ticker
	model        string
	mu           sync.Mutex // Added mutex for thread-safety
}

// Start the ticker with 5-minute inactivity timeout for client closure
func (a *AIModel) startClientDeadline() {
	// If the ticker is already running, just reset the ticker to restart the inactivity period
	if a.clientTicker != nil {
		a.clientTicker.Reset(5 * time.Minute)
		fmt.Println("Resetting deadline timer")
		return
	}

	// Create the ticker if it's not already running
	a.clientTicker = time.NewTicker(5 * time.Minute)
	fmt.Println("Created deadline timer")

	go func() {
		for { // nolint: all
			select {
			case <-a.clientTicker.C:
				// Close the client after 5 minutes of inactivity
				if a.client != nil {
					fmt.Println("Closing client due to inactivity.")
					a.mu.Lock()
					a.client.Close()
					a.client = nil
					a.mu.Unlock()
					a.clientTicker.Stop()
				}
				return
			}
		}
	}()
}

// Joins all Gemini text response content (in chunks) into a single string
func joinContentToString(resp *genai.GenerateContentResponse) string {
	var sb strings.Builder

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				sb.WriteString(fmt.Sprint(part))
			}
		}
	}

	return sb.String()
}

// Sends a text generation request to the model with the given prompt.
func (a *AIModel) GenerateText(prompt string) (string, error) {
	client, ctx := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	model := client.GenerativeModel(a.model)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		return "", err
	}

	return joinContentToString(resp), nil
}

// Sends a single file for analysis
func (a *AIModel) DescribeScreenshot(fileName string, prompt string) (string, error) {
	client, ctx := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	return sendFileToGemini(client, ctx, a.model, fileName, prompt)
}

// Sends multiple files for analysis
func (a *AIModel) DescribeBulkScreenshots(fileNames []string, prompt string) (string, error) {
	client, ctx := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	descriptions := make([]string, 0, len(fileNames)) // Proper initialization

	for _, fn := range fileNames {
		res, err := sendFileToGemini(client, ctx, a.model, fn, prompt)
		if err != nil {
			fmt.Printf("An error occurred sending file to Gemini: %v\n", err.Error())
			return "", err
		}
		descriptions = append(descriptions, res) // Use append properly
	}

	return strings.Join(descriptions, "\n"), nil
}

// Sends a file for analysis to the Gemini model
func sendFileToGemini(client *genai.Client, ctx context.Context, modelName string, fileName string, prompt string) (string, error) {
	file, err := client.UploadFileFromPath(ctx, filepath.Join(config.Config.ScrPath, fileName), nil)
	if err != nil {
		return "", err
	}

	defer func() {
		err = client.DeleteFile(ctx, file.Name)
		if err != nil {
			fmt.Printf("Failed to delete file from Gemini (was it already deleted?): %v\n", err.Error())
		}
	}() // Defer file deletion

	model := client.GenerativeModel(modelName)
	resp, err := model.GenerateContent(ctx, genai.FileData{URI: file.URI}, genai.Text(prompt))

	if err != nil {
		return "", err
	}

	return joinContentToString(resp), nil
}

// Initializes the genai client if it currently doesn't exist, or returns the existing client.
// Remember to call a.startClientDeadline(). This sets a 5-minute timer before the client is stopped
func (a *AIModel) generateClient() (*genai.Client, context.Context) {
	ctx := context.Background()

	a.mu.Lock() // Ensure that client creation is thread-safe
	defer a.mu.Unlock()

	if a.client == nil {
		client, err := genai.NewClient(ctx, option.WithAPIKey(config.Config.GeminiAPIKey))
		if err != nil {
			log.Fatal(err)
			return nil, ctx
		}
		a.client = client
	}

	return a.client, ctx
}

// CreateAPIClient is a factory method for creating the AI client
func CreateAPIClient(model string) *AIModel {
	return &AIModel{model: model}
}
