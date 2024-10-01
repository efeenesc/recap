package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rcallport/internal/utils"
	"strings"
	"sync"
	"time"
)

type AIModel struct {
	client       *http.Client
	clientTicker *time.Ticker
	model        string
	mu           sync.Mutex // Added mutex for thread-safety
}

// Starts or resets a timer that closes the HTTP client
// after 5 minutes of inactivity. If the client is already active, it resets the timer.
func (a *AIModel) startClientDeadline() {
	if a.clientTicker != nil {
		a.clientTicker.Reset(5 * time.Minute)
		fmt.Println("Resetting deadline timer")
		return
	}

	a.clientTicker = time.NewTicker(5 * time.Minute)
	fmt.Println("Created deadline timer")

	go func() {
		for {
			select {
			case <-a.clientTicker.C:
				a.mu.Lock()
				if a.client != nil {
					fmt.Println("Closing client due to inactivity.")
					a.client.CloseIdleConnections()
					a.client = nil
					a.clientTicker.Stop()
					a.clientTicker = nil // Set to nil after stopping
				}
				a.mu.Unlock()
				return
			}
		}
	}()
}

// Generates text based on the provided prompt using the AI client.
// It returns the generated text or an error if client creation fails.
func (a *AIModel) GenerateText(prompt string) (string, error) {
	client := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	return sendToOllama(client, a.model, nil, prompt)
}

// Generates a description for a screenshot specified by its filename
// using the AI client. It returns the description or an error if client creation fails.
func (a *AIModel) DescribeScreenshot(fileName string, prompt string) (string, error) {
	client := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	return sendToOllama(client, a.model, &fileName, prompt)
}

// Generates descriptions for multiple screenshots
// provided in the fileNames slice. It returns concatenated descriptions or an error.
func (a *AIModel) DescribeBulkScreenshots(fileNames []string, prompt string) (string, error) {
	client := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	var descriptions []string

	for _, fn := range fileNames {
		res, err := sendToOllama(client, a.model, &fn, prompt)
		if err != nil {
			return "", fmt.Errorf("error sending file to Ollama: %w", err)
		}
		descriptions = append(descriptions, res)
	}

	return strings.Join(descriptions, "\n"), nil
}

// Sends a request to the Ollama API with the specified client, model,
// and image data (if applicable). It returns the response from the API or an error.
func sendToOllama(client *http.Client, modelName string, fileName *string, prompt string) (string, error) {
	var images []string
	if fileName != nil {
		imageBase64 := utils.ReadImageToBase64(*fileName)
		if imageBase64 == "" {
			return "", fmt.Errorf("failed to read image file")
		}
		images = []string{imageBase64}
	}

	requestBody := OllamaRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: false,
		Images: &images,
	}

	preparedBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON request body: %w", err)
	}

	res, err := client.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(preparedBody))
	readRes, _ := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Error from Ollama: %s\n", readRes)
		return "", fmt.Errorf("error sending request to Ollama: %v", err.Error())
	}
	defer res.Body.Close()

	var ollamaResponse OllamaFullResponse
	err = json.Unmarshal(readRes, &ollamaResponse)
	if err != nil {
		fmt.Printf("error decoding Ollama response: %v", err.Error())
		return "", fmt.Errorf("error decoding Ollama response: %v", err.Error())
	}

	return ollamaResponse.Response, nil
}

// Creates and returns a new HTTP client if one does not already exist.
// It ensures thread-safe access to the client instance.
func (a *AIModel) generateClient() *http.Client {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.client == nil {
		a.client = &http.Client{}
	}

	return a.client
}

// Initializes a new AIModel instance with the specified model name.
func CreateAPIClient(model string) *AIModel {
	return &AIModel{model: model}
}
