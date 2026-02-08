package apiv1

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestStats_Run(t *testing.T) {
	db := &database.MockSqlite{}
	configuration := app.DefaultConfiguration()
	configuration.PostsPerPage = 25                 // Custom setting
	configuration.GoogleAnalyticsCode = "UA-123456" // Custom GA code
	container := &app.Container{Configuration: configuration, Db: db}
	response := &controller.MockResponse{}
	response.Reset()
	controller := &Stats{}
	controller.DisableTracking()

	// Test with journals
	db.Rows = &database.MockJournal_MultipleRows{}
	request := &http.Request{Method: "GET"}
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)

	if response.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}
	if response.Headers.Get("Content-Type") != "application/json" {
		t.Error("Expected JSON content type")
	}

	if !strings.Contains(response.Content, "count\":2,") {
		t.Errorf("Expected post count to be 2, got response %s", response.Content)
	}
	if !strings.Contains(response.Content, "posts_per_page\":25,") {
		t.Errorf("Expected posts per page to be 25, got response %s", response.Content)
	}
	if !strings.Contains(response.Content, "google_analytics\":true") {
		t.Error("Expected Google Analytics to be enabled")
	}

	// Now test with no journals
	response.Reset()
	db.Rows = &database.MockRowsEmpty{}
	controller.Run(response, request)

	if !strings.Contains(response.Content, "count\":0}") {
		t.Errorf("Expected post count to be 0, got response %s", response.Content)
	}
	if strings.Contains(response.Content, "first_post_date") {
		t.Error("Expected first_post_date to be omitted when no posts exist")
	}
}
