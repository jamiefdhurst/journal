package model

import (
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestVisits_CreateTable(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	visits := Visits{Container: container}
	visits.CreateTable()
	if db.Queries != 1 {
		t.Error("Expected 1 query to have been run")
	}
}
