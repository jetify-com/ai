package client

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ChatSettings holds configuration for OpenRouter chat completions.
type ChatSettings struct {
	// LogitBias modifies the likelihood of specified tokens appearing in the completion.
	// Maps token IDs to bias values from -100 to 100. Values between -1 and 1 decrease
	// or increase likelihood of selection; values like -100 or 100 result in a ban or
	// exclusive selection of the token.
	// Example: {50256: -100} prevents the <|endoftext|> token from being generated.
	LogitBias map[int]float64 `json:"logit_bias,omitempty"`

	// Logprobs controls returning log probabilities of the tokens.
	// When Enabled is true, returns all logprobs if TopK is 0,
	// or returns top K logprobs if TopK > 0.
	// Note: Including logprobs increases response size and can slow down response times.
	Logprobs *LogprobSettings `json:"logprobs,omitempty"`

	// ParallelToolCalls enables parallel function calling during tool use.
	// Defaults to true if not set.
	ParallelToolCalls *bool `json:"parallel_tool_calls,omitempty"`

	// User is a unique identifier representing the end-user, which helps OpenRouter
	// monitor and detect abuse.
	User *string `json:"user,omitempty"`

	// Models is a list of model IDs to try in order if the primary model fails.
	// Example: ["anthropic/claude-2", "gryphe/mythomax-l2-13b"]
	Models []string `json:"models,omitempty"`

	// IncludeReasoning requests the model to return extra reasoning text in the response,
	// if the model supports it.
	IncludeReasoning *bool `json:"include_reasoning,omitempty"`
}

// CompletionSettings holds configuration for OpenRouter completions.
type CompletionSettings struct {
	// Echo returns the prompt in addition to the completion
	Echo *bool `json:"echo,omitempty"`

	// LogitBias modifies the likelihood of specified tokens appearing in the completion.
	// Maps token IDs to bias values from -100 to 100. Values between -1 and 1 decrease
	// or increase likelihood of selection; values like -100 or 100 result in a ban or
	// exclusive selection of the token.
	// Example: {50256: -100} prevents the <|endoftext|> token from being generated.
	LogitBias map[int]float64 `json:"logit_bias,omitempty"`

	// Logprobs controls returning log probabilities of the tokens.
	// When Enabled is true, returns all logprobs if TopK is 0,
	// or returns top K logprobs if TopK > 0.
	// Note: Including logprobs increases response size and can slow down response times.
	Logprobs *LogprobSettings `json:"logprobs,omitempty"`

	// Suffix is appended after a completion of inserted text
	Suffix *string `json:"suffix,omitempty"`

	// User is a unique identifier representing the end-user, which helps OpenRouter
	// monitor and detect abuse.
	User *string `json:"user,omitempty"`

	// Models is a list of model IDs to try in order if the primary model fails.
	// Example: ["openai/gpt-4", "anthropic/claude-2"]
	Models []string `json:"models,omitempty"`

	// IncludeReasoning requests the model to return extra reasoning text in the response,
	// if the model supports it.
	IncludeReasoning *bool `json:"include_reasoning,omitempty"`
}

// LogprobSettings represents the configuration for token log probabilities.
// It can be configured either as a boolean flag or with a number for top-N logprobs.
type LogprobSettings struct {
	// Enabled indicates if logprobs should be returned
	Enabled bool
	// TopK specifies how many top logprobs to return (if > 0)
	TopK int
}

// MarshalJSON implements custom JSON marshaling for LogprobSettings
func (l *LogprobSettings) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("null"), nil
	}
	if l.TopK > 0 {
		return []byte(strconv.Itoa(l.TopK)), nil
	}
	return []byte(strconv.FormatBool(l.Enabled)), nil
}

// UnmarshalJSON implements custom JSON unmarshaling for LogprobSettings
func (l *LogprobSettings) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		l.Enabled = false
		l.TopK = 0
		return nil
	}

	// Try as boolean first
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		l.Enabled = b
		l.TopK = 0
		return nil
	}

	// Try as number
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		l.Enabled = true
		l.TopK = n
		return nil
	}

	return fmt.Errorf("logprobs must be boolean or number, got %s", string(data))
}
