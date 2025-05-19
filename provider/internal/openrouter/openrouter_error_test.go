package openrouter

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

func TestParseOpenRouterErrorJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    *openRouterErrorData
		wantErr bool
	}{
		{
			name: "valid error json",
			json: `{
				"error": {
					"message": "invalid request",
					"type": "invalid_request_error",
					"param": "model",
					"code": "model_not_found"
				}
			}`,
			want: &openRouterErrorData{
				Error: struct {
					Message string  `json:"message"`
					Type    string  `json:"type"`
					Param   any     `json:"param"`
					Code    *string `json:"code"`
				}{
					Message: "invalid request",
					Type:    "invalid_request_error",
					Param:   "model",
					Code:    strPtr("model_not_found"),
				},
			},
		},
		{
			name: "null fields",
			json: `{
				"error": {
					"message": "rate limited",
					"type": "rate_limit_error",
					"param": null,
					"code": null
				}
			}`,
			want: &openRouterErrorData{
				Error: struct {
					Message string  `json:"message"`
					Type    string  `json:"type"`
					Param   any     `json:"param"`
					Code    *string `json:"code"`
				}{
					Message: "rate limited",
					Type:    "rate_limit_error",
					Param:   nil,
					Code:    nil,
				},
			},
		},
		{
			name:    "invalid json",
			json:    `{"error": {`,
			wantErr: true,
		},
		{
			name: "empty json",
			json: `{}`,
			want: &openRouterErrorData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOpenRouterErrorJSON([]byte(tt.json))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOpenRouterFailedResponseHandler(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        string
		wantMessage string
	}{
		{
			name:       "valid error response",
			statusCode: http.StatusBadRequest,
			body: `{
				"error": {
					"message": "invalid request",
					"type": "invalid_request_error",
					"param": "model",
					"code": "model_not_found"
				}
			}`,
			wantMessage: "invalid request",
		},
		{
			name:       "rate limit error",
			statusCode: http.StatusTooManyRequests,
			body: `{
				"error": {
					"message": "rate limited",
					"type": "rate_limit_error",
					"param": null,
					"code": null
				}
			}`,
			wantMessage: "rate limited",
		},
		{
			name:        "invalid json falls back to status",
			statusCode:  http.StatusBadRequest,
			body:        `{"error": {`,
			wantMessage: "400 Bad Request",
		},
		{
			name:        "empty response falls back to status",
			statusCode:  http.StatusInternalServerError,
			body:        "",
			wantMessage: "500 Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			url := &url.URL{Scheme: "https", Host: "api.openrouter.ai", Path: "/api/v1/chat/completions"}
			req := &http.Request{URL: url}
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Status:     http.StatusText(tt.statusCode),
				Request:    req,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			// Call the handler
			err := OpenRouterFailedResponseHandler(resp, []byte(tt.body), nil)

			// Use errors.As instead of type assertion
			var apiErr *api.APICallError
			assert.True(t, errors.As(err, &apiErr))

			// Check the error details
			assert.Equal(t, tt.statusCode, apiErr.StatusCode)
			assert.Equal(t, url, apiErr.URL)
			assert.Equal(t, tt.wantMessage, apiErr.Error())

			// Check retryable status
			isRetryable := apiErr.IsRetryable()
			if tt.statusCode == 429 || tt.statusCode >= 500 {
				assert.True(t, isRetryable)
			} else {
				assert.False(t, isRetryable)
			}
		})
	}
}
