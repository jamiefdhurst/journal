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
	container := &app.Container{Db: db}
	response := controller.NewMockResponse()
	controller := &Index{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	controller.Init(container, []string{"", "0"})
	db.Rows = &database.MockJournal_MultipleRows{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title 2") {
		t.Error("Expected all journals to be displayed on screen")
	}

	// Test saved banner showing
	response.Reset()
	request, _ = http.NewRequest("GET", "/?saved=1", strings.NewReader(""))
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Journal saved") {
		t.Error("Expected saved banner to be displayed on screen")
	}
}
