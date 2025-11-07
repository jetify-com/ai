package openai

import (
	"go.jetify.com/ai/provider/openai/internal/codec"
)

// Tool IDs
var (
	FileSearchToolID  = codec.FileSearchToolID
	WebSearchToolID   = codec.WebSearchToolID
	ComputerUseToolID = codec.ComputerUseToolID
)

// Tool structs and related types for direct access
type (
	// Args structs for tool configuration
	FileSearchToolArgs  = codec.FileSearchToolArgs
	WebSearchToolArgs   = codec.WebSearchToolArgs
	ComputerUseToolArgs = codec.ComputerUseToolArgs

	// Option types for customization
	FileSearchToolOption  = codec.FileSearchToolOption
	WebSearchToolOption   = codec.WebSearchToolOption
	ComputerUseToolOption = codec.ComputerUseToolOption

	// Types for tool calls
	FileSearchToolCall = codec.FileSearchToolCall
	FileSearchResult   = codec.FileSearchResult

	// User location for web search
	WebSearchUserLocation = codec.WebSearchUserLocation

	ComputerToolCall    = codec.ComputerToolCall
	ComputerCoordinates = codec.ComputerCoordinates
	ComputerSafetyCheck = codec.ComputerSafetyCheck
)

// Constructor functions for creating tools
var (
	// FileSearchTool creates a new file search tool with the specified configuration.
	FileSearchTool = codec.FileSearchTool

	// WebSearchTool creates a new web search tool with the specified configuration.
	WebSearchTool = codec.WebSearchTool

	// ComputerUseTool creates a new computer use tool with the specified configuration.
	ComputerUseTool = codec.ComputerUseTool
)

// Option functions for tool customization
var (
	// File search options
	WithVectorStoreIDs = codec.WithVectorStoreIDs
	WithMaxNumResults  = codec.WithMaxNumResults

	// Web search options
	WithSearchContextSize = codec.WithSearchContextSize
	WithUserLocation      = codec.WithUserLocation

	// Computer use options
	WithDisplaySize = codec.WithDisplaySize
	WithEnvironment = codec.WithEnvironment
)
