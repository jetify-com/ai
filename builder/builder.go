package builder

import (
	"fmt"
	"slices"

	"go.jetify.com/ai/api"
)

type responseStateEnum int

const (
	noState responseStateEnum = iota
	textState
	reasoningState
	toolCallState
	sourceState
	fileState
	finishedState
)

// ResponseBuilder builds a Response from stream events.
// The builder is not safe for concurrent use. If you need to use it from multiple goroutines,
// you must synchronize access to the builder yourself.
type ResponseBuilder struct {
	resp api.Response
	// Map of tool call ID to index in Content array (for parallel tool calls)
	toolCallIndices map[string]int
	// Error encountered during event processing
	err error

	// State tracking
	currentState responseStateEnum
	responseID   string

	// Usage tracking
	usage api.Usage
}

// NewResponseBuilder creates a new ResponseBuilder.
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		resp: api.Response{
			// Initialize with empty Content slice
			Content:      []api.ContentBlock{},
			Warnings:     []api.CallWarning{},
			FinishReason: api.FinishReason(""),
		},
		toolCallIndices: make(map[string]int),
		currentState:    noState,
		err:             nil,
	}
}

// AddEvent processes a stream event and adds it to the builder.
func (b *ResponseBuilder) AddEvent(event api.StreamEvent) error {
	if event == nil {
		return fmt.Errorf("cannot add nil event")
	}

	// Validate that we haven't already finished
	if b.currentState == finishedState {
		return fmt.Errorf("cannot add events after finish event")
	}

	// TODO:
	// - Validate that we have a response ID if this isn't the first event
	// - Validate that we have a model ID if this isn't the first event
	// The reason we haven't implemented these yet is because I'm not sure if all providers have those guarantees.

	switch evt := event.(type) {
	// Handle pointer types
	case *api.TextDeltaEvent:
		return b.addTextDelta(evt)
	case *api.ReasoningEvent:
		return b.addReasoning(evt)
	case *api.ReasoningSignatureEvent:
		return b.addReasoningSignature(evt)
	case *api.RedactedReasoningEvent:
		return b.addRedactedReasoning(evt)
	case *api.ToolCallEvent:
		return b.addToolCall(evt)
	case *api.ToolCallDeltaEvent:
		return b.addToolCallDelta(evt)
	case *api.SourceEvent:
		return b.addSource(evt)
	case *api.FileEvent:
		return b.addFile(evt)
	case *api.ResponseMetadataEvent:
		return b.addResponseMetadata(evt)
	case *api.FinishEvent:
		return b.addFinish(evt)
	case *api.ErrorEvent:
		return b.addError(evt)

	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

// addTextDelta adds a text delta event to the response.
func (b *ResponseBuilder) addTextDelta(e *api.TextDeltaEvent) error {
	// Only concatenate with last block if the last content block is a TextBlock
	if len(b.resp.Content) > 0 {
		if lastBlock, ok := b.resp.Content[len(b.resp.Content)-1].(*api.TextBlock); ok {
			// Append to existing text block
			lastBlock.Text += e.TextDelta
			return nil
		}
	}

	// Create new text block
	b.currentState = textState
	b.resp.Content = append(b.resp.Content, &api.TextBlock{
		Text: e.TextDelta,
	})
	return nil
}

// addReasoning adds a reasoning event to the response.
func (b *ResponseBuilder) addReasoning(e *api.ReasoningEvent) error {
	// Only concatenate with last block if the last content block is a ReasoningBlock
	if len(b.resp.Content) > 0 {
		if lastBlock, ok := b.resp.Content[len(b.resp.Content)-1].(*api.ReasoningBlock); ok {
			// Append to existing reasoning block
			lastBlock.Text += e.TextDelta
			return nil
		}
	}

	// Validate state transition
	if b.currentState != noState && b.currentState != reasoningState {
		return fmt.Errorf("invalid state transition: cannot add reasoning in state %v", b.currentState)
	}

	// Create new reasoning block
	b.currentState = reasoningState
	b.resp.Content = append(b.resp.Content, &api.ReasoningBlock{
		Text: e.TextDelta,
	})
	return nil
}

// addReasoningSignature adds a reasoning signature event to the response.
func (b *ResponseBuilder) addReasoningSignature(e *api.ReasoningSignatureEvent) error {
	if len(b.resp.Content) == 0 {
		return fmt.Errorf("cannot add reasoning signature: no content blocks exist")
	}

	if block, ok := b.resp.Content[len(b.resp.Content)-1].(*api.ReasoningBlock); ok {
		block.Signature = e.Signature
		return nil
	}
	return fmt.Errorf("cannot add reasoning signature: last block is not a reasoning block")
}

// addRedactedReasoning adds a redacted reasoning event to the response.
func (b *ResponseBuilder) addRedactedReasoning(e *api.RedactedReasoningEvent) error {
	// Validate state transition
	if b.currentState != noState && b.currentState != reasoningState {
		return fmt.Errorf("invalid state transition: cannot add redacted reasoning in state %v", b.currentState)
	}

	b.currentState = reasoningState
	b.resp.Content = append(b.resp.Content, &api.RedactedReasoningBlock{
		Data: e.Data,
	})
	return nil
}

// addToolCall adds a tool call event to the response.
func (b *ResponseBuilder) addToolCall(e *api.ToolCallEvent) error {
	// Check for duplicate tool call ID
	if _, exists := b.toolCallIndices[e.ToolCallID]; exists {
		return fmt.Errorf("duplicate tool call ID: %s", e.ToolCallID)
	}

	b.currentState = toolCallState

	// Create new tool call block
	toolCall := &api.ToolCallBlock{
		ToolCallID: e.ToolCallID,
		ToolName:   e.ToolName,
		Args:       slices.Clone(e.Args),
	}
	b.resp.Content = append(b.resp.Content, toolCall)
	// Store index in the map
	b.toolCallIndices[e.ToolCallID] = len(b.resp.Content) - 1

	return nil
}

// addToolCallDelta adds a tool call delta event to the response.
func (b *ResponseBuilder) addToolCallDelta(e *api.ToolCallDeltaEvent) error {
	b.currentState = toolCallState

	// Get or create the tool call block
	idx, exists := b.toolCallIndices[e.ToolCallID]
	if !exists {
		// Create new tool call block if it doesn't exist
		toolCall := &api.ToolCallBlock{
			ToolCallID: e.ToolCallID,
			ToolName:   e.ToolName,
			Args:       make([]byte, 0),
		}
		b.resp.Content = append(b.resp.Content, toolCall)
		idx = len(b.resp.Content) - 1
		b.toolCallIndices[e.ToolCallID] = idx
	}

	// Append the new args to the existing args
	if toolCall, ok := b.resp.Content[idx].(*api.ToolCallBlock); ok {
		toolCall.Args = append(toolCall.Args, slices.Clone(e.ArgsDelta)...)
	}

	return nil
}

// addSource adds a source event to the response.
func (b *ResponseBuilder) addSource(e *api.SourceEvent) error {
	b.currentState = sourceState
	b.resp.Content = append(b.resp.Content, &api.SourceBlock{
		ID:               e.Source.ID,
		URL:              e.Source.URL,
		Title:            e.Source.Title,
		ProviderMetadata: &e.Source.ProviderMetadata,
	})
	return nil
}

// addFile adds a file event to the response.
func (b *ResponseBuilder) addFile(e *api.FileEvent) error {
	b.currentState = fileState
	b.resp.Content = append(b.resp.Content, &api.FileBlock{
		MediaType: e.MediaType,
		Data:      slices.Clone(e.Data),
	})
	return nil
}

// addResponseMetadata adds response metadata to the response.
func (b *ResponseBuilder) addResponseMetadata(e *api.ResponseMetadataEvent) error {
	// First event should set the ID
	if b.responseID == "" {
		b.responseID = e.ID
	} else if b.responseID != e.ID {
		return fmt.Errorf("response ID mismatch: expected %s, got %s", b.responseID, e.ID)
	}

	if b.resp.ResponseInfo == nil {
		b.resp.ResponseInfo = &api.ResponseInfo{}
	}
	b.resp.ResponseInfo.ID = e.ID
	b.resp.ResponseInfo.Timestamp = e.Timestamp
	b.resp.ResponseInfo.ModelID = e.ModelID
	return nil
}

// addFinish adds a finish event to the response.
func (b *ResponseBuilder) addFinish(e *api.FinishEvent) error {
	if b.currentState == finishedState {
		return fmt.Errorf("invalid state: finish event received after response was already finished")
	}

	b.currentState = finishedState
	b.resp.FinishReason = e.FinishReason

	// Update usage
	if !e.Usage.IsZero() {
		b.usage = e.Usage
		b.resp.Usage = b.usage
	}

	b.resp.ProviderMetadata = e.ProviderMetadata
	return nil
}

// addError adds an error event to the response.
func (b *ResponseBuilder) addError(e *api.ErrorEvent) error {
	// Store the error to be returned by Build
	b.err = e
	return nil
}

// Build creates a Response from the collected events.
func (b *ResponseBuilder) Build() (*api.Response, error) {
	// Return any error that was encountered during event processing
	if b.err != nil {
		return &b.resp, b.err
	}
	return &b.resp, nil
}

// AddMetadata adds metadata from a StreamResponse to the builder.
// It preserves existing non-nil slices in the response.
func (b *ResponseBuilder) AddMetadata(sr *api.StreamResponse) error {
	if sr.RequestInfo != nil {
		b.resp.RequestInfo = sr.RequestInfo
	}

	if sr.ResponseInfo != nil {
		b.resp.ResponseInfo = sr.ResponseInfo
	}

	return nil
}
