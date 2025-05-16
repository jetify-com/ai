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
	// Map of tool call ID to index in orderedToolCalls
	toolCallIndices map[string]int
	// Ordered slice of tool calls to maintain order
	orderedToolCalls []api.ToolCallBlock
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
			// Initialize with empty values
			Text:         "",
			Reasoning:    []api.Reasoning{},
			ToolCalls:    []api.ToolCallBlock{},
			Sources:      []api.Source{},
			Warnings:     []api.CallWarning{},
			FinishReason: api.FinishReason(""),
			Files:        []api.FileBlock{},
		},
		toolCallIndices:  make(map[string]int),
		orderedToolCalls: make([]api.ToolCallBlock, 0),
		currentState:     noState,
		err:              nil,
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

	switch e := event.(type) {
	case *api.TextDeltaEvent:
		return b.addTextDelta(e)
	case *api.ReasoningEvent:
		return b.addReasoning(e)
	case *api.ReasoningSignatureEvent:
		return b.addReasoningSignature(e)
	case *api.RedactedReasoningEvent:
		return b.addRedactedReasoning(e)
	case *api.ToolCallEvent:
		return b.addToolCall(e)
	case *api.ToolCallDeltaEvent:
		return b.addToolCallDelta(e)
	case *api.SourceEvent:
		return b.addSource(e)
	case *api.FileEvent:
		return b.addFile(e)
	case *api.ResponseMetadataEvent:
		return b.addResponseMetadata(e)
	case *api.FinishEvent:
		return b.addFinish(e)
	case *api.ErrorEvent:
		return b.addError(e)
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

// addTextDelta adds a text delta event to the response.
func (b *ResponseBuilder) addTextDelta(e *api.TextDeltaEvent) error {
	b.currentState = textState
	b.resp.Text += e.TextDelta
	return nil
}

// addReasoning adds a reasoning event to the response.
func (b *ResponseBuilder) addReasoning(e *api.ReasoningEvent) error {
	// Validate state transition
	if b.currentState != noState && b.currentState != reasoningState {
		return fmt.Errorf("invalid state transition: cannot add reasoning in state %v", b.currentState)
	}

	b.currentState = reasoningState
	b.resp.Reasoning = append(b.resp.Reasoning, &api.ReasoningBlock{
		Text: e.TextDelta,
	})
	return nil
}

// addReasoningSignature adds a reasoning signature event to the response.
func (b *ResponseBuilder) addReasoningSignature(e *api.ReasoningSignatureEvent) error {
	if len(b.resp.Reasoning) == 0 {
		return fmt.Errorf("cannot add reasoning signature: no reasoning blocks exist")
	}

	if block, ok := b.resp.Reasoning[len(b.resp.Reasoning)-1].(*api.ReasoningBlock); ok {
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
	b.resp.Reasoning = append(b.resp.Reasoning, &api.RedactedReasoningBlock{
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
	b.orderedToolCalls = append(b.orderedToolCalls, api.ToolCallBlock{
		ToolCallID: e.ToolCallID,
		ToolName:   e.ToolName,
		Args:       slices.Clone(e.Args),
	})
	// Store index in the map
	b.toolCallIndices[e.ToolCallID] = len(b.orderedToolCalls) - 1

	return nil
}

// addToolCallDelta adds a tool call delta event to the response.
func (b *ResponseBuilder) addToolCallDelta(e *api.ToolCallDeltaEvent) error {
	b.currentState = toolCallState

	// Get or create the tool call block
	idx, exists := b.toolCallIndices[e.ToolCallID]
	if !exists {
		// Create new tool call block if it doesn't exist
		b.orderedToolCalls = append(b.orderedToolCalls, api.ToolCallBlock{
			ToolCallID: e.ToolCallID,
			ToolName:   e.ToolName,
			Args:       make([]byte, 0),
		})
		idx = len(b.orderedToolCalls) - 1
		b.toolCallIndices[e.ToolCallID] = idx
	}

	// Append the new args to the existing args
	b.orderedToolCalls[idx].Args = append(b.orderedToolCalls[idx].Args, slices.Clone(e.ArgsDelta)...)

	return nil
}

// addSource adds a source event to the response.
func (b *ResponseBuilder) addSource(e *api.SourceEvent) error {
	b.currentState = sourceState
	// Create a copy of the source to prevent mutation
	sourceCopy := e.Source
	b.resp.Sources = append(b.resp.Sources, sourceCopy)
	return nil
}

// addFile adds a file event to the response.
func (b *ResponseBuilder) addFile(e *api.FileEvent) error {
	b.currentState = fileState
	b.resp.Files = append(b.resp.Files, api.FileBlock{
		MimeType: e.MimeType,
		Data:     slices.Clone(e.Data),
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
	if e.Usage != nil {
		b.usage = *e.Usage
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
func (b *ResponseBuilder) Build() (api.Response, error) {
	// Copy the ordered tool calls to the response
	b.resp.ToolCalls = slices.Clone(b.orderedToolCalls)

	// Return any error that was encountered during event processing
	if b.err != nil {
		return b.resp, b.err
	}
	return b.resp, nil
}

// AddMetadata adds metadata from a StreamResponse to the builder.
// It preserves existing non-nil slices in the response.
func (b *ResponseBuilder) AddMetadata(sr *api.StreamResponse) error {
	// Only copy warnings if the source has warnings
	if len(sr.Warnings) > 0 {
		b.resp.Warnings = sr.Warnings
	}

	// Copy raw call info
	b.resp.RawCall = sr.RawCall

	// Copy other metadata if not nil
	if sr.RawResponse != nil {
		b.resp.RawResponse = sr.RawResponse
	}

	if sr.RequestInfo != nil {
		b.resp.RequestInfo = sr.RequestInfo
	}

	if sr.ProviderMetadata != nil {
		b.resp.ProviderMetadata = sr.ProviderMetadata
	}

	return nil
}
