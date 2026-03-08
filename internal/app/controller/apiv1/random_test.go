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
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	random := &Random{}
	random.DisableTracking()

	// Test with a journal found
	response := controller.NewMockResponse()
	db.Rows = &database.MockJournal_SingleRow{}
	request, _ := http.NewRequest("GET", "/api/v1/post/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	random.Run(response, request)

	if response.Headers.Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json content type, got %s", response.Headers.Get("Content-Type"))
	}
	if !strings.Contains(response.Content, `"url"`) || !strings.Contains(response.Content, `"Title"`) {
		t.Errorf("Expected JSON response with journal data, got: %s", response.Content)
	}

	// Test with no journal found
	response = controller.NewMockResponse()
	db.Rows = &database.MockRowsEmpty{}
	request, _ = http.NewRequest("GET", "/api/v1/post/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	random.Run(response, request)

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 not found, got status %d", response.StatusCode)
	}
}
