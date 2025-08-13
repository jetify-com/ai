package aitesting

import (
	"bytes"

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
// For content blocks, it checks that each expected content block has a corresponding
// match in the actual response content (order and additional blocks don't matter).
// Within each content block, only fields that are set in the expected block are compared,
// allowing for partial content block matching as well.
//
// Parameters:
//   - testingT: The testing.T instance for reporting test failures
//   - expected: The expected Response containing the fields to check
//   - contains: The actual Response to check against the expected values
//
// Examples:
//
//	// Check for any text block with specific text
//	expected := api.Response{
//	    Content: []api.ContentBlock{
//	        &api.TextBlock{Text: "Hello world"},
//	    },
//	    Usage: api.Usage{InputTokens: 10},
//	}
//	ResponseContains(t, expected, actualResponse)
//
//	// Check for a tool call with specific name (ignoring args and ID)
//	expected := api.Response{
//	    Content: []api.ContentBlock{
//	        &api.ToolCallBlock{ToolName: "get_weather"},
//	    },
//	}
//	ResponseContains(t, expected, actualResponse)
//
//	// Check for any image block (ignoring all fields)
//	expected := api.Response{
//	    Content: []api.ContentBlock{
//	        &api.ImageBlock{}, // matches any image block
//	    },
//	}
//	ResponseContains(t, expected, actualResponse)
//
// The first example verifies that actualResponse has a text block with "Hello world"
// and 10 input tokens, while ignoring all other fields and content blocks.
// The second example only checks for the presence of a tool call named "get_weather".
// The third example checks for the presence of any image block regardless of its content.
func ResponseContains(testingT T, expected, contains api.Response) {
	// Compare content blocks if set
	if len(expected.Content) > 0 {
		// Check that we have enough blocks in the actual response
		if len(contains.Content) < len(expected.Content) {
			assert.Fail(testingT, "Not enough content blocks",
				"Expected %d content blocks, but actual response only has %d", len(expected.Content), len(contains.Content))
			return
		}

		// Compare blocks in order
		for i, expectedBlock := range expected.Content {
			actualBlock := contains.Content[i]
			if !contentBlocksEqual(testingT, expectedBlock, actualBlock) {
				assert.Fail(testingT, "Content block mismatch",
					"Content block at index %d does not match. Expected: %T, Actual: %T", i, expectedBlock, actualBlock)
			}
		}
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
}

// contentBlocksEqual compares two content blocks for equality
// Only compares fields that are set in the expected block (contains semantics)
func contentBlocksEqual(testingT T, expected, actual api.ContentBlock) bool {
	// First check if they're the same type
	if expected.Type() != actual.Type() {
		return false
	}

	// Then compare based on the specific type using helper functions
	switch expectedBlock := expected.(type) {
	case *api.TextBlock:
		if actualBlock, ok := actual.(*api.TextBlock); ok {
			return textBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.ReasoningBlock:
		if actualBlock, ok := actual.(*api.ReasoningBlock); ok {
			return reasoningBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.RedactedReasoningBlock:
		if actualBlock, ok := actual.(*api.RedactedReasoningBlock); ok {
			return redactedReasoningBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.ImageBlock:
		if actualBlock, ok := actual.(*api.ImageBlock); ok {
			return imageBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.FileBlock:
		if actualBlock, ok := actual.(*api.FileBlock); ok {
			return fileBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.ToolCallBlock:
		if actualBlock, ok := actual.(*api.ToolCallBlock); ok {
			return toolCallBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	case *api.SourceBlock:
		if actualBlock, ok := actual.(*api.SourceBlock); ok {
			return sourceBlocksEqual(testingT, expectedBlock, actualBlock)
		}
	}

	return false
}

// textBlocksEqual compares two text blocks using contains semantics
func textBlocksEqual(testingT T, expected, actual *api.TextBlock) bool {
	allMatch := true
	// Only check text if it's set in expected
	if expected.Text != "" {
		if !assert.Equal(testingT, expected.Text, actual.Text, "TextBlock.Text mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// reasoningBlocksEqual compares two reasoning blocks using contains semantics
func reasoningBlocksEqual(testingT T, expected, actual *api.ReasoningBlock) bool {
	allMatch := true
	// Only check text if it's set in expected
	if expected.Text != "" {
		if !assert.Equal(testingT, expected.Text, actual.Text, "ReasoningBlock.Text mismatch") {
			allMatch = false
		}
	}
	// Only check signature if it's set in expected
	if expected.Signature != "" {
		if !assert.Equal(testingT, expected.Signature, actual.Signature, "ReasoningBlock.Signature mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// redactedReasoningBlocksEqual compares two redacted reasoning blocks using contains semantics
func redactedReasoningBlocksEqual(testingT T, expected, actual *api.RedactedReasoningBlock) bool {
	allMatch := true
	// Only check data if it's set in expected
	if expected.Data != "" {
		if !assert.Equal(testingT, expected.Data, actual.Data, "RedactedReasoningBlock.Data mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// imageBlocksEqual compares two image blocks using contains semantics
func imageBlocksEqual(testingT T, expected, actual *api.ImageBlock) bool {
	allMatch := true
	// Only check URL if it's set in expected
	if expected.URL != "" {
		if !assert.Equal(testingT, expected.URL, actual.URL, "ImageBlock.URL mismatch") {
			allMatch = false
		}
	}
	// Only check data if it's set in expected
	if len(expected.Data) > 0 {
		if !assert.True(testingT, bytes.Equal(expected.Data, actual.Data), "ImageBlock.Data mismatch") {
			allMatch = false
		}
	}
	// Only check media type if it's set in expected
	if expected.MediaType != "" {
		if !assert.Equal(testingT, expected.MediaType, actual.MediaType, "ImageBlock.MediaType mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// fileBlocksEqual compares two file blocks using contains semantics
func fileBlocksEqual(testingT T, expected, actual *api.FileBlock) bool {
	allMatch := true
	// Only check filename if it's set in expected
	if expected.Filename != "" {
		if !assert.Equal(testingT, expected.Filename, actual.Filename, "FileBlock.Filename mismatch") {
			allMatch = false
		}
	}
	// Only check URL if it's set in expected
	if expected.URL != "" {
		if !assert.Equal(testingT, expected.URL, actual.URL, "FileBlock.URL mismatch") {
			allMatch = false
		}
	}
	// Only check data if it's set in expected
	if len(expected.Data) > 0 {
		if !assert.True(testingT, bytes.Equal(expected.Data, actual.Data), "FileBlock.Data mismatch") {
			allMatch = false
		}
	}
	// Only check media type if it's set in expected
	if expected.MediaType != "" {
		if !assert.Equal(testingT, expected.MediaType, actual.MediaType, "FileBlock.MediaType mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// toolCallBlocksEqual compares two tool call blocks using contains semantics
func toolCallBlocksEqual(testingT T, expected, actual *api.ToolCallBlock) bool {
	allMatch := true
	// Only check tool call ID if it's set in expected
	if expected.ToolCallID != "" {
		if !assert.Equal(testingT, expected.ToolCallID, actual.ToolCallID, "ToolCallBlock.ToolCallID mismatch") {
			allMatch = false
		}
	}
	// Only check tool name if it's set in expected
	if expected.ToolName != "" {
		if !assert.Equal(testingT, expected.ToolName, actual.ToolName, "ToolCallBlock.ToolName mismatch") {
			allMatch = false
		}
	}
	// Only check args if it's set in expected
	if len(expected.Args) > 0 {
		if !assert.Equal(testingT, string(expected.Args), string(actual.Args), "ToolCallBlock.Args mismatch") {
			allMatch = false
		}
	}
	return allMatch
}

// sourceBlocksEqual compares two source blocks using contains semantics
func sourceBlocksEqual(testingT T, expected, actual *api.SourceBlock) bool {
	allMatch := true
	// Only check ID if it's set in expected
	if expected.ID != "" {
		if !assert.Equal(testingT, expected.ID, actual.ID, "SourceBlock.ID mismatch") {
			allMatch = false
		}
	}
	// Only check URL if it's set in expected
	if expected.URL != "" {
		if !assert.Equal(testingT, expected.URL, actual.URL, "SourceBlock.URL mismatch") {
			allMatch = false
		}
	}
	// Only check title if it's set in expected
	if expected.Title != "" {
		if !assert.Equal(testingT, expected.Title, actual.Title, "SourceBlock.Title mismatch") {
			allMatch = false
		}
	}
	return allMatch
}
