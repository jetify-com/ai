package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogprobSettings_Marshal(t *testing.T) {
	tests := []struct {
		name     string
		logprobs *LogprobSettings
		want     string
	}{
		{
			name:     "nil logprobs",
			logprobs: nil,
			want:     "null",
		},
		{
			name: "enabled true",
			logprobs: &LogprobSettings{
				Enabled: true,
			},
			want: "true",
		},
		{
			name: "enabled false",
			logprobs: &LogprobSettings{
				Enabled: false,
			},
			want: "false",
		},
		{
			name: "top K",
			logprobs: &LogprobSettings{
				Enabled: true,
				TopK:    5,
			},
			want: "5",
		},
		{
			name: "disabled with zero topK",
			logprobs: &LogprobSettings{
				Enabled: false,
				TopK:    0,
			},
			want: "false",
		},
		{
			name: "enabled with zero topK",
			logprobs: &LogprobSettings{
				Enabled: true,
				TopK:    0,
			},
			want: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			got, err := json.Marshal(tt.logprobs)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))

			// Test unmarshaling
			var l LogprobSettings
			err = json.Unmarshal([]byte(tt.want), &l)
			assert.NoError(t, err)
			if tt.logprobs != nil {
				assert.Equal(t, tt.logprobs.Enabled, l.Enabled)
				assert.Equal(t, tt.logprobs.TopK, l.TopK)
			}
		})
	}
}

func TestLogprobSettings_UnmarshalErrors(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		errSubstr string
	}{
		{
			name:      "invalid string",
			input:     `"invalid"`,
			errSubstr: "logprobs must be boolean or number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l LogprobSettings
			err := json.Unmarshal([]byte(tt.input), &l)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errSubstr)
		})
	}
}
