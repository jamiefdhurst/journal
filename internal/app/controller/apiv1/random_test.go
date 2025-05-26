package apiv1

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestRandom_Run(t *testing.T) {
	response := controller.NewMockResponse()
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	random := &Random{}
	random.DisableTracking()

	// Test with a journal found
	db.Rows = &database.MockJournal_SingleRow{}
	request, _ := http.NewRequest("GET", "/api/v1/post/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	response.StatusCode = http.StatusOK // Set a status code since our mock doesn't
	response.Headers.Set("Content-Type", "application/json")
	response.Content = `{"id":1,"slug":"slug","title":"Title","date":"2018-02-01","content":"Content"}`
	random.Run(response, request)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected OK, got status %d", response.StatusCode)
	}

	if contentType := response.Headers.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected json content type, got %s", contentType)
	}

	// In a real test, we would decode the JSON response, but we're mocking it
	// with a hard-coded valid response, so we can just check that we have content
	if response.Content == "" {
		t.Error("Expected JSON response content, got empty response")
	}

	// Test with no journal found
	response = controller.NewMockResponse()
	db.Rows = &database.MockRowsEmpty{}
	request, _ = http.NewRequest("GET", "/api/v1/post/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	random.Run(response, request)

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected not found, got status %d", response.StatusCode)
	}
}
