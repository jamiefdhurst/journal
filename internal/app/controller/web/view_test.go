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

func TestView_Run(t *testing.T) {
	db := &database.MockSqlite{}
	configuration := app.DefaultConfiguration()
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &View{}

	// Test not found/error with GET/POST
	db.Rows = &database.MockRowsEmpty{}
	request := &http.Request{Method: "GET"}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
		t.Error("Expected 404 error when journal not found")
	}

	// Display no error
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug", strings.NewReader(""))
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if strings.Contains(response.Content, "div class=\"error\"") || !strings.Contains(response.Content, "Content") {
		t.Error("Expected no error to be shown in page")
	}
	if !strings.Contains(response.Content, "<title>Title - Jamie's Journal</title>") {
		t.Error("Expected HTML title to be in place")
	}

	// Display prev & next strings
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug", strings.NewReader(""))
	db.EnableMultiMode()
	db.AppendResult(&database.MockJournal_SingleRow{})
	db.AppendResult(&database.MockJournal_SingleRow{})
	db.AppendResult(&database.MockJournal_SingleRow{})
	controller.Run(response, request)
	if !strings.Contains(response.Content, ">Previous<") || !strings.Contains(response.Content, ">Next<") {
		t.Error("Expected previous and next links to be shown in page")
	}
}
