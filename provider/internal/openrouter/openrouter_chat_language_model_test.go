package openrouter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/aitesting"
	"go.jetify.com/ai/api"
	"go.jetify.com/ai/provider/internal/openrouter/client"
	"go.jetify.com/ai/provider/internal/openrouter/codec"
	"go.jetify.com/pkg/httpmock"
)

var testPrompt = []api.Message{
	&api.UserMessage{
		Content: api.ContentFromText("Hello"),
	},
}

var testLogprobs = &client.LogProbs{
	Content: []client.LogProb{
		{
			Token:   "Hello",
			LogProb: -0.0009994634,
			TopLogProbs: []client.TopLogProb{
				{Token: "Hello", LogProb: -0.0009994634},
			},
		},
		{
			Token:   "!",
			LogProb: -0.13410144,
			TopLogProbs: []client.TopLogProb{
				{Token: "!", LogProb: -0.13410144},
			},
		},
		{
			Token:   " How",
			LogProb: -0.0009250381,
			TopLogProbs: []client.TopLogProb{
				{Token: " How", LogProb: -0.0009250381},
			},
		},
		{
			Token:   " can",
			LogProb: -0.047709424,
			TopLogProbs: []client.TopLogProb{
				{Token: " can", LogProb: -0.047709424},
			},
		},
		{
			Token:   " I",
			LogProb: -0.000009014684,
			TopLogProbs: []client.TopLogProb{
				{Token: " I", LogProb: -0.000009014684},
			},
		},
		{
			Token:   " assist",
			LogProb: -0.009125131,
			TopLogProbs: []client.TopLogProb{
				{Token: " assist", LogProb: -0.009125131},
			},
		},
		{
			Token:   " you",
			LogProb: -0.0000066306106,
			TopLogProbs: []client.TopLogProb{
				{Token: " you", LogProb: -0.0000066306106},
			},
		},
		{
			Token:   " today",
			LogProb: -0.00011093382,
			TopLogProbs: []client.TopLogProb{
				{Token: " today", LogProb: -0.00011093382},
			},
		},
		{
			Token:   "?",
			LogProb: -0.00004596782,
			TopLogProbs: []client.TopLogProb{
				{Token: "?", LogProb: -0.00004596782},
			},
		},
	},
}

func TestDoGenerate(t *testing.T) {
	t.Skip("skipping test")
	defaultModel := "anthropic/claude-3.5-sonnet"

	tests := []struct {
		name         string
		modelID      string               // Optional: override default model
		settings     *client.ChatSettings // Optional: custom settings
		expectedReq  httpmock.Request
		mockResp     responseValues
		expectedResp api.Response
		wantErr      bool
	}{
		{
			name: "should extract text response",
			mockResp: responseValues{
				Content: "Hello, World!",
			},
			expectedResp: api.Response{
				Text: "Hello, World!",
			},
		},
		{
			name: "should extract usage",
			mockResp: responseValues{
				Content: "",
				Usage: &client.Usage{
					PromptTokens:     20,
					TotalTokens:      25,
					CompletionTokens: 5,
				},
			},
			expectedResp: api.Response{
				Usage: api.Usage{
					InputTokens:  20,
					OutputTokens: 5,
					TotalTokens:  25,
				},
			},
		},
		{
			name: "should extract logprobs",
			mockResp: responseValues{
				LogProbs: testLogprobs,
			},
			expectedResp: api.Response{
				LogProbs: codec.DecodeLogProbs(testLogprobs),
			},
		},
		{
			name: "should extract finish reason",
			mockResp: responseValues{
				Content:      "",
				FinishReason: client.FinishReasonStop,
			},
			expectedResp: api.Response{
				FinishReason: api.FinishReasonStop,
			},
		},
		{
			name: "should support unknown finish reason",
			mockResp: responseValues{
				Content:      "",
				FinishReason: "eos",
			},
			expectedResp: api.Response{
				FinishReason: api.FinishReasonUnknown,
			},
		},
		{
			name: "should pass the model and messages",
			mockResp: responseValues{
				Content: "",
			},
			expectedReq: httpmock.Request{
				Body: map[string]interface{}{
					"model": "anthropic/claude-3.5-sonnet",
					"messages": []map[string]string{
						{
							"role":    "user",
							"content": "Hello",
						},
					},
				},
			},
		},
		{
			name: "should pass the models array when provided",
			settings: &client.ChatSettings{
				Models: []string{
					"anthropic/claude-2",
					"gryphe/mythomax-l2-13b",
				},
			},
			mockResp: responseValues{
				Content: "",
			},
			expectedReq: httpmock.Request{
				Body: map[string]interface{}{
					"model": "anthropic/claude-3.5-sonnet",
					"models": []string{
						"anthropic/claude-2",
						"gryphe/mythomax-l2-13b",
					},
					"messages": []map[string]string{
						{
							"role":    "user",
							"content": "Hello",
						},
					},
				},
			},
		},
		{
			name:    "should pass settings",
			modelID: "openai/gpt-3.5-turbo",
			settings: &client.ChatSettings{
				LogitBias: map[int]float64{50256: -100},
				Logprobs: &client.LogprobSettings{
					Enabled: true,
					TopK:    2,
				},
				ParallelToolCalls: boolPtr(false),
				User:              stringPtr("test-user-id"),
			},
			mockResp: responseValues{
				Content: "",
			},
			expectedReq: httpmock.Request{
				Body: map[string]interface{}{
					"model": "openai/gpt-3.5-turbo",
					"messages": []map[string]string{
						{
							"role":    "user",
							"content": "Hello",
						},
					},
					"logprobs":            2, // LogprobSettings with TopK=2 serializes to number
					"logit_bias":          map[string]float64{"50256": -100},
					"parallel_tool_calls": false,
					"user":                "test-user-id",
				},
			},
		},
	}
	// TODO: add tests for "should pass tools and toolChoice"
	// and "should pass headers"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultReq := httpmock.Request{
				Method: http.MethodPost,
				Path:   "/api/v1/chat/completions",
			}

			server := httpmock.NewServer(t, []httpmock.Exchange{
				{
					Request:  httpmock.MergeRequests(defaultReq, tt.expectedReq),
					Response: newResponse(tt.mockResp),
				},
			})
			defer server.Close()

			provider := NewOpenRouterProvider(
				server.BaseURL(),
				"test-api-key",
			)

			modelID := defaultModel
			if tt.modelID != "" {
				modelID = tt.modelID
			}

			model := NewOpenRouterChatLanguageModel(provider, modelID, tt.settings)

			got, err := model.DoGenerate(t.Context(), testPrompt, api.CallOptions{})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			aitesting.ResponseContains(t, tt.expectedResp, got)
		})
	}
}

type responseValues struct {
	Content      string            `json:"content"`
	Usage        *client.Usage     `json:"usage,omitempty"`
	LogProbs     *client.LogProbs  `json:"logprobs,omitempty"`
	FinishReason string            `json:"finish_reason,omitempty"`
	Headers      map[string]string // Response headers
}

func defaultClientResponse() client.Response {
	return client.Response{
		ID:                "chatcmpl-95ZTZkhr0mHNKqerQfiwkuox3PHAd",
		Object:            "chat.completion",
		Created:           1711115037,
		Model:             "gpt-3.5-turbo-0125",
		SystemFingerprint: "fp_3bc1b5746c",
		Usage: &client.Usage{
			PromptTokens:     4,
			TotalTokens:      34,
			CompletionTokens: 30,
		},
		Choices: []client.Choice{
			{
				Index:        0,
				Message:      client.AssistantMessage{},
				FinishReason: "stop",
			},
		},
	}
}

func newResponse(values responseValues) httpmock.Response {
	resp := defaultClientResponse()

	// Override defaults with provided values
	if values.Content != "" {
		resp.Choices[0].Message.Content = values.Content
	}
	if values.Usage != nil {
		resp.Usage = values.Usage
	}
	if values.LogProbs != nil {
		resp.Choices[0].LogProbs = *values.LogProbs
	}
	if values.FinishReason != "" {
		resp.Choices[0].FinishReason = values.FinishReason
	}

	headers := map[string]string{}

	// Add any custom headers
	for k, v := range values.Headers {
		headers[k] = v
	}

	return httpmock.Response{
		StatusCode: http.StatusOK,
		Body:       resp,
		Headers:    headers,
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
