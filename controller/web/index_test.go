package web

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

func TestIndex_Run(t *testing.T) {
	database := &model.MockSqlite{}
	response := controller.NewMockResponse()
	controller := &Index{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	controller.Init(database, []string{"", "0"})
	database.Rows = &model.MockJournal_MultipleRows{}
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
