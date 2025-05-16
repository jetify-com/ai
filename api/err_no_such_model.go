package api

import (
	"fmt"
)

// ModelType represents the type of AI model
type ModelType string

const (
	// LanguageModelType represents a language model type
	LanguageModelType ModelType = "languageModel"
	// TextEmbeddingModelType represents a text embedding model type
	TextEmbeddingModelType ModelType = "textEmbeddingModel"
	// ImageModelType represents an image model type
	ImageModelType ModelType = "imageModel"
)

// NoSuchModelError indicates that the requested model does not exist
type NoSuchModelError struct {
	*AISDKError

	// ModelID is the identifier of the model that was not found
	ModelID string

	// ModelType is the type of model that was requested
	ModelType ModelType
}

// NewNoSuchModelError creates a new NoSuchModelError instance
// Parameters:
//   - modelID: The identifier of the model that was not found
//   - modelType: The type of model that was requested
func NewNoSuchModelError(modelID string, modelType ModelType) *NoSuchModelError {
	message := fmt.Sprintf("No such %s: %s", modelType, modelID)
	return &NoSuchModelError{
		AISDKError: NewAISDKError("AI_NoSuchModelError", message, nil),
		ModelID:    modelID,
		ModelType:  modelType,
	}
}
