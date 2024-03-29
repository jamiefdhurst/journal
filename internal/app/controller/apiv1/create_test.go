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

func TestCreate_Run(t *testing.T) {
	db := &database.MockSqlite{}
	db.Result = &database.MockResult{}
	db.Rows = &database.MockRowsEmpty{}
	container := &app.Container{Configuration: app.DefaultConfiguration(), Db: db}
	response := controller.NewMockResponse()
	response.Reset()
	controller := &Create{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test forbidden
	container.Configuration.EnableCreate = false
	request, _ := http.NewRequest("POST", "/new", strings.NewReader("{\"not\":\"valid\":\"json\"}"))
	request.Header.Add("Content-Type", "application/json")
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 403 {
		t.Error("Expected 403 error when creation is disabled")
	}

	// Test invalid JSON
	container.Configuration.EnableCreate = true
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("{\"not\":\"valid\":\"json\"}"))
	request.Header.Add("Content-Type", "application/json")
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when invalid JSON provided")
	}

	// Test missing JSON
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("{\"title\":\"only\"}"))
	request.Header.Add("Content-Type", "application/json")
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when missing JSON provided")
	}

	// Test Journal is retrieved on save
	response.Reset()
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("{\"title\":\"Something New\",\"date\":\"2018-01-01\",\"content\":\"New\"}"))
	request.Header.Add("Content-Type", "application/json")
	db.Result = &database.MockResult{}
	controller.Run(response, request)
	if response.StatusCode != 201 || !strings.Contains(response.Content, "Something New") {
		t.Error("Expected new title to be within content")
	}
}
