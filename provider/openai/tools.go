package openai

import "go.jetify.com/ai/provider/openai/internal/codec"

type (
	FileSearchTool     = codec.FileSearchTool
	FileSearchToolCall = codec.FileSearchToolCall
	FileSearchResult   = codec.FileSearchResult
)

type (
	WebSearchTool         = codec.WebSearchTool
	WebSearchUserLocation = codec.WebSearchUserLocation
)

type (
	ComputerUseTool  = codec.ComputerUseTool
	ComputerToolCall = codec.ComputerToolCall
)
