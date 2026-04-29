package router

import (
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
)

func TestNewRouter(t *testing.T) {
    config := app.DefaultConfiguration()
    container := &app.Container{Configuration: config}
    rtr := NewRouter(container)

    if rtr == nil {
        t.Fatal("Expected router to be created, got nil")
    }
    if len(rtr.Routes) == 0 {
        t.Error("Expected routes to be registered")
    }
    if rtr.ErrorController == nil {
        t.Error("Expected error controller to be set")
    }
    if len(rtr.StaticPaths) == 0 {
        t.Error("Expected static paths to be configured")
    }
}
