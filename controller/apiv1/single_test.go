package apiv1

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
	controller := &Single{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET
	controller.Init(database, []string{"", "0"})
	database.Rows = &MockEmptyRows{}
	request := &http.Request{Method: "GET"}
	controller.Run(response, request)
	if response.StatusCode != 404 {
		t.Error("Expected 404 error when journal not found")
	}

	// Test return
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug", strings.NewReader(""))
	database.Rows = &MockJournalSingleRow{}
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected content to be returned")
	}
}
