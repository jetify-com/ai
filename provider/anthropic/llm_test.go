package anthropic

import (
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.jetify.com/ai/api"
	"go.jetify.com/pkg/httpmock"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name        string
		prompt      []api.Message
		exchanges   []httpmock.Exchange
		expectError bool
	}{
		{
			name: "successful generation with user message",
			prompt: []api.Message{
				&api.UserMessage{
					Content: api.ContentFromText("Hello, how are you?"),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/v1/messages",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body: &anthropic.Message{
							Content: []anthropic.ContentBlock{
								{
									Text: "I'm doing well, thank you for asking!",
									Type: anthropic.ContentBlockTypeText,
								},
							},
							Role: anthropic.MessageRoleAssistant,
						},
					},
				},
			},
		},
		{
			name: "successful generation with system message",
			prompt: []api.Message{
				&api.SystemMessage{Content: "You are a helpful assistant"},
				&api.UserMessage{
					Content: api.ContentFromText("What's 2+2?"),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/v1/messages",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusOK,
						Body: &anthropic.Message{
							Content: []anthropic.ContentBlock{
								{
									Text: "4",
									Type: anthropic.ContentBlockTypeText,
								},
							},
							Role: anthropic.MessageRoleAssistant,
						},
					},
				},
			},
		},
		{
			name: "api error",
			prompt: []api.Message{
				&api.UserMessage{
					Content: api.ContentFromText("Hello"),
				},
			},
			exchanges: []httpmock.Exchange{
				{
					Request: httpmock.Request{
						Method: http.MethodPost,
						Path:   "/v1/messages",
					},
					Response: httpmock.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       map[string]interface{}{"error": "internal server error"},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httpmock.NewServer(t, tt.exchanges)
			defer server.Close()

			// Create client with mock server URL and test API key
			client := anthropic.NewClient(
				option.WithBaseURL(server.BaseURL()),
				option.WithAPIKey("test-key"),
				option.WithMaxRetries(0), // Disable retries
			)

			// Create model with mocked client
			model := NewLanguageModel("claude-3", WithClient(client))

			// Call Generate with empty CallOptions
			resp, err := model.Generate(t.Context(), tt.prompt, api.CallOptions{})

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			// For successful cases, verify response content
			message := tt.exchanges[0].Response.Body.(*anthropic.Message)
			if len(message.Content) > 0 {
				assert.Equal(t, message.Content[0].Text, resp.Text)
			}
		})
	}
}
