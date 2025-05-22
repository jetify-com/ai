package api

// ImageCallWarningType represents the type of warning from the model provider.
type ImageCallWarningType string

const (
	// ImageCallWarningTypeUnsupportedSetting indicates a warning about an unsupported setting.
	ImageCallWarningTypeUnsupportedSetting ImageCallWarningType = "unsupported-setting"

	// ImageCallWarningTypeOther indicates a generic warning.
	ImageCallWarningTypeOther ImageCallWarningType = "other"
)

// ImageCallWarning represents a warning from the model provider for this call.
// The call will proceed, but e.g. some settings might not be supported,
// which can lead to suboptimal results.
type ImageCallWarning interface {
	// isImageCallWarning is a marker method to ensure type safety
	isImageCallWarning()
}

// UnsupportedSettingWarning represents a warning about an unsupported setting.
type UnsupportedSettingWarning struct {
	// Type is always ImageCallWarningTypeUnsupportedSetting for this warning
	Type ImageCallWarningType `json:"type"`

	// Setting is the name of the unsupported setting from ImageCallOptions
	Setting string `json:"setting"`

	// Details provides additional information about why the setting is unsupported
	Details *string `json:"details,omitzero"`
}

func (UnsupportedSettingWarning) isImageCallWarning() {}

// NewUnsupportedSettingWarning creates a new UnsupportedSettingWarning with the given setting and optional details
func NewUnsupportedSettingWarning(setting string, details *string) UnsupportedSettingWarning {
	return UnsupportedSettingWarning{
		Type:    ImageCallWarningTypeUnsupportedSetting,
		Setting: setting,
		Details: details,
	}
}

// OtherWarning represents a generic warning with a message.
type OtherWarning struct {
	// Type is always ImageCallWarningTypeOther for this warning
	Type ImageCallWarningType `json:"type"`

	// Message describes the warning
	Message string `json:"message"`
}

func (OtherWarning) isImageCallWarning() {}

// NewOtherWarning creates a new OtherWarning with the given message
func NewOtherWarning(message string) OtherWarning {
	return OtherWarning{
		Type:    ImageCallWarningTypeOther,
		Message: message,
	}
}
