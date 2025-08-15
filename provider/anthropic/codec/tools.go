package codec

import (
	"encoding/json"

	"go.jetify.com/ai/api"
)

// Default versions and recommended settings for Anthropic tools
const (
	DefaultToolVersion       = "20250124"
	RecommendedDisplayWidth  = 1280
	RecommendedDisplayHeight = 800
)

// ComputerAction is a ComputerToolCall action.
type ComputerAction string

const (
	// ActionKey presses the key or key-combination specfied in
	// [ComputerToolCall.Text]. It supports xdotool's key syntax.
	ActionKey ComputerAction = "key"

	// ActionHoldKey holds down the key or key-combination in
	// [ComputerToolCall.Text] for a duration of
	// [ComputerToolCall.Duration]. It supports the same syntax as
	// [ActionKey].
	ActionHoldKey ComputerAction = "hold_key"

	// ActionType types the string specified by [ComputerToolCall.Text] on
	// the keyboard.
	ActionType ComputerAction = "type"

	// ActionCursorPosition reports the current (x, y) pixel coordinates of
	// the cursor on the screen.
	ActionCursorPosition ComputerAction = "cursor_position"

	// ActionMouseMove moves the cursor to the pixel specified by
	// [ComputerToolCall.Coordinate].
	ActionMouseMove ComputerAction = "mouse_move"

	// ActionLeftMouseDown presses down the left mouse button without
	// releasing it.
	ActionLeftMouseDown ComputerAction = "left_mouse_down"

	// ActionLeftMouseUp releases the left mouse button.
	ActionLeftMouseUp ComputerAction = "left_mouse_up"

	// ActionLeftClick clicks the left mouse button at
	// [ComputerToolCall.Coordinate] while optionally holding down the keys
	// in [ComputerToolCall.Text].
	ActionLeftClick ComputerAction = "left_click"

	// ActionLeftClickDrag clicks and drags the cursor from
	// [ComputerToolCall.StartCoordinate] to [ComputerToolCall.Coordinate].
	ActionLeftClickDrag ComputerAction = "left_click_drag"

	// ActionRightClick clicks the right mouse button at
	// [ComputerToolCall.Coordinate].
	ActionRightClick ComputerAction = "right_click"

	// ActionMiddleClick clicks the middle mouse button at
	// [ComputerToolCall.Coordinate].
	ActionMiddleClick ComputerAction = "middle_click"

	// ActionDoubleClick double-clicks the left mouse button at
	// [ComputerToolCall.Coordinate].
	ActionDoubleClick ComputerAction = "double_click"

	// ActionTripleClick triple-clicks the left mouse button at
	// [ComputerToolCall.Coordinate].
	ActionTripleClick ComputerAction = "triple_click"

	// ActionScroll turns the mouse scroll wheel by
	// [ComputerToolCall.ScrollAmount] in the direction of
	// [ComputerToolCall.ScrollDirection] at [ComputerToolCall.Coordinate].
	ActionScroll ComputerAction = "scroll"

	// ActionWait pauses execution for [ComputerToolCall.Duration].
	ActionWait ComputerAction = "wait"

	// ActionScreenshot takes a screenshot.
	ActionScreenshot ComputerAction = "screenshot"
)

// ScrollDirection is a direction to scroll the screen.
type ScrollDirection string

const (
	ScrollUp    = "up"
	ScrollDown  = "down"
	ScrollLeft  = "left"
	ScrollRight = "right"
)

// TODO(gcurtis): make ComputerToolCall implement json.Unmarshaler so it can
// have better types for some of its fields:
//
// 	- dedicated coordinate type
// 	- duration should be time.Duration
// 	- use ints instead of json.Number while still being flexible about
// 	  accepting ints, floats, or number strings

// ComputerToolCall contains the parameters of a call to [ComputerTool].
type ComputerToolCall struct {
	// Action is the action to perform. It is the only mandatory field.
	Action ComputerAction `json:"action"`

	// Text is a key, key-combination, or string literal to type on the
	// keyboard. It specifies individual key presses or key-combinations
	// using an xdotool-style syntax. Examples include "a", "Return",
	// "alt+Tab", "ctrl+s", "Up", "KP_0" (for numpad 0). The ActionType,
	// ActionKey, and ActionHoldKey actions require a non-empty Text value.
	// Click or scroll actions may optionally set Text to specify keys to
	// hold down keys during the click or scroll.
	Text string `json:"text,omitzero"`

	// Coordinate is a pair of (x, y) on-screen pixel coordinates for cursor
	// actions. (0, 0) is the top-left corner of the screen. The
	// ActionMouseMove and ActionLeftClickDrag actions require a coordinate.
	Coordinate [2]json.Number `json:"coordinate,omitzero"`

	// StartCoordinate is the starting point for mouse drag actions.
	StartCoordinate [2]json.Number `json:"start_coordinate,omitzero"`

	// Duration is the number of seconds to hold down keys or pause
	// execution. The ActionHoldKey and ActionWait actions require a
	// non-zero Duration.
	Duration json.Number `json:"duration,omitzero"`

	// ScrollAmount is the number of mouse wheel "clicks" to scroll.
	// ActionScroll requires a non-zero ScrollAmount.
	ScrollAmount json.Number `json:"scroll_amount,omitzero"`

	// ScrollDirection is the direction to scroll. ActionScroll requires
	// a ScrollDirection.
	ScrollDirection ScrollDirection `json:"scroll_direction,omitzero"`
}

// ComputerToolArgs contains the configuration arguments for the computer use tool.
// See the [computer use guide](https://docs.anthropic.com/en/docs/agents-and-tools/computer-use) for more details.
type ComputerToolArgs struct {
	// The version of the computer tool to use.
	// Optional field, defaults to the latest version. Possible values are: "20250124", "20241022".
	Version string `json:"version"`

	// The width of the display being controlled by the model in pixels.
	// Required field. We recommend setting it to 1280.
	DisplayWidthPx int `json:"display_width_px"`
	// The height of the display being controlled by the model in pixels.
	// Required field. We recommend setting it to 800.
	DisplayHeightPx int `json:"display_height_px"`

	// The display number to control (only relevant for X11 environments).
	// Optional field, if specified, the tool will be provided a display number in the tool definition.
	DisplayNumber int `json:"display_number,omitzero"`
}

// ComputerToolOption allows customizing computer tool configuration.
type ComputerToolOption func(*ComputerToolArgs)

// WithComputerVersion sets the computer tool version.
func WithComputerVersion(version string) ComputerToolOption {
	return func(args *ComputerToolArgs) {
		args.Version = version
	}
}

// WithDisplayNumber sets the display number for X11 environments.
func WithDisplayNumber(displayNum int) ComputerToolOption {
	return func(args *ComputerToolArgs) {
		args.DisplayNumber = displayNum
	}
}

// ComputerTool creates a new computer use tool with the specified configuration.
// ComputerTool is a built-in tool that can be used to control a computer.
// It allows the model to use a mouse and keyboard and to take screenshots.
// See the [computer use guide](https://docs.anthropic.com/en/docs/agents-and-tools/computer-use) for more details.
func ComputerTool(displayWidth, displayHeight int, options ...ComputerToolOption) *api.ProviderDefinedTool {
	args := &ComputerToolArgs{
		DisplayWidthPx:  displayWidth,
		DisplayHeightPx: displayHeight,
		Version:         DefaultToolVersion,
	}

	// Apply options
	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   "anthropic.computer",
		Name: "computer",
		Args: args,
	}
}

// BashToolArgs contains the configuration arguments for the bash tool.
// See the [computer use guide](https://docs.anthropic.com/en/docs/agents-and-tools/computer-use) for more details.
type BashToolArgs struct {
	// The version of the bash tool to use.
	// Optional field, defaults to the latest version. Possible values are: "20250124", "20241022".
	Version string `json:"version"`
}

// BashToolOption allows customizing bash tool configuration.
type BashToolOption func(*BashToolArgs)

// WithBashVersion sets the bash tool version.
func WithBashVersion(version string) BashToolOption {
	return func(args *BashToolArgs) {
		args.Version = version
	}
}

// BashTool creates a new bash tool with the specified configuration.
// BashTool is a built-in tool that can be used to run shell commands.
// See the [computer use guide](https://docs.anthropic.com/en/docs/agents-and-tools/computer-use) for more details.
func BashTool(options ...BashToolOption) *api.ProviderDefinedTool {
	args := &BashToolArgs{
		Version: DefaultToolVersion,
	}

	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   "anthropic.bash",
		Name: "bash",
		Args: args,
	}
}

// TextEditorToolArgs contains the configuration arguments for the text editor tool.
// See the [text editor guide](https://docs.anthropic.com/en/docs/build-with-claude/tool-use/text-editor-tool) for more details.
type TextEditorToolArgs struct {
	// The version of the text editor tool to use.
	// Optional field, defaults to the latest version. Possible values are: "20250124", "20241022".
	Version string `json:"version"`
}

// TextEditorToolOption allows customizing text editor tool configuration.
type TextEditorToolOption func(*TextEditorToolArgs)

// WithTextEditorVersion sets the text editor tool version.
func WithTextEditorVersion(version string) TextEditorToolOption {
	return func(args *TextEditorToolArgs) {
		args.Version = version
	}
}

// TextEditorTool creates a new text editor tool with the specified configuration.
// TextEditorTool is a built-in tool that can be used to view, create and edit text files.
// See the [text editor guide](https://docs.anthropic.com/en/docs/build-with-claude/tool-use/text-editor-tool) for more details.
func TextEditorTool(options ...TextEditorToolOption) *api.ProviderDefinedTool {
	args := &TextEditorToolArgs{
		Version: DefaultToolVersion,
	}

	for _, opt := range options {
		opt(args)
	}

	return &api.ProviderDefinedTool{
		ID:   "anthropic.text_editor",
		Name: "str_replace_editor",
		Args: args,
	}
}

// TODO: Add predefined tool call blocks for the different built-in tools.
