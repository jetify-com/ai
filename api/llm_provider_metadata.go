package api

import (
	"encoding/json"
)

// ProviderMetadata provides access to provider-specific metadata structures.
// It stores and retrieves strongly-typed metadata for specific providers.
type ProviderMetadata struct {
	data map[string]any
}

// NewProviderMetadata creates a new ProviderMetadata with the given initial data.
// If data is nil, an empty map will be created.
func NewProviderMetadata(data map[string]any) *ProviderMetadata {
	if data == nil {
		return &ProviderMetadata{data: make(map[string]any)}
	}
	return &ProviderMetadata{data: data}
}

// Get retrieves the metadata for a specific provider.
// Returns the metadata and a boolean indicating whether the provider was found.
func (p *ProviderMetadata) Get(provider string) (any, bool) {
	if p == nil || p.data == nil {
		return nil, false
	}

	metadata, exists := p.data[provider]
	return metadata, exists
}

// Set stores metadata for a specific provider.
func (p *ProviderMetadata) Set(provider string, metadata any) {
	if p.data == nil {
		p.data = make(map[string]any)
	}
	p.data[provider] = metadata
}

// Has checks if metadata exists for a specific provider.
func (p *ProviderMetadata) Has(provider string) bool {
	if p == nil || p.data == nil {
		return false
	}
	_, exists := p.data[provider]
	return exists
}

// MarshalJSON implements the json.Marshaler interface.
// It serializes the underlying data map to JSON.
func (p *ProviderMetadata) MarshalJSON() ([]byte, error) {
	if p == nil || p.data == nil {
		return json.Marshal(make(map[string]any))
	}
	return json.Marshal(p.data)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It deserializes JSON into the underlying data map.
func (p *ProviderMetadata) UnmarshalJSON(data []byte) error {
	if p.data == nil {
		p.data = make(map[string]any)
	}
	return json.Unmarshal(data, &p.data)
}

// IsZero returns true if this ProviderMetadata is a zero value or contains no data.
func (p *ProviderMetadata) IsZero() bool {
	return p == nil || len(p.data) == 0
}

type MetadataSource interface {
	GetProviderMetadata() *ProviderMetadata
}

// GetMetadata is a generic helper function to retrieve provider-specific
// metadata as a pointer to the requested type.
//
// If the provider is not found or the type doesn't match, it returns nil.
//
// We recommend providers use this helper to expose predefined metadata functions.
func GetMetadata[T any](provider string, source MetadataSource) *T {
	if source == nil {
		return nil
	}

	pm := source.GetProviderMetadata()
	if pm == nil {
		return nil
	}

	metadata, ok := pm.Get(provider)
	if !ok {
		return nil
	}

	// First try to get the pointer:
	ptr, ok := metadata.(*T)
	if ok {
		return ptr
	}

	// If that fails, try the value:
	value, ok := metadata.(T)
	if ok {
		return &value
	}

	// We couldn't cast it to the right type, return nil:
	return nil
}
