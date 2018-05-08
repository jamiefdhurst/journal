package controller

import (
	"database/sql"
	"testing"

	"github.com/jamiefdhurst/journal/model"
)

type FakeDatabase struct{}

func (f *FakeDatabase) Close() {}

func (f *FakeDatabase) Connect() error {
	return nil
}

func (f *FakeDatabase) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (f *FakeDatabase) Query(sql string, args ...interface{}) (model.Rows, error) {
	return nil, nil
}

func TestInit(t *testing.T) {
	db := &FakeDatabase{}
	params := []string{
		"param1", "param2", "param3", "param4",
	}
	controller := Super{}
	controller.Init(db, params)
	if controller.Db != db || controller.Params[2] != "param3" {
		t.Error("Expected values were not passed into struct")
	}
}
