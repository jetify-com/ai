package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/anthropic/codec"
)

// ModelOption is a function type that modifies a LanguageModel.
type ModelOption func(*LanguageModel)

// WithClient returns a ModelOption that sets the client.
func WithClient(client anthropic.Client) ModelOption {
	// TODO: Instead of only supporting an anthropic.Client, we can "flatten"
	// the options supported by the Anthropic SDK.
	return func(m *LanguageModel) {
		m.client = client
	}
}

// LanguageModel represents an Anthropic language model.
type LanguageModel struct {
	modelID string
	client  anthropic.Client
}

var _ api.LanguageModel = &LanguageModel{}

// NewLanguageModel creates a new Anthropic language model.
func NewLanguageModel(modelID string, opts ...ModelOption) *LanguageModel {
	// Create model with default settings
	model := &LanguageModel{
		modelID: modelID,
		client:  anthropic.NewClient(), // Default client
	}

	// Apply options
	for _, opt := range opts {
		opt(model)
	}

	return model
}

func (m *LanguageModel) ProviderName() string {
	return ProviderName
}

func (m *LanguageModel) ModelID() string {
	return m.modelID
}

func (m *LanguageModel) SupportedUrls() []api.SupportedURL {
	// TODO: Make configurable via the constructor.
	return []api.SupportedURL{
		{
			MediaType: "image/*",
			URLPatterns: []string{
				"^https?://.*",
			},
		},
	}
}

func (m *LanguageModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.Response, error) {
	params, warnings, err := codec.EncodeParams(m.modelID, prompt, opts)
	if err != nil {
		return nil, err
	}

	message, err := m.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return nil, err
	}

	response, err := codec.DecodeResponse(message)
	if err != nil {
		return nil, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (*api.StreamResponse, error) {
	return nil, api.NewUnsupportedFunctionalityError("streaming generation", "")
}
