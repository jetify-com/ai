package api

import (
	"context"
	"net/http"
	"time"
)

// ImageModel is a specification for an image generation model that implements
// the image model interface version 1.
type ImageModel interface {
	// SpecificationVersion returns which image model interface version is implemented.
	// This will allow us to evolve the image model interface and retain backwards
	// compatibility. The different implementation versions can be handled as a
	// discriminated union on our side.
	SpecificationVersion() string

	// ProviderName returns the name of the provider for logging purposes.
	ProviderName() string

	// ModelID returns the provider-specific model ID for logging purposes.
	ModelID() string

	// MaxImagesPerCall returns the limit of how many images can be generated in a single API call.
	// If undefined, we will max generate one image per call.
	MaxImagesPerCall() *int

	// DoGenerate generates an array of images based on the given prompt.
	DoGenerate(ctx context.Context, prompt string, opts ...ImageCallOption) ImageResponse
}

// ImageResponse represents the response from generating images.
type ImageResponse struct {
	// Images are the generated images as base64 encoded strings or binary data.
	// The images should be returned without any unnecessary conversion.
	// If the API returns base64 encoded strings, the images should be returned
	// as base64 encoded strings. If the API returns binary data, the images should
	// be returned as binary data.
	Images []ImageData

	// Warnings for the call, e.g. unsupported settings.
	Warnings []ImageCallWarning

	// Response contains information for telemetry and debugging purposes.
	Response ImageResponseMetadata
}

// ImageData represents either a base64 encoded string or binary data for an image
type ImageData interface {
	// IsImageData is a marker method to ensure type safety
	IsImageData()
}

// Base64Image represents an image as a base64 encoded string
type Base64Image string

func (Base64Image) IsImageData() {}

// BinaryImage represents an image as binary data
type BinaryImage []byte

func (BinaryImage) IsImageData() {}

// ImageResponseMetadata contains response information for telemetry and debugging purposes.
type ImageResponseMetadata struct {
	// Timestamp is the timestamp for the start of the generated response.
	Timestamp time.Time

	// ModelID is the ID of the response model that was used to generate the response.
	ModelID string

	// Headers are the response headers.
	Headers http.Header
}
