package client

// LogProbs represents the logprobs structure returned by the OpenRouter
// chat API.
type LogProbs struct {
	Content []LogProb `json:"content,omitempty"`
}

// LogProb represents a single token's log probability information
type LogProb struct {
	Token       string       `json:"token"`
	LogProb     float64      `json:"logprob"`
	TopLogProbs []TopLogProb `json:"top_logprobs,omitempty"`
}

// TopLogProb represents a single top logprob entry
type TopLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
}

// CompletionLogProbs represents the logprobs structure returned by OpenRouter
// completions API.
type CompletionLogProbs struct {
	Tokens        []string             `json:"tokens"`
	TokenLogProbs []float64            `json:"token_logprobs"`
	TopLogProbs   []map[string]float64 `json:"top_logprobs,omitempty"`
}
