package controller

import (
	"net/http"
	"strings"
	"testing"
)

type BlankInterface struct{}

func TestInit(t *testing.T) {
	container := BlankInterface{}
	params := []string{
		"param1", "param2", "param3", "param4",
	}
	controller := Super{}
	request, _ := http.NewRequest("GET", "/", strings.NewReader(""))
	controller.Init(container, params, request)
	if controller.Container != container || controller.Params[2] != "param3" {
		t.Error("Expected values were not passed into struct")
	}
}
