package controller

import (
	"testing"

	"github.com/jamiefdhurst/journal/model"
)

func TestInit(t *testing.T) {
	db := &model.MockDatabase{}
	params := []string{
		"param1", "param2", "param3", "param4",
	}
	controller := Super{}
	controller.Init(db, params)
	if controller.Db != db || controller.Params[2] != "param3" {
		t.Error("Expected values were not passed into struct")
	}
}
