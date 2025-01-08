package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
)

// MockClient Mock an HTTP client
type MockClient struct {
	ErrorMode bool
	Response  string
}

// Get Trigger a fake GET, with JSON
func (m *MockClient) Get(url string, destination interface{}) error {
	if m.ErrorMode {
		return errors.New("Simulated error")
	}

	if m.Response != "" {
		json.NewDecoder(bytes.NewBufferString(m.Response)).Decode(&destination)
	}

	return nil
}
