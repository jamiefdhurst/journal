package controller

import (
	"testing"
)

type BlankInterface struct{}

func TestInit(t *testing.T) {
	container := BlankInterface{}
	params := []string{
		"param1", "param2", "param3", "param4",
	}
	controller := Super{}
	controller.Init(container, params)
	if controller.Container != container || controller.Params[2] != "param3" {
		t.Error("Expected values were not passed into struct")
	}
}
