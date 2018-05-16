package apiv1

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

func TestNew_Run(t *testing.T) {
	database := &model.MockSqlite{}
	response := controller.NewMockResponse()
	response.Reset()
	controller := &Create{}
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Test invalid JSON
	controller.Init(database, []string{"", "0"})
	request, _ := http.NewRequest("POST", "/new", strings.NewReader("{\"not\":\"valid\":\"json\"}"))
	request.Header.Add("Content-Type", "application/json")
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when invalid JSON provided")
	}

	// Test missing JSON
	controller.Init(database, []string{"", "0"})
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("{\"title\":\"only\"}"))
	request.Header.Add("Content-Type", "application/json")
	controller.Run(response, request)
	if response.StatusCode != 400 {
		t.Error("Expected 400 error when missing JSON provided")
	}

	// Test Journal is retrieved on save
	response.Reset()
	request, _ = http.NewRequest("POST", "/new", strings.NewReader("{\"title\":\"Something New\",\"date\":\"2018-01-01\",\"content\":\"New\"}"))
	request.Header.Add("Content-Type", "application/json")
	database.Result = &model.MockResult{}
	controller.Run(response, request)
	if response.StatusCode != 200 || !strings.Contains(response.Content, "Something New") {
		t.Error("Expected new title to be within content")
	}
}
