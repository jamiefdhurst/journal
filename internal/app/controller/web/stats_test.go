package web

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
	configuration.PostsPerPage = 25
	configuration.GoogleAnalyticsCode = "UA-123456"
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &Stats{}
	controller.DisableTracking()

	// Test with journals
	db.Rows = &database.MockJournal_MultipleRows{}
	request, _ := http.NewRequest("GET", "/stats", strings.NewReader(""))
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)

	if !strings.Contains(response.Content, "<h1>Stats</h1>") {
		t.Error("Expected stats page title to be displayed")
	}

	if !strings.Contains(response.Content, "<dt>Total Posts</dt>\n        <dd>2</dd>") {
		t.Error("Expected post count to be displayed")
	}

	if !strings.Contains(response.Content, "<dt>First Post Date</dt>") {
		t.Error("Expected first post date to be displayed")
	}

	if !strings.Contains(response.Content, "<dt>Posts Per Page</dt>\n        <dd>25</dd>") {
		t.Error("Expected custom posts per page setting to be displayed")
	}

	if !strings.Contains(response.Content, "<dt>Google Analytics</dt>\n        <dd>Enabled</dd>") {
		t.Error("Expected GA code to be displayed as enabled")
	}

	response.Reset()
	db.Rows = &database.MockRowsEmpty{}
	controller.Run(response, request)

	if !strings.Contains(response.Content, "<dt>Total Posts</dt>\n        <dd>0</dd>") {
		t.Error("Expected post count to be 0")
	}

	if !strings.Contains(response.Content, "<dt>First Post Date</dt>\n        <dd>No posts yet</dd>") {
		t.Error("Expected 'No posts yet' message for first post date")
	}
}
