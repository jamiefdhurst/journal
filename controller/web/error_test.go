package web

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestError_Run(t *testing.T) {
	response := &FakeResponse{}
	controller := &Error{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test header and response
	controller.Run(response, &http.Request{})
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

}
