package api

import (
	"context"
)

// Embedding is a vector, i.e. an array of numbers.
// It is e.g. used to represent a text as a vector of word embeddings.
type Embedding []float64

// EmbeddingModel is a specification for an embedding model that implements the embedding model
// interface version 1.
//
// T is the type of the values that the model can embed.
// This will allow us to go beyond text embeddings in the future,
// e.g. to support image embeddings
type EmbeddingModel[T any] interface {
	// SpecificationVersion returns which embedding model interface version is implemented.
	// This will allow us to evolve the embedding model interface and retain backwards
	// compatibility. The different implementation versions can be handled as a discriminated
	// union on our side.
	SpecificationVersion() string

	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model ID for logging purposes.
	ModelID() string

	// MaxEmbeddingsPerCall returns the limit of how many embeddings can be generated in a single API call.
	MaxEmbeddingsPerCall() *int

	// SupportsParallelCalls returns if the model can handle multiple embedding calls in parallel.
	SupportsParallelCalls() bool

	// DoEmbed generates a list of embeddings for the given input values.
	//
	// Naming: "do" prefix to prevent accidental direct usage of the method
	// by the user.
	DoEmbed(ctx context.Context, values []T, opts ...EmbeddingOption) EmbeddingResponse
}

// EmbeddingResponse represents the response from generating embeddings.
type EmbeddingResponse struct {
	// Embeddings are the generated embeddings. They are in the same order as the input values.
	Embeddings []Embedding

	// Usage contains token usage information. We only have input tokens for embeddings.
	Usage *EmbeddingUsage

	// RawResponse contains optional raw response information for debugging purposes.
	RawResponse *EmbeddingRawResponse
}

// EmbeddingUsage represents token usage information.
type EmbeddingUsage struct {
	Tokens int
}

// EmbeddingRawResponse contains raw response information for debugging.
type EmbeddingRawResponse struct {
	// Headers are the response headers.
	Headers map[string]string
}
