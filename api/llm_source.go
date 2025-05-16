package api

// Source represents a source that has been used as input to generate the response.
type Source struct {
	// SourceType indicates the type of source. Currently only "url" is supported.
	SourceType string `json:"source_type"`

	// ID is the unique identifier of the source.
	ID string `json:"id"`

	// URL is the URL of the source.
	URL string `json:"url"`

	// Title is the optional title of the source.
	Title string `json:"title,omitempty"`

	// ProviderMetadata contains additional provider-specific metadata.
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}
