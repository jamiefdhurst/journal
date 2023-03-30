package web

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestEdit_Run(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	response := controller.NewMockResponse()
	controller := &Edit{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET/POST
	db.Rows = &database.MockRowsEmpty{}
	request := &http.Request{Method: "GET"}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

	response.Reset()
	request = &http.Request{Method: "POST"}
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

	// Display error when passed through
	response.Reset()
	request, _ = http.NewRequest("GET", "/test/edit?error=1", strings.NewReader(""))
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if !controller.Error || !strings.Contains(response.Content, "div class=\"error\"") {
		t.Error("Expected error to be shown in form")
	}

	// Display no error
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug/edit", strings.NewReader(""))
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Error = false
	controller.Run(response, request)
	if controller.Error || strings.Contains(response.Content, "div class=\"error\"") {
		t.Error("Expected no error to be shown in form")
	}

	// Redirect if empty content on POST
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("title=&date=&content="))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/slug/edit?error=1" {
		t.Error("Expected redirect back to same page")
	}

	// Redirect on success
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("title=Title&date=2018-02-01&content=Test+again"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/?saved=1" {
		t.Error("Expected redirect back to home with saved flag")
	}
}
