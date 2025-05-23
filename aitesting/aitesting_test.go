package aitesting

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

// mockT implements the T interface for testing
type mockT struct {
	failed    bool
	failNowed bool
	errors    []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	m.errors = append(m.errors, format)
}

func (m *mockT) FailNow() {
	m.failNowed = true
}

func TestResponseContains(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		expected api.Response
		contains api.Response
		wantFail bool
	}{
		{
			name: "matching text",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
			wantFail: false,
		},
		{
			name: "mismatched text",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello world"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "goodbye world"},
				},
			},
			wantFail: true,
		},
		{
			name: "matching usage",
			expected: api.Response{
				Usage: api.Usage{
					InputTokens:       10,
					OutputTokens:      20,
					TotalTokens:       30,
					ReasoningTokens:   10,
					CachedInputTokens: 5,
				},
			},
			contains: api.Response{
				Usage: api.Usage{
					InputTokens:       10,
					OutputTokens:      20,
					TotalTokens:       30,
					ReasoningTokens:   10,
					CachedInputTokens: 5,
				},
			},
			wantFail: false,
		},
		{
			name: "exact usage check",
			expected: api.Response{
				Usage: api.Usage{
					InputTokens: 10,
					// All other fields should be 0
				},
			},
			contains: api.Response{
				Usage: api.Usage{
					InputTokens: 10,
					// All other fields should be 0
				},
			},
			wantFail: false,
		},
		{
			name: "mismatched usage",
			expected: api.Response{
				Usage: api.Usage{
					InputTokens: 10,
				},
			},
			contains: api.Response{
				Usage: api.Usage{
					InputTokens:  10,
					OutputTokens: 20, // This will cause the test to fail since OutputTokens doesn't match
				},
			},
			wantFail: true,
		},
		{
			name: "matching response info",
			expected: api.Response{
				ResponseInfo: &api.ResponseInfo{
					ID:        "test-id",
					ModelID:   "test-model",
					Timestamp: now,
				},
			},
			contains: api.Response{
				ResponseInfo: &api.ResponseInfo{
					ID:        "test-id",
					ModelID:   "test-model",
					Timestamp: now,
				},
			},
			wantFail: false,
		},
		{
			name: "nil response info when expected",
			expected: api.Response{
				ResponseInfo: &api.ResponseInfo{
					ID: "test-id",
				},
			},
			contains: api.Response{
				ResponseInfo: nil,
			},
			wantFail: true,
		},
		{
			name: "ignores unset fields",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
					&api.SourceBlock{ID: "src-1", URL: "http://example.com"},
				},
				ResponseInfo: &api.ResponseInfo{
					ModelID: "different-model",
				},
			},
			wantFail: false,
		},
		{
			name: "matching tool calls",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "test-id",
						ToolName:   "test",
						Args:       json.RawMessage(`"args"`),
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "test-id",
						ToolName:   "test",
						Args:       json.RawMessage(`"args"`),
					},
				},
			},
			wantFail: false,
		},
		{
			name: "matching files",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Filename:  "test.txt",
						Data:      []byte("test content"),
						MediaType: "text/plain",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Filename:  "test.txt",
						Data:      []byte("test content"),
						MediaType: "text/plain",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "mismatched files",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Filename:  "test.txt",
						Data:      []byte("test content"),
						MediaType: "text/plain",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.FileBlock{
						Filename:  "different.txt",
						Data:      []byte("different content"),
						MediaType: "text/plain",
					},
				},
			},
			wantFail: true,
		},
		{
			name: "matching image blocks",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL:       "https://example.com/image.jpg",
						MediaType: "image/jpeg",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ImageBlock{
						URL:       "https://example.com/image.jpg",
						MediaType: "image/jpeg",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "matching reasoning blocks",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{
						Text:      "Let me think about this...",
						Signature: "sig123",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ReasoningBlock{
						Text:      "Let me think about this...",
						Signature: "sig123",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "matching redacted reasoning blocks",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.RedactedReasoningBlock{
						Data: "redacted_data_123",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.RedactedReasoningBlock{
						Data: "redacted_data_123",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "matching source blocks",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.SourceBlock{
						ID:    "src-1",
						URL:   "https://example.com",
						Title: "Example Source",
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.SourceBlock{
						ID:    "src-1",
						URL:   "https://example.com",
						Title: "Example Source",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "partial field matching - tool call only name",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolName: "test",
						// ID and Args are not set, so they won't be compared
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{
						ToolCallID: "some-id",
						ToolName:   "test",
						Args:       json.RawMessage(`{"key": "value"}`),
					},
				},
			},
			wantFail: false,
		},
		{
			name: "partial field matching - text block empty text",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{
						// Text is empty, so it won't be compared
					},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{
						Text: "any text here",
					},
				},
			},
			wantFail: false,
		},
		{
			name: "mixed content blocks",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
					&api.ToolCallBlock{ToolName: "test"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
					&api.ImageBlock{URL: "https://example.com/img.jpg"},
					&api.ToolCallBlock{
						ToolCallID: "id-123",
						ToolName:   "test",
						Args:       json.RawMessage(`{}`),
					},
					&api.SourceBlock{ID: "src-1"},
				},
			},
			wantFail: true,
		},
		{
			name: "wrong content block type",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ImageBlock{URL: "https://example.com/img.jpg"},
				},
			},
			wantFail: true,
		},
		{
			name: "empty content blocks",
			expected: api.Response{
				Content: []api.ContentBlock{},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "hello"},
				},
			},
			wantFail: false,
		},
		{
			name: "content blocks must match in order",
			expected: api.Response{
				Content: []api.ContentBlock{
					&api.TextBlock{Text: "first"},
					&api.ToolCallBlock{ToolName: "test"},
				},
			},
			contains: api.Response{
				Content: []api.ContentBlock{
					&api.ToolCallBlock{ToolName: "test"},
					&api.TextBlock{Text: "first"},
				},
			},
			wantFail: true, // Should fail because order is wrong
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockT{}
			ResponseContains(mock, tt.expected, tt.contains)
			assert.Equal(t, tt.wantFail, mock.failed, "ResponseContains() failed = %v, want %v", mock.failed, tt.wantFail)
		})
	}
}
