package web

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
)

func TestError_Run(t *testing.T) {
	response := controller.NewMockResponse()
	configuration := app.DefaultConfiguration()
	container := &app.Container{Configuration: configuration}
	controller := &BadRequest{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test header and response
	controller.Init(container, []string{}, request)
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}
	if !strings.Contains(response.Content, "<title>Page Not Found - Jamie's Journal</title>") {
		t.Error("Expected HTML title to be in place")
	}

}
