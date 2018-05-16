package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
)

// MockGiphyAdapter Mock the Giphy adapter
type MockGiphyAdapter struct {
	ErrorMode bool
}

// SearchForID Present to search
func (m MockGiphyAdapter) SearchForID(s string) (string, error) {
	if m.ErrorMode {
		return "", errors.New("Simulating error")
	}
	if s == "testsearch" {
		return "9991234", nil
	}

	return "0000000", nil
}

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
