package aitesting

import (
	"github.com/stretchr/testify/assert"
	"go.jetify.com/ai/api"
)

// T is an interface that captures the testing.T methods we need
type T interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

// ResponseContains checks if a response contains all the expected fields.
// It only compares fields that are set in the expected response, ignoring unset fields.
// This allows for partial matching of responses where only specific fields need to be verified.
//
// Parameters:
//   - testingT: The testing.T instance for reporting test failures
//   - expected: The expected Response containing the fields to check
//   - contains: The actual Response to check against the expected values
//
// Example:
//
//	expected := api.Response{
//	    Text: "Hello world",
//	    Usage: api.Usage{PromptTokens: 10},
//	}
//	ResponseContains(t, expected, actualResponse)
//
// The above example will verify that actualResponse has the text "Hello world"
// and 10 prompt tokens, while ignoring all other fields.
func ResponseContains(testingT T, expected api.Response, contains api.Response) {
	// Compare text if set
	if expected.Text != "" {
		assert.Equal(testingT, expected.Text, contains.Text, "Text mismatch")
	}

	// Compare reasoning if set
	if len(expected.Reasoning) > 0 {
		assert.Equal(testingT, expected.Reasoning, contains.Reasoning, "Reasoning mismatch")
	}

	// Compare files if set
	if len(expected.Files) > 0 {
		assert.Equal(testingT, expected.Files, contains.Files, "Files mismatch")
	}

	// Compare tool calls if set
	if len(expected.ToolCalls) > 0 {
		assert.Equal(testingT, expected.ToolCalls, contains.ToolCalls, "ToolCalls mismatch")
	}

	// Compare usage if set
	if !expected.Usage.IsZero() {
		assert.Equal(testingT, expected.Usage, contains.Usage, "Usage mismatch")
	}

	// Compare finish reason if set
	if expected.FinishReason != "" {
		assert.Equal(testingT, expected.FinishReason, contains.FinishReason, "FinishReason mismatch")
	}

	// Compare request info if set
	if expected.RequestInfo != nil {
		assert.NotNil(testingT, contains.RequestInfo, "RequestInfo should not be nil")
		if expected.RequestInfo.Body != nil {
			assert.Equal(testingT, expected.RequestInfo.Body, contains.RequestInfo.Body, "RequestInfo.Body mismatch")
		}
	}

	// Compare response info if set
	if expected.ResponseInfo != nil {
		if contains.ResponseInfo == nil {
			assert.Fail(testingT, "ResponseInfo should not be nil")
			return
		}
		if expected.ResponseInfo.ID != "" {
			assert.Equal(testingT, expected.ResponseInfo.ID, contains.ResponseInfo.ID, "ResponseInfo.ID mismatch")
		}
		if expected.ResponseInfo.ModelID != "" {
			assert.Equal(testingT, expected.ResponseInfo.ModelID, contains.ResponseInfo.ModelID, "ResponseInfo.ModelID mismatch")
		}
		if !expected.ResponseInfo.Timestamp.IsZero() {
			assert.Equal(testingT, expected.ResponseInfo.Timestamp, contains.ResponseInfo.Timestamp, "ResponseInfo.Timestamp mismatch")
		}
	}

	// Compare warnings if set
	if len(expected.Warnings) > 0 {
		assert.Equal(testingT, expected.Warnings, contains.Warnings, "Warnings mismatch")
	}

	// Compare provider metadata if set
	if !expected.ProviderMetadata.IsZero() {
		assert.Equal(testingT, expected.ProviderMetadata, contains.ProviderMetadata, "ProviderMetadata mismatch")
	}

	// Compare sources if set
	if len(expected.Sources) > 0 {
		assert.Equal(testingT, expected.Sources, contains.Sources, "Sources mismatch")
	}

	// Compare log probs if set
	if expected.LogProbs != nil {
		assert.NotNil(testingT, contains.LogProbs, "LogProbs should not be nil")
		assert.Equal(testingT, expected.LogProbs, contains.LogProbs, "LogProbs mismatch")
	}
}
