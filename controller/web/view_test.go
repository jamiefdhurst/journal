package web

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestView_Run(t *testing.T) {
	database := &FakeSqlite{}
	response := &FakeResponse{}
	response.Reset()
	controller := &View{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET/POST
	controller.Init(database, []string{"", "0"})
	database.Rows = &MockEmptyRows{}
	request := &http.Request{Method: "GET"}
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

	// Display no error
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug", strings.NewReader(""))
	database.Rows = &MockJournalSingleRow{}
	controller.Run(response, request)
	if strings.Contains(response.Content, "div class=\"error\"") || !strings.Contains(response.Content, "Content") {
		t.Error("Expected no error to be shown in form")
	}
}
