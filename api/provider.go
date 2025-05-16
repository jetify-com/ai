package api

// Provider is a provider for language and text embedding models.
type Provider interface {
	// LanguageModel returns the language model with the given id.
	// The model id is then passed to the provider function to get the model.
	//
	// Parameters:
	//   modelID: The id of the model to return.
	//
	// Returns:
	//   The language model associated with the id
	//   error of type NoSuchModelError if no such model exists
	LanguageModel(modelID string) (LanguageModel, error)

	// TextEmbeddingModel returns the text embedding model with the given id.
	// The model id is then passed to the provider function to get the model.
	//
	// Parameters:
	//   modelID: The id of the model to return.
	//
	// Returns:
	//   The text embedding model associated with the id
	//   error of type NoSuchModelError if no such model exists
	TextEmbeddingModel(modelID string) (EmbeddingModel[string], error)

	// ImageModel returns the image model with the given id.
	// The model id is then passed to the provider function to get the model.
	// This method is optional and may return nil if image models are not supported.
	//
	// Parameters:
	//   modelID: The id of the model to return.
	//
	// Returns:
	//   The image model associated with the id, or nil if image models are not supported
	//   error of type NoSuchModelError if no such model exists and image models are supported
	ImageModel(modelID string) (*ImageModel, error)
}
