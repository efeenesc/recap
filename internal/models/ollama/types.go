package ollama

type OllamaRequest struct {
	Model   string    `json:"model"`
	Prompt  string    `json:"prompt"`
	Stream  bool      `json:"stream"`
	Images  *[]string `json:"images"`
	Options struct {
		Temperature float32 `json:""`
	}
}

type OllamaFullResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"` // Date field. Parse if needed
	Done               bool   `json:"done"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}
