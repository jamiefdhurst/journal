package apiv1

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/model"
)

type FakeSqlite struct {
	Closed           bool
	Connected        bool
	ErrorAtQuery     int
	ErrorMode        bool
	ExpectedArgument string
	Queries          int
	Result           sql.Result
	Rows             model.Rows
}

func (f *FakeSqlite) Close() {
	f.Closed = true
}

func (f *FakeSqlite) Connect() error {
	f.Connected = true
	return nil
}

func (f *FakeSqlite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	f.Queries++
	if f.ErrorMode || f.ErrorAtQuery == f.Queries {
		return nil, errors.New("Simulating error")
	}
	if f.ExpectedArgument != "" && !f.inArgs(args) {
		return nil, errors.New("Expected " + f.ExpectedArgument + " in query")
	}
	return f.Result, nil
}

func (f *FakeSqlite) Query(sql string, args ...interface{}) (model.Rows, error) {
	f.Queries++
	if f.ErrorMode || f.ErrorAtQuery == f.Queries {
		return nil, errors.New("Simulating error")
	}
	if f.ExpectedArgument != "" && !f.inArgs(args) {
		return nil, errors.New("Expected " + f.ExpectedArgument + " in query")
	}
	return f.Rows, nil
}

func (f *FakeSqlite) inArgs(slice []interface{}) bool {
	for _, v := range slice {
		if v.(string) == f.ExpectedArgument {
			return true
		}
	}
	return false
}

type FakeResponse struct {
	Content    string
	Headers    http.Header
	StatusCode int
}

func (f *FakeResponse) Header() http.Header {
	return f.Headers
}

func (f *FakeResponse) Reset() {
	f.Content = ""
	f.Headers = make(http.Header)
	f.StatusCode = 200
}

func (f *FakeResponse) Write(b []byte) (int, error) {
	f.Content = strings.Join([]string{f.Content, string(b[:])}, "")
	return len(b), nil
}

func (f *FakeResponse) WriteHeader(statusCode int) {
	f.StatusCode = statusCode
}

type MockEmptyRows struct{}

func (m *MockEmptyRows) Close() error {
	return nil
}

func (m *MockEmptyRows) Columns() ([]string, error) {
	return []string{}, nil
}

func (m *MockEmptyRows) Next() bool {
	return false
}

func (m *MockEmptyRows) Scan(dest ...interface{}) error {
	return nil
}

type MockJournalSingleRow struct {
	MockEmptyRows
	RowNumber int
}

func (m *MockJournalSingleRow) Next() bool {
	m.RowNumber++
	if m.RowNumber < 2 {
		return true
	}
	return false
}

func (m *MockJournalSingleRow) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "slug"
		*dest[2].(*string) = "Title"
		*dest[3].(*string) = "2018-02-01"
		*dest[4].(*string) = "Content"
	}
	return nil
}

func TestUpdate_Run(t *testing.T) {
	database := &FakeSqlite{}
	response := &FakeResponse{}
	response.Reset()
	controller := &Update{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test not found/error with GET/POST
	controller.Init(database, []string{"", "0"})
	database.Rows = &MockEmptyRows{}
	request := &http.Request{Method: "POST"}
	controller.Run(response, request)
	if response.StatusCode != 404 {
		t.Error("Expected 404 error when journal not found")
	}

	// Test for bad request on invalid JSON
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("{\"not\":\"valid\":\"json\"}"))
	request.Header.Add("Content-Type", "application/json")
	database.Rows = &MockJournalSingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when invalid JSON provided")
	}

	// Test Journal is retrieved on save
	response.Reset()
	request, _ = http.NewRequest("POST", "/slug/edit", strings.NewReader("{\"title\":\"Something New\",\"date\":\"2018-01-01\",\"content\":\"New\"}"))
	request.Header.Add("Content-Type", "application/json")
	database.Rows = &MockJournalSingleRow{}
	controller.Run(response, request)
	if response.StatusCode != 200 || !strings.Contains(response.Content, "Something New") {
		t.Error("Expected new title to be within content")
	}
}
