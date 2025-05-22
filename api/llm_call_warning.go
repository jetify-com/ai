package api

// CallWarning represents a warning from the model provider for a call.
// The call will proceed, but some settings might not be supported,
// which can lead to suboptimal results.
type CallWarning struct {
	// TODO: We might want to turn Type into an enum
	// OR we could make Warning an interface with different concrete types.

	// Type indicates the kind of warning: "unsupported-setting", "unsupported-tool", or "other"
	Type string `json:"type"`

	// TODO: These are usually called Configs or Options in go ... should we update
	// the name of the field?

	// Setting contains the name of the unsupported setting when Type is "unsupported-setting"
	Setting string `json:"setting,omitzero"`

	// Tool contains the unsupported tool when Type is "unsupported-tool"
	Tool ToolDefinition `json:"tool,omitzero"`

	// Details provides additional information about the warning
	Details string `json:"details,omitzero"`

	// Message contains a human-readable warning message
	Message string `json:"message,omitzero"`
}
