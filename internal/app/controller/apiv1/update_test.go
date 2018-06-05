package apiv1

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestUpdate_Run(t *testing.T) {
	db := &database.MockSqlite{}
	response := &controller.MockResponse{}
	response.Reset()
	controller := &Update{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET/POST
	controller.Init(db, []string{"", "0"})
	db.Rows = &database.MockRowsEmpty{}
	request := &http.Request{Method: "POST"}
	controller.Run(response, request)
	if response.StatusCode != 404 {
		t.Error("Expected 404 error when journal not found")
	}

	// Test for bad request on invalid JSON
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("{\"not\":\"valid\":\"json\"}"))
	request.Header.Add("Content-Type", "application/json")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when invalid JSON provided")
	}

	// Test Journal is retrieved on save
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("{\"title\":\"Something New\",\"date\":\"2018-01-01\",\"content\":\"New\"}"))
	request.Header.Add("Content-Type", "application/json")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 200 || !strings.Contains(response.Content, "Something New") {
		t.Error("Expected new title to be within content")
	}
}