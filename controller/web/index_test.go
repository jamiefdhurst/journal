package web

import (
	"net/http"
	"os"
	"strings"
	"testing"
)

type MockJournalMultipleRows struct {
	MockEmptyRows
	RowNumber int
}

func (m *MockJournalMultipleRows) Next() bool {
	m.RowNumber++
	if m.RowNumber < 3 {
		return true
	}
	return false
}

func (m *MockJournalMultipleRows) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "slug"
		*dest[2].(*string) = "Title"
		*dest[3].(*string) = "2018-02-01"
		*dest[4].(*string) = "Content"
	} else if m.RowNumber == 2 {
		*dest[0].(*int) = 2
		*dest[1].(*string) = "slug-2"
		*dest[2].(*string) = "Title 2"
		*dest[3].(*string) = "2018-03-01"
		*dest[4].(*string) = "Content 2"
	}
	return nil
}
func TestIndex_Run(t *testing.T) {
	database := &FakeSqlite{}
	response := &FakeResponse{}
	response.Reset()
	controller := &Index{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test showing all Journals
	controller.Init(database, []string{"", "0"})
	database.Rows = &MockJournalMultipleRows{}
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
