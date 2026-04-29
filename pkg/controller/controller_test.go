package controller

import (
    "net/http"
    "strings"
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
    mockCtrl "github.com/jamiefdhurst/journal/test/mocks/controller"
    mockDb "github.com/jamiefdhurst/journal/test/mocks/database"
)

type BlankInterface struct{}

func TestInit(t *testing.T) {
    t.Run("Init with blank interface", func(t *testing.T) {
        container := BlankInterface{}
        params := []string{
            "param1", "param2", "param3", "param4",
        }
        controller := Super{}
        request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
        request.Host = "foobar.com"
        controller.Init(container, params, request)
        if controller.Container() != container || controller.Params()[2] != "param3" || controller.Host() != "foobar.com" {
            t.Error("Expected values were not passed into struct")
        }
    })

    t.Run("Init with app container and session config", func(t *testing.T) {
        container := &app.Container{
            Configuration: app.Configuration{
                SessionKey:       "12345678901234567890123456789012",
                SessionName:      "test-session",
                CookieDomain:     "example.com",
                CookieMaxAge:     3600,
                CookieSecure:     true,
                CookieHTTPOnly:   true,
            },
        }
        params := []string{"param1", "param2"}
        controller := Super{}
        request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
        request.Host = "test.com"

        controller.Init(container, params, request)

        if controller.Container() != container {
            t.Error("Expected container to be set")
        }
        if controller.Host() != "test.com" {
            t.Error("Expected host to be set")
        }
        if len(controller.Params()) != 2 {
            t.Error("Expected params to be set")
        }
        if controller.Session() == nil {
            t.Error("Expected session to be initialized")
        }
    })
}

func TestDisableTracking(t *testing.T) {
    db := &mockDb.MockSqlite{}
    db.Rows = &mockDb.MockRowsEmpty{}
    db.Result = &mockDb.MockResult{}
    container := &app.Container{
        Configuration: app.Configuration{
            SessionKey:  "12345678901234567890123456789012",
            SessionName: "test-session",
        },
        Db: db,
    }
    c := Super{}
    c.DisableTracking()
    request, _ := http.NewRequest("GET", "/test-path", strings.NewReader(""))
    c.Init(container, []string{}, request)
    if db.Queries != 0 {
        t.Errorf("Expected no DB queries when tracking disabled, got %d", db.Queries)
    }
}

func TestTrackVisit(t *testing.T) {
    db := &mockDb.MockSqlite{}
    db.Rows = &mockDb.MockRowsEmpty{}
    db.Result = &mockDb.MockResult{}
    container := &app.Container{
        Configuration: app.Configuration{
            SessionKey:  "12345678901234567890123456789012",
            SessionName: "test-session",
        },
        Db: db,
    }
    c := Super{}
    request, _ := http.NewRequest("GET", "/test-path", strings.NewReader(""))
    c.Init(container, []string{}, request)
    if db.Queries == 0 {
        t.Error("Expected DB queries for visit tracking when not disabled")
    }
}

func TestTrackVisit_NilContainer(t *testing.T) {
    c := Super{}
    request, _ := http.NewRequest("GET", "/test-path", strings.NewReader(""))
    // Should not panic with nil container
    c.trackVisit(request)
}

func TestSaveSession(t *testing.T) {
    // Without a session store — should be a no-op
    c := Super{}
    response := mockCtrl.NewMockResponse()
    c.SaveSession(response)
    if response.Headers.Get("Set-Cookie") != "" {
        t.Error("Expected no Set-Cookie header when there is no session store")
    }

    // With a session store
    container := &app.Container{
        Configuration: app.Configuration{
            SessionKey:  "12345678901234567890123456789012",
            SessionName: "test-session",
            CookieMaxAge: 3600,
            CookieHTTPOnly: true,
        },
    }
    c2 := Super{}
    c2.DisableTracking()
    request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
    c2.Init(container, []string{}, request)
    c2.Session().Set("key", "value")
    response2 := mockCtrl.NewMockResponse()
    c2.SaveSession(response2)
    if response2.Headers.Get("Set-Cookie") == "" {
        t.Error("Expected Set-Cookie header after saving session")
    }
}
