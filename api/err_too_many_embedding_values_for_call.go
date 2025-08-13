package api

import "fmt"

// TooManyEmbeddingValuesForCallError indicates that too many values were provided for a single embedding call
type TooManyEmbeddingValuesForCallError struct {
	*AISDKError

	// Provider is the name of the AI provider
	Provider string

	// ModelID is the identifier of the model
	ModelID string

	// MaxEmbeddingsPerCall is the maximum number of embeddings allowed per call
	MaxEmbeddingsPerCall int

	// Values are the embedding values that were provided
	Values []any
}

// NewTooManyEmbeddingValuesForCallError creates a new TooManyEmbeddingValuesForCallError instance
// Parameters:
//   - provider: The name of the AI provider
//   - modelID: The identifier of the model
//   - maxEmbeddingsPerCall: The maximum number of embeddings allowed per call
//   - values: The embedding values that were provided
func NewTooManyEmbeddingValuesForCallError(provider, modelID string, maxEmbeddingsPerCall int, values []any) *TooManyEmbeddingValuesForCallError {
	message := fmt.Sprintf(
		"Too many values for a single embedding call. The %s model \"%s\" can only embed up to %d values per call, but %d values were provided.",
		provider,
		modelID,
		maxEmbeddingsPerCall,
		len(values),
	)

	return &TooManyEmbeddingValuesForCallError{
		AISDKError:           NewAISDKError("AI_TooManyEmbeddingValuesForCallError", message, nil),
		Provider:             provider,
		ModelID:              modelID,
		MaxEmbeddingsPerCall: maxEmbeddingsPerCall,
		Values:               values,
	}
}
