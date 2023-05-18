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

func TestIndex_Run(t *testing.T) {
	db := &database.MockSqlite{}
	configuration := app.DefaultConfiguration()
	configuration.ArticlesPerPage = 2
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &Index{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	db.EnableMultiMode()
	db.AppendResult(&database.MockPagination_Result{TotalResults: 2})
	db.AppendResult(&database.MockJournal_MultipleRows{})
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title 2") {
		t.Error("Expected all journals to be displayed on screen")
	}
	if !strings.Contains(response.Content, "<title>Jamie's Journal</title>") {
		t.Error("Expected default HTML title to be in place")
	}
	if !strings.Contains(response.Content, "<meta name=\"description\" content=\"A private journal") {
		t.Error("Expected default HTML description to be in place")
	}

	// Test pagination
	db.EnableMultiMode()
	db.AppendResult(&database.MockPagination_Result{TotalResults: 4})
	db.AppendResult(&database.MockJournal_MultipleRows{})
	request, _ = http.NewRequest("GET", "/?page=2", strings.NewReader(""))
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title 2") {
		t.Error("Expected all journals to be displayed on screen")
	}
	if !strings.Contains(response.Content, `class="current"`) {
		t.Error("Expected pagination to work")
	}

	// Test saved banner showing
	response.Reset()
	db.AppendResult(&database.MockPagination_Result{TotalResults: 2})
	db.AppendResult(&database.MockJournal_MultipleRows{})
	controller.Init(container, []string{"", "0"}, request)
	controller.Session.AddFlash("saved")
	controller.SessionStore.Save(response)
	request, _ = http.NewRequest("GET", "/", strings.NewReader(""))
	request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Journal saved") {
		t.Error("Expected saved banner to be displayed on screen")
	}
	response.Reset()
	db.AppendResult(&database.MockPagination_Result{TotalResults: 2})
	db.AppendResult(&database.MockJournal_MultipleRows{})
	request, _ = http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Run(response, request)
	if strings.Contains(response.Content, "Journal saved") {
		t.Error("Expected saved banner to be hidden, but it is showing")
	}
}
