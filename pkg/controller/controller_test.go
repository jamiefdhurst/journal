package controller

import (
    "net/http"
    "strings"
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
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
        if controller.sessionStore == nil {
            t.Error("Expected session store to be initialized")
        }
        if controller.session == nil {
            t.Error("Expected session to be initialized")
        }
    })
}
