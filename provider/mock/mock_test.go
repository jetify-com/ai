package mock

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestGenerateModel(t *testing.T) {
	ctx := context.Background()
	messages := []api.Message{
		&api.UserMessage{Content: []api.ContentBlock{&api.TextBlock{Text: "test"}}},
	}
	opts := api.CallOptions{}

	tests := []struct {
		name    string
		results []MockResult
	}{
		{
			name: "single successful result",
			results: []MockResult{
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "Success"}},
				}},
			},
		},
		{
			name: "single error result",
			results: []MockResult{
				{Error: errors.New("rate limit exceeded")},
			},
		},
		{
			name: "multiple results in order",
			results: []MockResult{
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "First"}},
				}},
				{Error: errors.New("second fails")},
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "Third"}},
				}},
			},
		},
		{
			name: "multiple successful results",
			results: []MockResult{
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "First"}},
				}},
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "Second"}},
				}},
				{Response: &api.Response{
					Content: []api.ContentBlock{&api.TextBlock{Text: "Third"}},
				}},
			},
		},
		{
			name: "multiple error results",
			results: []MockResult{
				{Error: errors.New("first error")},
				{Error: errors.New("second error")},
				{Error: errors.New("third error")},
			},
		},
		{
			name:    "empty results slice",
			results: []MockResult{},
		},
		{
			name:    "nil results",
			results: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := NewGenerateModel(test.results)

			// Accumulate responses
			var accumulated []MockResult
			if test.results != nil {
				accumulated = make([]MockResult, 0, len(test.results))
			}

			// Check AssertCount before any calls
			mockT := &mockTestingT{}
			model.AssertCount(mockT)
			if len(test.results) == 0 {
				assert.False(t, mockT.failed, "AssertCount should have passed with 0 calls when expecting 0 results")
			} else {
				assert.True(t, mockT.failed, "AssertCount should have failed with 0 calls when expecting %d results", len(test.results))
			}

			// Call Generate() len(results) + 1 times
			totalCalls := len(test.results) + 1
			for i := 0; i < totalCalls; i++ {
				resp, err := model.Generate(ctx, messages, opts)

				// Accumulate results up to len(results)
				if i < len(test.results) {
					accumulated = append(accumulated, MockResult{
						Response: resp,
						Error:    err,
					})
				}

				// Check AssertCount after each call
				mockT = &mockTestingT{}
				model.AssertCount(mockT)

				// Should pass only when we've made exactly len(results) calls
				// i+1 is the number of calls made so far
				if i+1 == len(test.results) {
					assert.False(t, mockT.failed, "AssertCount should have passed at call %d (exactly %d results)", i+1, len(test.results))
				} else {
					assert.True(t, mockT.failed, "AssertCount should have failed at call %d (expected %d results)", i+1, len(test.results))
				}
			}

			// Compare accumulated results with original
			assert.Equal(t, test.results, accumulated)
		})
	}
}

func TestStreamCallTracking(t *testing.T) {
	ctx := context.Background()
	messages := []api.Message{
		&api.UserMessage{Content: []api.ContentBlock{&api.TextBlock{Text: "test"}}},
	}
	opts := api.CallOptions{}

	t.Run("AssertCount fails when Stream is called", func(t *testing.T) {
		model := NewGenerateModel([]MockResult{})

		// Call Stream method
		_, err := model.Stream(ctx, messages, opts)
		assert.Error(t, err, "Stream should return an error")
		assert.Contains(t, err.Error(), "not implemented", "Stream should return 'not implemented' error")

		// AssertCount should fail because Stream was called
		mockT := &mockTestingT{}
		model.AssertCount(mockT)
		assert.True(t, mockT.failed, "AssertCount should fail when Stream was called")
	})

	t.Run("AssertCount passes when Stream is never called", func(t *testing.T) {
		model := NewGenerateModel([]MockResult{})

		// Don't call Stream, only check AssertCount
		mockT := &mockTestingT{}
		model.AssertCount(mockT)
		assert.False(t, mockT.failed, "AssertCount should pass when Stream was never called")
	})

	t.Run("AssertCount fails when Stream is called even with correct Generate calls", func(t *testing.T) {
		model := NewGenerateModel([]MockResult{
			{Response: &api.Response{Content: []api.ContentBlock{&api.TextBlock{Text: "test"}}}},
		})

		// Make the expected Generate call
		_, err := model.Generate(ctx, messages, opts)
		assert.NoError(t, err)

		// Also call Stream
		_, err = model.Stream(ctx, messages, opts)
		assert.Error(t, err)

		// AssertCount should fail because Stream was called, even though Generate calls match
		mockT := &mockTestingT{}
		model.AssertCount(mockT)
		assert.True(t, mockT.failed, "AssertCount should fail when Stream was called, even with correct Generate calls")
	})
}

// mockTestingT implements the mock.T interface for testing
type mockTestingT struct {
	failed bool
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.failed = true
}

func (m *mockTestingT) FailNow() {
	m.failed = true
}

func (m *mockTestingT) Helper() {}

func TestGenerateModelOptions(t *testing.T) {
	tests := []struct {
		name         string
		opts         []GenerateModelOption
		wantProvider string
		wantModelID  string
	}{
		{
			name:         "default values",
			opts:         nil,
			wantProvider: "mock-provider",
			wantModelID:  "mock-model",
		},
		{
			name:         "custom provider only",
			opts:         []GenerateModelOption{WithProviderName("custom-provider")},
			wantProvider: "custom-provider",
			wantModelID:  "mock-model",
		},
		{
			name:         "custom model ID only",
			opts:         []GenerateModelOption{WithModelID("custom-model")},
			wantProvider: "mock-provider",
			wantModelID:  "custom-model",
		},
		{
			name: "both custom provider and model ID",
			opts: []GenerateModelOption{
				WithProviderName("openai"),
				WithModelID("gpt-4"),
			},
			wantProvider: "openai",
			wantModelID:  "gpt-4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			model := NewGenerateModel([]MockResult{}, test.opts...)

			assert.Equal(t, test.wantProvider, model.ProviderName())
			assert.Equal(t, test.wantModelID, model.ModelID())
		})
	}
}
