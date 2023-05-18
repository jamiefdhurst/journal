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
	configuration := app.DefaultConfiguration()
	configuration.EnableEdit = true
	container := &app.Container{Configuration: configuration, Db: db}
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
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

	// Display error when cookie was set
	response.Reset()
	controller.Init(container, []string{"", "0"}, request)
	controller.Session.AddFlash("error")
	controller.SessionStore.Save(response)
	request, _ = http.NewRequest("GET", "/test/edit", strings.NewReader(""))
	request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
	response.Reset()
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if !controller.Error || !strings.Contains(response.Content, "div class=\"error\"") {
		t.Error("Expected error to be shown in form")
	}

	// Display no error
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug/edit", strings.NewReader(""))
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Error = false
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if controller.Error || strings.Contains(response.Content, "div class=\"error\"") {
		t.Error("Expected no error to be shown in form")
	}
	if !strings.Contains(response.Content, "<title>Edit Title - Jamie's Journal</title>") {
		t.Error("Expected HTML title to be in place")
	}

	// Redirect if empty content on POST
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("title=&date=&content="))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/slug/edit" {
		t.Error("Expected redirect back to same page")
	}

	// Validate error cookie on redirect
	request, _ = http.NewRequest("GET", "/", strings.NewReader(""))
	request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
	controller.Init(container, []string{"", "0"}, request)
	flash := controller.Session.GetFlash()
	if flash == nil || flash[0] != "error" {
		t.Error("Expected cookie to contain error value")
	}

	// Redirect on success
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("title=Title&date=2018-02-01&content=Test+again"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/" {
		t.Error("Expected redirect back to home with saved banner shown")
	}

	// Validate saved cookie on redirect
	request, _ = http.NewRequest("GET", "/", strings.NewReader(""))
	request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
	controller.Init(container, []string{"", "0"}, request)
	flash = controller.Session.GetFlash()
	if flash == nil || flash[0] != "saved" {
		t.Error("Expected cookie to contain saved value")
	}
}
