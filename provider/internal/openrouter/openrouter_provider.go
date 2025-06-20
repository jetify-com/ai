package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenRouterProvider represents the OpenRouter API provider.
type OpenRouterProvider struct {
	baseURL string
	apiKey  string
	client  *http.Client
	headers http.Header
}

// NewOpenRouterProvider creates a new OpenRouter provider.
func NewOpenRouterProvider(baseURL string, apiKey string, opts ...ProviderOption) *OpenRouterProvider {
	p := &OpenRouterProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  http.DefaultClient,
		headers: make(http.Header),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// ProviderOption configures the OpenRouter provider.
type ProviderOption func(*OpenRouterProvider)

// WithClient sets a custom HTTP client.
func WithClient(client *http.Client) ProviderOption {
	return func(p *OpenRouterProvider) {
		p.client = client
	}
}

// WithHeaders sets custom headers for API requests.
func WithHeaders(headers http.Header) ProviderOption {
	return func(p *OpenRouterProvider) {
		for k, values := range headers {
			for _, v := range values {
				p.headers.Add(k, v)
			}
		}
	}
}

// doJSONRequest makes a JSON request to the OpenRouter API.
func (p *OpenRouterProvider) doJSONRequest(ctx context.Context, method, path string, body any, extraHeaders http.Header) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, p.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Set provider headers
	for k, values := range p.headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	// Set request-specific headers
	for k, values := range extraHeaders {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}
