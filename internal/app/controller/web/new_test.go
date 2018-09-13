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

func TestNew_Run(t *testing.T) {
	db := &database.MockSqlite{}
	db.Result = &database.MockResult{}
	db.Rows = &database.MockRowsEmpty{}
	container := &app.Container{Db: db}
	response := controller.NewMockResponse()
	controller := &New{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Display form
	controller.Init(container, []string{"", "0"})
	request, _ := http.NewRequest("GET", "/new", strings.NewReader(""))
	controller.Run(response, request)
	if controller.Error || !strings.Contains(response.Content, "<form") {
		t.Error("Expected form to be shown")
	}

	// Display error
	response.Reset()
	request, _ = http.NewRequest("GET", "/new?error=1", strings.NewReader(""))
	controller.Run(response, request)
	if !strings.Contains(response.Content, "<form") || !strings.Contains(response.Content, "error") {
		t.Error("Expected form and error to be shown")
	}

	// Redirect if empty content on POST
	response.Reset()
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=&date=&content="))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/new?error=1" {
		t.Error("Expected redirect back to same page")
	}

	// Redirect on success
	response.Reset()
	db.Result = &database.MockResult{}
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=Title&date=2018-02-01&content=Test+again"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/?saved=1" {
		t.Error("Expected redirect back to home with saved flag")
	}
}
