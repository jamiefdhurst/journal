package web

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestNew_Run(t *testing.T) {
	db := &database.MockSqlite{}
	db.Result = &database.MockResult{}
	db.Rows = &database.MockRowsEmpty{}
	configuration := app.DefaultConfiguration()
	configuration.EnableCreate = true
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &New{}

	// Display form
	request, _ := http.NewRequest("GET", "/new", strings.NewReader(""))
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if controller.Error || !strings.Contains(response.Content, "<form") {
		t.Error("Expected form to be shown")
	}
	if !strings.Contains(response.Content, "<title>Create New Post - Jamie's Journal</title>") {
		t.Error("Expected HTML title to be in place")
	}

	// Display error when cookie was set
	response.Reset()
	controller.Init(container, []string{"", "0"}, request)
	controller.Session.AddFlash("error")
	controller.SessionStore.Save(response)
	request, _ = http.NewRequest("GET", "/new", strings.NewReader(""))
	request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
	controller.Init(container, []string{"", "0"}, request)
	response.Reset()
	controller.Run(response, request)
	if !strings.Contains(response.Content, "<form") || !strings.Contains(response.Content, "div class=\"error\"") {
		t.Error("Expected form and error to be shown")
	}

	// Redirect if empty content on POST
	response.Reset()
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=&date=&content="))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 302 || response.Headers.Get("Location") != "/new" {
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
	db.Result = &database.MockResult{}
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=Title&date=2018-02-01&content=Test+again"))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
