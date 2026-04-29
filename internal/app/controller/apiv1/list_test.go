package apiv1

import (
    "net/http"
    "strings"
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/test/mocks/controller"
    "github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestList_Run(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Configuration: app.DefaultConfiguration(), Db: db}
    response := controller.NewMockResponse()
    controller := &List{}
    controller.DisableTracking()

    // Test showing all Journals
    db.EnableMultiMode()
    db.AppendResult(&database.MockPagination_Result{TotalResults: 2})
    db.AppendResult(&database.MockJournal_MultipleRows{})
    request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if !strings.Contains(response.Content, "Title 2") {
        t.Error("Expected all journals to be returned")
    }

    // Test with page parameter
    response.Reset()
    db.EnableMultiMode()
    db.AppendResult(&database.MockPagination_Result{TotalResults: 25})
    db.AppendResult(&database.MockJournal_MultipleRows{})
    request, _ = http.NewRequest("GET", "/?page=2", strings.NewReader(""))
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if !strings.Contains(response.Content, "Title 2") {
        t.Error("Expected journals to be returned for page 2")
    }
    if !strings.Contains(response.Content, `"current_page":2`) {
        t.Error("Expected pagination to reflect page 2")
    }
}
