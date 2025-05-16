package anthropic

// ProviderName is the name of the Anthropic provider.
const ProviderName = "anthropic"

const (
	// Claude 3.7 Models

	ModelClaude3_7SonnetLatest  = "claude-3-7-sonnet-latest"
	ModelClaude37Sonnet20250219 = "claude-3-7-sonnet-20250219"

	// Claude 3.5 Models

	ModelClaude35HaikuLatest   = "claude-3-5-haiku-latest"
	ModelClaude35Haiku20241022 = "claude-3-5-haiku-20241022"

	ModelClaude35SonnetLatest   = "claude-3-5-sonnet-latest"
	ModelClaude35Sonnet20241022 = "claude-3-5-sonnet-20241022"
	ModelClaude35Sonnet20240620 = "claude-3-5-sonnet-20240620"

	// Claude 3.0 Models

	ModelClaude3OpusLatest   = "claude-3-opus-latest"
	ModelClaude3Opus20240229 = "claude-3-opus-20240229"
	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude3Sonnet20240229 = "claude-3-sonnet-20240229"
	ModelClaude3Haiku20240307  = "claude-3-haiku-20240307"

	// Claude 2 Models

	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude21 = "claude-2.1"
	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude20 = "claude-2.0"
)
