package api

// TokenLogProb represents the log probability for a single token.
type TokenLogProb struct {
	// Token is the text of the token
	Token string `json:"token"`

	// LogProb is the log probability of the token
	LogProb float64 `json:"logprob"`
}

// LogProb represents the log probability information for a token,
// including its top alternative tokens.
type LogProb struct {
	// Token is the text of the token
	Token string `json:"token"`

	// LogProb is the log probability of the token
	LogProb float64 `json:"logprob"`

	// TopLogProbs contains the log probabilities of alternative tokens
	TopLogProbs []TokenLogProb `json:"top_logprobs"`
}

// LogProbs represents a sequence of token log probabilities
type LogProbs []LogProb
