// Package mock provides mock implementations for testing.
//
// Example usage:
//
//	mock := mock.NewGenerateModel([]mock.MockResult{
//		{Response: &api.Response{Content: []api.ContentBlock{&api.TextBlock{Text: "Hello"}}}},
//		{Error: errors.New("rate limit exceeded")},
//		{Response: &api.Response{Content: []api.ContentBlock{&api.TextBlock{Text: "World"}}}},
//	})
//
//	// First call returns success
//	resp1, err1 := mock.Generate(ctx, messages, opts) // returns "Hello", nil
//
//	// Second call returns error
//	resp2, err2 := mock.Generate(ctx, messages, opts) // returns {}, error
//
//	// Third call returns success
//	resp3, err3 := mock.Generate(ctx, messages, opts) // returns "World", nil
//
//	// Verify all expectations were met
//	mock.AssertCount(t)
package mock

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

// MockResult represents either a successful response or an error
type MockResult struct {
	Response *api.Response
	Error    error
}

type GenerateModel struct {
	results         []MockResult
	callCount       atomic.Int32
	streamCallCount atomic.Int32
	providerName    string
	modelID         string
}

// T is an interface that captures the testing.T methods we need
type T interface {
	Errorf(format string, args ...any)
	FailNow()
	Helper()
}

// GenerateModelOption is a functional option for configuring a GenerateModel
type GenerateModelOption func(*GenerateModel)

// WithProviderName sets the provider name for the mock model
func WithProviderName(name string) GenerateModelOption {
	return func(m *GenerateModel) {
		m.providerName = name
	}
}

// WithModelID sets the model ID for the mock model
func WithModelID(id string) GenerateModelOption {
	return func(m *GenerateModel) {
		m.modelID = id
	}
}

// NewGenerateModel creates a new mock GenerateModel with the given results.
// The results will be returned in order as Generate is called.
// If results is nil, it will be treated as an empty slice.
// Optional functional options can be provided to customize provider name and model ID.
func NewGenerateModel(results []MockResult, opts ...GenerateModelOption) *GenerateModel {
	if results == nil {
		results = []MockResult{}
	}

	m := &GenerateModel{
		results:      results,
		providerName: "mock-provider", // Default
		modelID:      "mock-model",    // Default
	}

	// Apply options
	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *GenerateModel) ProviderName() string {
	return m.providerName
}

func (m *GenerateModel) ModelID() string {
	return m.modelID
}

func (m *GenerateModel) SupportedUrls() []api.SupportedURL {
	return nil
}

func (m *GenerateModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.Response, error) {
	// Atomically increment and get the new count
	newCount := m.callCount.Add(1)
	index := newCount - 1

	if int(index) < len(m.results) {
		result := m.results[index]
		return result.Response, result.Error
	}
	return &api.Response{}, nil
}

func (m *GenerateModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.StreamResponse, error) {
	m.streamCallCount.Add(1)
	return nil, errors.New("Stream: not implemented")
}

// AssertCount verifies that all expected results have been consumed
// and that Stream was never called.
// It fails the test if there are unused results or if Stream was called.
func (m *GenerateModel) AssertCount(t T) {
	t.Helper()
	callCount := int(m.callCount.Load())
	streamCallCount := int(m.streamCallCount.Load())
	assert.Equal(t, len(m.results), callCount,
		"GenerateModel: expected %d Generate calls, but got %d", len(m.results), callCount)
	assert.Equal(t, 0, streamCallCount,
		"GenerateModel: expected 0 Stream calls, but got %d", streamCallCount)
}
