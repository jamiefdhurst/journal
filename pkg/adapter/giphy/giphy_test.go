package giphy

import (
	"testing"

	"github.com/jamiefdhurst/journal/test/mocks/adapter"
)

func TestGiphy_SearchForID(t *testing.T) {

	var (
		response string
		err      error
	)

	// Test no API key
	client := Client{}
	_, err = client.SearchForID("test 1")
	if err == nil {
		t.Error("Expected API key error was not achieved")
	}

	// Test error
	mockClient := &adapter.MockClient{}
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
