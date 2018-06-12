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

func TestSingle_Run(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	response := &controller.MockResponse{}
	response.Reset()
	controller := &Single{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET
	controller.Init(container, []string{"", "0"})
	db.Rows = &database.MockRowsEmpty{}
	request := &http.Request{Method: "GET"}
	controller.Run(response, request)
	if response.StatusCode != 404 {
		t.Error("Expected 404 error when journal not found")
	}

	// Test return
	response.Reset()
	request, _ = http.NewRequest("GET", "/slug", strings.NewReader(""))
	db.Rows = &database.MockJournal_SingleRow{}
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected content to be returned")
	}
}
