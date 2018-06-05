package apiv1

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestIndex_Run(t *testing.T) {
	db := &database.MockSqlite{}
	response := &controller.MockResponse{}
	response.Reset()
	controller := &List{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	controller.Init(db, []string{"", "0"})
	db.Rows = &database.MockJournal_MultipleRows{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title 2") {
		t.Error("Expected all journals to be returned")
	}
}
