package api

import "net/http"

// EmbeddingOption represent the options for generating embeddings.
type EmbeddingOption func(*EmbeddingOptions)

// WithEmbeddingHeaders sets HTTP headers to be sent with the request.
// Only applicable for HTTP-based providers.
func WithEmbeddingHeaders(headers http.Header) EmbeddingOption {
	return func(o *EmbeddingOptions) {
		o.Headers = headers
	}
}

// EmbeddingOptions represents the options for generating embeddings.
type EmbeddingOptions struct {
	// Headers are additional HTTP headers to be sent with the request.
	// Only applicable for HTTP-based providers.
	Headers http.Header
}
