package web

import (
    "net/http"
    "strings"
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/test/mocks/controller"
    "github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestNew_Run(t *testing.T) {
    db := &database.MockSqlite{}
    db.Result = &database.MockResult{}
    db.Rows = &database.MockRowsEmpty{}
    configuration := app.DefaultConfiguration()
    configuration.EnableCreate = true
    configuration.SessionKey = "12345678901234567890123456789012"
    container := &app.Container{Configuration: configuration, Db: db}
    response := controller.NewMockResponse()
    controller := &New{}
    controller.DisableTracking()

    // Test forbidden when creation is disabled
    container.Configuration.EnableCreate = false
    request, _ := http.NewRequest("GET", "/new", strings.NewReader(""))
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if response.StatusCode != 404 || !strings.Contains(response.Content, "Page Not Found") {
        t.Error("Expected error page when creation is disabled")
    }
    container.Configuration.EnableCreate = true

    // Display form
    response.Reset()
    request, _ = http.NewRequest("GET", "/new", strings.NewReader(""))
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if !strings.Contains(response.Content, "<form") {
        t.Error("Expected form to be shown")
    }
    if !strings.Contains(response.Content, "<title>Create New Post - A Fantastic Journal</title>") {
        t.Error("Expected HTML title to be in place")
    }

    // Display error when cookie was set
    response.Reset()
    controller.Init(container, []string{"", "0"}, request)
    controller.Session().AddFlash("error")
    controller.SaveSession(response)
    request, _ = http.NewRequest("GET", "/new", strings.NewReader(""))
    request.Header.Add("Cookie", response.Headers.Get("Set-Cookie"))
    controller.Init(container, []string{"", "0"}, request)
    response.Reset()
    controller.Run(response, request)
    if !strings.Contains(response.Content, "<form") || !strings.Contains(response.Content, "div class=\"error\"") {
        t.Error("Expected form and error to be shown")
    }

    // Redirect if empty content on POST
    response.Reset()
    request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=&date=&content="))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if response.StatusCode != 302 || response.Headers.Get("Location") != "/new" {
        t.Error("Expected redirect back to same page")
    }

    // Test form data preservation when validation fails
    response.Reset()
    // Create a new controller instance for this test
    prevController := &New{}
    // Submit a form with a missing field (date is empty)
    request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=Test+Title&date=&content=Test+Content"))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    prevController.Init(container, []string{"", "0"}, request)

    // Verify form data is saved in session
    prevController.Run(response, request)
    if response.StatusCode != 302 || response.Headers.Get("Location") != "/new" {
        t.Error("Expected redirect back to new page")
    }

    // Check if form_data was set in the session
    formData := prevController.Session().Get("form_data")
    if formData == nil {
        t.Error("Expected form_data to be set in session")
    } else {
        // Cast and verify form data values
        formMap := formData.(map[string]string)
        if formMap["title"] != "Test Title" {
            t.Errorf("Expected title to be 'Test Title', got '%s'", formMap["title"])
        }
        if formMap["content"] != "Test Content" {
            t.Errorf("Expected content to be 'Test Content', got '%s'", formMap["content"])
        }
        if formMap["date"] != "" {
            t.Errorf("Expected date to be empty, got '%s'", formMap["date"])
        }
    }

    // Redirect on success
    response.Reset()
    db.Result = &database.MockResult{}
    request, _ = http.NewRequest("POST", "/new", strings.NewReader("title=Title&date=2018-02-01&content=Test+again"))
    request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    controller.Init(container, []string{"", "0"}, request)
    controller.Run(response, request)
    if response.StatusCode != 302 || response.Headers.Get("Location") != "/" {
        t.Error("Expected redirect back to home with saved banner shown")
    }

}
