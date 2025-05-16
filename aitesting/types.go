package aitesting

import "go.jetify.com/ai/api"

// MockMetadataSource implements api.MetadataSource for testing
type MockMetadataSource struct {
	ProviderMetadata *api.ProviderMetadata
}

func (m *MockMetadataSource) GetProviderMetadata() *api.ProviderMetadata {
	return m.ProviderMetadata
}

// MockUnsupportedMessage implements api.Message but is not a known type
type MockUnsupportedMessage struct{}

func (m *MockUnsupportedMessage) Role() api.MessageRole                      { return "unsupported" }
func (m *MockUnsupportedMessage) GetProviderMetadata() *api.ProviderMetadata { return nil }

// MockUnsupportedBlock implements api.ContentBlock for testing unsupported content types
type MockUnsupportedBlock struct{}

func (m *MockUnsupportedBlock) Type() api.ContentBlockType                 { return "unsupported" }
func (m *MockUnsupportedBlock) GetProviderMetadata() *api.ProviderMetadata { return nil }
