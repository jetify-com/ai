package client

// Response represents a response from the OpenRouter chat API
type Response struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

// Choice represents a single completion choice in the response
type Choice struct {
	Index        int              `json:"index"`
	Message      AssistantMessage `json:"message"`
	LogProbs     LogProbs         `json:"logprobs,omitempty"`
	FinishReason string           `json:"finish_reason"`
}

// Usage represents token usage information in a response
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}
