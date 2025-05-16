package builder

import "go.jetify.com/ai/api"

func StreamToResponse(stream *api.StreamResponse) (*api.Response, error) {
	if stream == nil {
		return nil, nil
	}

	builder := NewResponseBuilder()

	// Add any metadata from the stream response
	if err := builder.AddMetadata(stream); err != nil {
		return nil, err
	}

	// Process each event in the stream
	for event := range stream.Events {
		if err := builder.AddEvent(event); err != nil {
			return nil, err
		}
	}

	// Build the final response
	resp, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
