package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

type MockClient struct {
	ErrorMode bool
	Response  string
}

func (m *MockClient) Get(url string, destination interface{}) error {
	if m.ErrorMode {
		return errors.New("Simulated error")
	}

	if m.Response != "" {
		json.NewDecoder(bytes.NewBufferString(m.Response)).Decode(&destination)
	}

	return nil
}

func TestGiphy_SearchForID(t *testing.T) {

	var (
		response string
		err      error
	)

	// Test no API key
	client := GiphyClient{}
	_, err = client.SearchForID("test 1")
	if err == nil {
		t.Error("Expected API key error was not achieved")
	}

	// Test error
	mockClient := &MockClient{}
	mockClient.ErrorMode = true
	client.APIKey = "API123456"
	client.Client = mockClient
	_, err = client.SearchForID("test 2")
	if err == nil {
		t.Error("Expected clietn error was not achieved")
	}

	// Test valid response
	validJSON := "{\"data\":{\"id\":\"testing123\"}}"
	mockClient.ErrorMode = false
	mockClient.Response = validJSON
	response, err = client.SearchForID("test 3")
	if err != nil || response != "testing123" {
		t.Error("Expected ID to be returned")
	}

	// Test invalid response
	invalidJSON := "{\"data\":{}}"
	mockClient.Response = invalidJSON
	response, err = client.SearchForID("test 4")
	if err == nil || response != "" {
		t.Error("Expected error to be returned with invalid response")
	}
}
