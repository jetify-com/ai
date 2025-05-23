package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"time"

	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/responses"
	"go.jetify.com/ai/api"
)

// StreamReader is an interface for reading from an SSE stream.
// This abstraction makes testing easier as we can mock this interface instead of the concrete ssestream.Stream type.
type StreamReader interface {
	Next() bool
	Current() responses.ResponseStreamEventUnion
	Err() error
}

// DecodeStream converts an OpenAI SSE stream to our API's StreamResponse.
// This is the main entry point for decoding OpenAI streams.
func DecodeStream(stream StreamReader) (api.StreamResponse, error) {
	decoder := &streamDecoder{}
	return decoder.DecodeStream(stream)
}

// streamDecoder maintains state while decoding a stream of OpenAI events.
type streamDecoder struct {
	// Map from output index to tool call information
	ongoingToolCalls map[int64]toolCallInfo

	// Tracking response metadata
	responseID string

	// Usage statistics
	promptTokens     int
	completionTokens int
	totalTokens      int

	// Advanced usage statistics
	inputCachedTokens     int
	outputReasoningTokens int

	// Flags
	hasToolCalls bool

	// Incomplete details for determining finish reason
	incompleteReason string

	// Counter for source annotations
	annotationCounter int
}

// toolCallInfo tracks information about an ongoing tool call.
type toolCallInfo struct {
	toolName   string
	toolCallID string
}

// DecodeStream processes an OpenAI stream and returns our API stream format.
func (d *streamDecoder) DecodeStream(stream StreamReader) (api.StreamResponse, error) {
	if d.ongoingToolCalls == nil {
		d.ongoingToolCalls = make(map[int64]toolCallInfo)
	}

	return api.StreamResponse{
		Stream: d.decodeEvents(stream),
	}, nil
}

// decodeEvents returns an iterator that yields events from the OpenAI stream.
func (d *streamDecoder) decodeEvents(stream StreamReader) iter.Seq[api.StreamEvent] {
	return func(yield func(api.StreamEvent) bool) {
		// Process all events directly in the iterator function
		for stream.Next() {
			// Get the current event
			event := stream.Current()

			// Process the event
			decodedEvent := d.decodeEvent(event)

			// Only send non-nil events (this handles events that are processed internally but don't yield output,
			// and also events that are explicitly decoded to an api.ErrorEvent)
			if decodedEvent != nil {
				if !yield(decodedEvent) {
					return
				}
			}
		}

		// Check if we encountered an error from the underlying stream
		if err := stream.Err(); err != nil && !errors.Is(err, io.EOF) {
			// Yield this as a final error event
			if !yield(&api.ErrorEvent{Err: err}) {
				// Consumer doesn't want more, even this error
				return
			}
		}

		// Create provider metadata with full usage information
		metadata := api.NewProviderMetadata(map[string]any{
			"openai": &Metadata{
				ResponseID: d.responseID,
				Usage: Usage{
					InputTokens:           d.promptTokens,
					OutputTokens:          d.completionTokens,
					InputCachedTokens:     d.inputCachedTokens,
					OutputReasoningTokens: d.outputReasoningTokens,
				},
			},
		})

		// Determine finish reason based on decoder state
		finishReason := decodeFinishReason(d.incompleteReason, d.hasToolCalls)

		// Send the final finish event
		yield(&api.FinishEvent{
			FinishReason: finishReason,
			Usage: api.Usage{
				InputTokens:       d.promptTokens,
				OutputTokens:      d.completionTokens,
				TotalTokens:       d.totalTokens,
				ReasoningTokens:   d.outputReasoningTokens,
				CachedInputTokens: d.inputCachedTokens,
			},
			ProviderMetadata: metadata,
		})
	}
}

// decodeEvent translates an OpenAI event to our API event format.
// It now returns a single api.StreamEvent, which can be an api.ErrorEvent
// if an internal processing error occurs or if an OpenAI error event is decoded.
// It returns nil if the event is known but intentionally not exposed to clients.
func (d *streamDecoder) decodeEvent(event responses.ResponseStreamEventUnion) api.StreamEvent {
	switch event.Type {
	case "response.output_text.delta":
		return d.decodeTextDelta(event)
	case "response.output_item.added":
		return d.decodeOutputItemAdded(event)
	case "response.function_call_arguments.delta":
		return d.decodeFunctionCallArgumentsDelta(event)
	case "response.output_item.done":
		return d.decodeOutputItemDone(event)
	case "response.created":
		return d.decodeResponseCreated(event)
	case "response.completed":
		return d.decodeResponseCompleted(event)
	case "response.failed", "response.incomplete":
		return d.decodeResponseFailedOrIncomplete(event)
	case "response.reasoning_summary_text.delta":
		return d.decodeReasoningSummaryTextDelta(event)
	case "response.output_text.annotation.added":
		return d.decodeOutputTextAnnotationAdded(event)
	case "response.text_annotation.delta":
		return d.decodeTextAnnotationDelta(event)
	case "error":
		return d.decodeError(event)
	// Event types that we're aware of but don't yet expose to clients:
	case "response.in_progress",
		"response.content_part.done",
		"response.content_part.added",
		"response.output_text.done",
		"response.refusal.delta",
		"response.refusal.done",
		"response.function_call_arguments.done",
		"response.file_search_call.in_progress",
		"response.file_search_call.searching",
		"response.file_search_call.completed",
		"response.web_search_call.in_progress",
		"response.web_search_call.searching",
		"response.web_search_call.completed",
		"response.reasoning_summary_part.added",
		"response.reasoning_summary_part.done",
		"response.reasoning_summary_text.done",
		"response.audio.delta",
		"response.audio.done",
		"response.audio.transcript.delta",
		"response.audio.transcript.done",
		"response.code_interpreter_call.code.delta",
		"response.code_interpreter_call.code.done",
		"response.code_interpreter_call.completed",
		"response.code_interpreter_call.in_progress",
		"response.code_interpreter_call.interpreting":
		// We're aware these events exist but we don't expose them to clients yet.
		return nil
	default:
		// Ignore any other unsupported event types
		// TODO: Return a warning as an api.ErrorEvent?
		// For now, returning nil to maintain current behavior of ignoring.
		return nil
	}
}

// decodeTextDelta handles text delta events
func (d *streamDecoder) decodeTextDelta(event responses.ResponseStreamEventUnion) api.StreamEvent {
	textDelta := event.AsResponseOutputTextDelta()
	return &api.TextDeltaEvent{
		TextDelta: textDelta.Delta,
	}
}

// decodeOutputItemAdded handles output item added events
func (d *streamDecoder) decodeOutputItemAdded(event responses.ResponseStreamEventUnion) api.StreamEvent {
	itemAdded := event.AsResponseOutputItemAdded()
	item := itemAdded.Item

	if item.Type == "function_call" {
		funcCall := item.AsFunctionCall()

		// Store the tool call information for later deltas
		d.ongoingToolCalls[itemAdded.OutputIndex] = toolCallInfo{
			toolName:   funcCall.Name,
			toolCallID: funcCall.CallID,
		}
		d.hasToolCalls = true

		return &api.ToolCallDeltaEvent{
			ToolCallID: funcCall.CallID,
			ToolName:   funcCall.Name,
			ArgsDelta:  []byte(funcCall.Arguments),
		}
	}

	if item.Type == "reasoning" {
		return nil
	}

	return nil
}

// decodeFunctionCallArgumentsDelta handles function call arguments delta events
func (d *streamDecoder) decodeFunctionCallArgumentsDelta(event responses.ResponseStreamEventUnion) api.StreamEvent {
	argsDelta := event.AsResponseFunctionCallArgumentsDelta()
	toolCall, exists := d.ongoingToolCalls[argsDelta.OutputIndex]

	if !exists {
		return &api.ErrorEvent{Err: fmt.Errorf("received function call arguments delta for unknown output index: %d", argsDelta.OutputIndex)}
	}

	return &api.ToolCallDeltaEvent{
		ToolCallID: toolCall.toolCallID,
		ToolName:   toolCall.toolName,
		ArgsDelta:  []byte(argsDelta.Delta),
	}
}

// decodeOutputItemDone handles output item done events
func (d *streamDecoder) decodeOutputItemDone(event responses.ResponseStreamEventUnion) api.StreamEvent {
	itemDone := event.AsResponseOutputItemDone()
	item := itemDone.Item

	if item.Type == "function_call" {
		funcCall := item.AsFunctionCall()
		delete(d.ongoingToolCalls, itemDone.OutputIndex)
		return &api.ToolCallEvent{
			ToolCallID: funcCall.CallID,
			ToolName:   funcCall.Name,
			Args:       json.RawMessage(funcCall.Arguments),
		}
	}
	return nil
}

// decodeResponseCreated handles response created events
func (d *streamDecoder) decodeResponseCreated(event responses.ResponseStreamEventUnion) api.StreamEvent {
	created := event.AsResponseCreated()
	d.responseID = created.Response.ID
	timestamp := time.Unix(int64(created.Response.CreatedAt), 0).UTC()
	return &api.ResponseMetadataEvent{
		ID:        created.Response.ID,
		Timestamp: timestamp,
		ModelID:   created.Response.Model,
	}
}

// decodeResponseCompleted processes the response.completed event which has final usage statistics
func (d *streamDecoder) decodeResponseCompleted(event responses.ResponseStreamEventUnion) api.StreamEvent {
	completed := event.AsResponseCompleted()

	// Update usage statistics from the response if available
	usage := completed.Response.Usage
	if !param.IsOmitted(usage) {
		d.promptTokens = int(usage.InputTokens)
		d.completionTokens = int(usage.OutputTokens)

		// If totalTokens is provided and non-zero, use it; otherwise compute it
		totalTokens := int(usage.TotalTokens)
		if totalTokens == 0 {
			totalTokens = d.promptTokens + d.completionTokens
		}
		d.totalTokens = totalTokens

		// Also update advanced usage statistics if available
		if !param.IsOmitted(usage.InputTokensDetails) {
			d.inputCachedTokens = int(usage.InputTokensDetails.CachedTokens)
		}
		if !param.IsOmitted(usage.OutputTokensDetails) {
			d.outputReasoningTokens = int(usage.OutputTokensDetails.ReasoningTokens)
		}
	}

	// Preserve the existing behavior for incomplete reason
	if completed.Response.IncompleteDetails.Reason != "" {
		d.incompleteReason = completed.Response.IncompleteDetails.Reason
	}

	return nil
}

// decodeResponseFailedOrIncomplete handles response failed or incomplete events
// This function updates internal state and does not yield an event itself.
// The final finish event's reason is affected, and errors might be reported via FinishEvent.
func (d *streamDecoder) decodeResponseFailedOrIncomplete(event responses.ResponseStreamEventUnion) api.StreamEvent {
	var reason string
	if event.Type == "response.failed" {
		failedEvent := event.AsResponseFailed()
		reason = failedEvent.Response.IncompleteDetails.Reason
		// Potentially, if failedEvent.Response.Error is not nil, we could also emit an ErrorEvent here.
		// For now, sticking to updating incompleteReason for finish_reason consistency.
		// if errDetails := failedEvent.Response.Error; errDetails.Code != "" || errDetails.Message != "" {
		//  // This would be a place to consider if a direct error event is also needed
		// }
	} else if event.Type == "response.incomplete" {
		incompleteEvent := event.AsResponseIncomplete()
		reason = incompleteEvent.Response.IncompleteDetails.Reason
	}
	d.incompleteReason = reason
	return nil // No event yielded directly from here
}

// decodeReasoningSummaryTextDelta handles reasoning summary text delta events
func (d *streamDecoder) decodeReasoningSummaryTextDelta(event responses.ResponseStreamEventUnion) api.StreamEvent {
	return &api.ReasoningEvent{
		TextDelta: event.Delta,
	}
}

// decodeOutputTextAnnotationAdded handles response.output_text.annotation.added events
func (d *streamDecoder) decodeOutputTextAnnotationAdded(event responses.ResponseStreamEventUnion) api.StreamEvent {
	addedEvent := event.AsResponseOutputTextAnnotationAdded()
	if addedEvent.Annotation.Type == "url_citation" {
		citation := addedEvent.Annotation.AsURLCitation()
		sourceEvent := &api.SourceEvent{
			Source: api.Source{
				SourceType: "url",
				ID:         fmt.Sprintf("source-%d", d.annotationCounter),
				URL:        citation.URL,
				Title:      citation.Title,
			},
		}
		d.annotationCounter++
		return sourceEvent
	}
	return nil
}

// decodeTextAnnotationDelta handles response.text_annotation.delta events
func (d *streamDecoder) decodeTextAnnotationDelta(event responses.ResponseStreamEventUnion) api.StreamEvent {
	if event.Annotation.Type == "url_citation" {
		citation := event.Annotation.AsURLCitation()
		sourceEvent := &api.SourceEvent{
			Source: api.Source{
				SourceType: "url",
				ID:         fmt.Sprintf("source-%d", d.annotationCounter),
				URL:        citation.URL,
				Title:      citation.Title,
			},
		}
		d.annotationCounter++
		return sourceEvent
	}
	return nil
}

// decodeError handles error events from the OpenAI stream
func (d *streamDecoder) decodeError(event responses.ResponseStreamEventUnion) api.StreamEvent {
	errorEvent := event.AsError()
	return &api.ErrorEvent{
		Err: fmt.Errorf("%s: %s", errorEvent.Code, errorEvent.Message),
	}
}
