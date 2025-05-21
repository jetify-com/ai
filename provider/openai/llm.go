package openai

import (
	"context"
	"net/url"

	"github.com/openai/openai-go"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/openai/internal/codec"
)

// ModelOption is a function type that modifies a LanguageModel.
type ModelOption func(*LanguageModel)

// WithClient returns a ModelOption that sets the client.
func WithClient(client openai.Client) ModelOption {
	// TODO: Instead of only supporting a single client, we can "flatten"
	// the options supported by the OpenAI SDK.
	return func(m *LanguageModel) {
		m.client = client
	}
}

// LanguageModel represents an OpenAI language model.
type LanguageModel struct {
	modelID string
	client  openai.Client
}

var _ api.LanguageModel = &LanguageModel{}

// NewLanguageModel creates a new OpenAI language model.
func NewLanguageModel(modelID string, opts ...ModelOption) *LanguageModel {
	// Create model with default settings
	model := &LanguageModel{
		modelID: modelID,
		client:  openai.NewClient(), // Default client
	}

	// Apply options
	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (m *LanguageModel) SpecificationVersion() string {
	return "v1"
}

func (m *LanguageModel) ProviderName() string {
	return "openai"
}

func (m *LanguageModel) ModelID() string {
	return m.modelID
}

func (m *LanguageModel) DefaultObjectGenerationMode() api.ObjectGenerationMode {
	return api.ObjectGenerationModeJSON
}

func (m *LanguageModel) SupportsImageURLs() bool {
	return true
}

func (m *LanguageModel) SupportsStructuredOutputs() bool {
	return false
}

func (m *LanguageModel) SupportsURL(u *url.URL) bool {
	return true
}

func (m *LanguageModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (api.Response, error) {
	params, warnings, err := codec.Encode(m.modelID, prompt, opts)
	if err != nil {
		return api.Response{}, err
	}

	openaiResponse, err := m.client.Responses.New(ctx, params)
	if err != nil {
		return api.Response{}, err
	}

	response, err := codec.DecodeResponse(openaiResponse)
	if err != nil {
		return api.Response{}, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (api.StreamResponse, error) {
	params, warnings, err := codec.Encode(m.modelID, prompt, opts)
	if err != nil {
		return api.StreamResponse{}, err
	}

	stream := m.client.Responses.NewStreaming(ctx, params)
	response, err := codec.DecodeStream(stream)
	if err != nil {
		return api.StreamResponse{}, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}
