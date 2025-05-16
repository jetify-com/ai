package anthropic

import "go.jetify.com/ai/provider/anthropic/codec"

type (
	ComputerUseTool  = codec.ComputerUseTool
	ComputerToolCall = codec.ComputerToolCall
)

type BashTool = codec.BashTool

type TextEditorTool = codec.TextEditorTool
