package anthropic

import (
	"context"
	"net/url"

	"github.com/anthropics/anthropic-sdk-go"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/anthropic/codec"
)

// ModelOption is a function type that modifies a LanguageModel.
type ModelOption func(*LanguageModel)

// WithClient returns a ModelOption that sets the client.
func WithClient(client *anthropic.Client) ModelOption {
	// TODO: Instead of only supporting an anthropic.Client, we can "flatten"
	// the options supported by the Anthropic SDK.
	return func(m *LanguageModel) {
		m.client = client
	}
}

// LanguageModel represents an Anthropic language model.
type LanguageModel struct {
	modelID string
	client  *anthropic.Client
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

func (m *LanguageModel) SpecificationVersion() string {
	return "v1"
}

func (m *LanguageModel) ProviderName() string {
	return ProviderName
}

func (m *LanguageModel) ModelID() string {
	return m.modelID
}

func (m *LanguageModel) DefaultObjectGenerationMode() api.ObjectGenerationMode {
	return api.ObjectGenerationModeTool // Anthropic models support tool mode by default
}

func (m *LanguageModel) SupportsImageURLs() bool {
	return true // Claude 3 models support image URLs
}

func (m *LanguageModel) SupportsStructuredOutputs() bool {
	return true // Claude models support structured JSON outputs
}

func (m *LanguageModel) SupportsURL(u *url.URL) bool {
	// TODO: Double check if we should only return true for a subset of URLs
	return true // Anthropic models support URLs
}

func (m *LanguageModel) Generate(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (api.Response, error) {
	params, warnings, err := codec.EncodeParams(prompt, opts)
	if err != nil {
		return api.Response{}, err
	}

	message, err := m.client.Beta.Messages.New(ctx, params)
	if err != nil {
		return api.Response{}, err
	}

	response, err := codec.DecodeResponse(message)
	if err != nil {
		return api.Response{}, err
	}

	response.Warnings = append(response.Warnings, warnings...)
	return response, nil
}

func (m *LanguageModel) Stream(
	ctx context.Context, prompt []api.Message, opts api.CallOptions,
) (api.StreamResponse, error) {
	return api.StreamResponse{}, api.NewUnsupportedFunctionalityError("streaming generation", "")
}
