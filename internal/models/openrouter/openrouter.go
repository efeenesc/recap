package openrouter

import (
	"recap/internal/config"
	"recap/internal/models"
	"recap/internal/models/openai"
)

/**
OpenRouter API is essentially the same as OpenAI's API.
This creates an OpenAI API connector instance with a different ApiName and Endpoint parameter,
pointing the connector to OpenRouter's server
*/

// Initializes a new AIModel instance with the specified model name.
func CreateAPIClient(model string) models.TextVisionAPI {
	return &openai.AIModel{ApiName: "OpenRouter", Endpoint: "https://openrouter.ai/api/v1/chat/completions", Model: model, ApiKeyPtr: &config.Config.OpenRouterAPIKey}
}

func init() {
	models.RegisterAPI("OpenRouter", CreateAPIClient)
}
