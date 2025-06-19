package anthropic

import "go.jetify.com/ai/provider/anthropic/codec"

// Factory functions for creating Anthropic tools
var (
	// ComputerTool creates a new computer use tool with the specified configuration.
	ComputerTool = codec.ComputerTool

	// BashTool creates a new bash tool with the specified configuration.
	BashTool = codec.BashTool

	// TextEditorTool creates a new text editor tool with the specified configuration.
	TextEditorTool = codec.TextEditorTool
)

// Option functions for customizing tools
var (
	// WithComputerVersion sets the computer tool version.
	WithComputerVersion = codec.WithComputerVersion

	// WithDisplayNumber sets the display number for X11 environments.
	WithDisplayNumber = codec.WithDisplayNumber

	// WithBashVersion sets the bash tool version.
	WithBashVersion = codec.WithBashVersion

	// WithTextEditorVersion sets the text editor tool version.
	WithTextEditorVersion = codec.WithTextEditorVersion
)

// Tool call types for handling responses
type ComputerToolCall = codec.ComputerToolCall

// Recommended constants
const (
	DefaultToolVersion       = codec.DefaultToolVersion
	RecommendedDisplayWidth  = codec.RecommendedDisplayWidth
	RecommendedDisplayHeight = codec.RecommendedDisplayHeight
)
