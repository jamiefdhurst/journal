package apiv1

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestList_Run(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	response := &controller.MockResponse{}
	response.Reset()
	controller := &List{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	db.Rows = &database.MockJournal_MultipleRows{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title 2") {
		t.Error("Expected all journals to be returned")
	}
}
