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
				Text: "hello world",
			},
			contains: api.Response{
				Text: "hello world",
			},
			wantFail: false,
		},
		{
			name: "mismatched text",
			expected: api.Response{
				Text: "hello world",
			},
			contains: api.Response{
				Text: "goodbye world",
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
				Text: "hello",
			},
			contains: api.Response{
				Text: "hello",
				ResponseInfo: &api.ResponseInfo{
					ModelID: "different-model",
				},
				Sources:  []api.Source{},
				LogProbs: api.LogProbs{},
			},
			wantFail: false,
		},
		{
			name: "matching tool calls",
			expected: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
						ToolCallID: "test-id",
						ToolName:   "test",
						Args:       json.RawMessage(`"args"`),
					},
				},
			},
			contains: api.Response{
				ToolCalls: []api.ToolCallBlock{
					{
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
				Files: []api.FileBlock{
					{
						Filename:  "test.txt",
						Data:      []byte("test content"),
						MediaType: "text/plain",
					},
				},
			},
			contains: api.Response{
				Files: []api.FileBlock{
					{
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
				Files: []api.FileBlock{
					{
						Filename:  "test.txt",
						Data:      []byte("test content"),
						MediaType: "text/plain",
					},
				},
			},
			contains: api.Response{
				Files: []api.FileBlock{
					{
						Filename:  "different.txt",
						Data:      []byte("different content"),
						MediaType: "text/plain",
					},
				},
			},
			wantFail: true,
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
