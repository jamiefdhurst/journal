package web

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
	
	// Test with a journal found
	db.Rows = &database.MockJournal_SingleRow{}
	request, _ := http.NewRequest("GET", "/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	random.Run(response, request)
	
	if response.StatusCode != http.StatusFound {
		t.Errorf("Expected redirect, got status %d", response.StatusCode)
	}
	
	if location := response.Headers.Get("Location"); location != "/slug" {
		t.Errorf("Expected redirect to /slug, got %s", location)
	}
	
	// Test with no journal found
	response = controller.NewMockResponse()
	db.Rows = &database.MockRowsEmpty{}
	request, _ = http.NewRequest("GET", "/random", strings.NewReader(""))
	random.Init(container, []string{}, request)
	random.Run(response, request)
	
	if response.StatusCode != http.StatusFound {
		t.Errorf("Expected redirect, got status %d", response.StatusCode)
	}
	
	if location := response.Headers.Get("Location"); location != "/" {
		t.Errorf("Expected redirect to /, got %s", location)
	}
}