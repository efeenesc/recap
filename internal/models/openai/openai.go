package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"recap/internal/config"
	"recap/internal/models"
	"recap/internal/utils"
	"strings"
	"sync"
	"time"
)

type AIModel struct {
	ApiName      string
	ApiKeyPtr    *string
	client       *http.Client
	clientTicker *time.Ticker
	Endpoint     string // Used to set endpoint to OpenAI's API, or OpenRouter's
	Model        string
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
		for { // nolint: all
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

	return sendToOpenAI(client, a.Model, nil, prompt, a.Endpoint, *a.ApiKeyPtr)
}

// Generates a description for a screenshot specified by its filename
// using the AI client. It returns the description or an error if client creation fails.
func (a *AIModel) DescribeScreenshot(fileName string, prompt string) (string, error) {
	client := a.generateClient()
	if client == nil {
		return "", fmt.Errorf("failed to create client")
	}
	defer a.startClientDeadline()

	return sendToOpenAI(client, a.Model, &fileName, prompt, a.Endpoint, *a.ApiKeyPtr)
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
		res, err := sendToOpenAI(client, a.Model, &fn, prompt, a.Endpoint, *a.ApiKeyPtr)
		if err != nil {
			return "", fmt.Errorf("error sending file to OpenAI: %w", err)
		}
		descriptions = append(descriptions, res)
	}

	return strings.Join(descriptions, "\n"), nil
}

// Sends a request to the OpenAI API with the specified client, model,
// and image data (if applicable). It returns the response from the API or an error.
func sendToOpenAI(client *http.Client, modelName string, fileName *string, prompt string, endpoint string, apiKey string) (string, error) {
	var images []string
	if fileName != nil {
		imageBase64 := utils.ReadImageToBase64(*fileName)
		if imageBase64 == "" {
			return "", fmt.Errorf("failed to read image file")
		}
		images = append(images, imageBase64)
	}

	apiBearerAuth := fmt.Sprintf("Bearer %s", apiKey)

	requestBody := OpenAIRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: false,
		Images: &images,
	}

	preparedBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON request body: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(preparedBody))
	if err != nil {
		return "", fmt.Errorf("error creating request to OpenAI: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiBearerAuth) // Replace with your actual API key

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to OpenAI: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		errorMessage, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("OpenAI API returned non-2xx status: %d - %s", res.StatusCode, string(errorMessage))
	}

	readRes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response from OpenAI: %w", err)
	}

	var openaiResponse OpenAIFullResponse
	err = json.Unmarshal(readRes, &openaiResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding OpenAI response: %v", err.Error())
	}

	return openaiResponse.Response, nil
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

func (a *AIModel) GetAPIName() string {
	return a.ApiName
}

func (a *AIModel) GetAPIModelName() string {
	return a.Model
}

// Initializes a new AIModel instance with the specified model name.
func CreateAPIClient(model string) models.TextVisionAPI {
	return &AIModel{ApiName: "OpenAI", Endpoint: "https://api.openai.com/v1/completions", Model: model, ApiKeyPtr: &config.Config.OpenAIAPIKey}
}

func init() {
	models.RegisterAPI("OpenAI", CreateAPIClient)
}
