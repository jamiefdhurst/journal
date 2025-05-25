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

	err := visits.CreateTable()

	if err != nil {
		t.Errorf("Expected no error creating table, got: %s", err)
	}
}

func TestVisits_FindByDateAndURL(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	visits := Visits{Container: container}

	db.Rows = &database.MockVisit_SingleRow{}
	visit := visits.FindByDateAndURL("2023-01-01", "/test")

	if visit.ID != 1 {
		t.Errorf("Expected visit ID to be 1, got %d", visit.ID)
	}
	if visit.URL != "/test" {
		t.Errorf("Expected visit URL to be /test, got %s", visit.URL)
	}
	if visit.Hits != 5 {
		t.Errorf("Expected visit hits to be 5, got %d", visit.Hits)
	}

	// Test with no visit found
	db.Rows = &database.MockRowsEmpty{}
	emptyVisit := visits.FindByDateAndURL("2023-01-01", "/nonexistent")

	if emptyVisit.ID != 0 {
		t.Errorf("Expected empty visit ID to be 0, got %d", emptyVisit.ID)
	}
}

func TestVisits_RecordVisit(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	visits := Visits{Container: container}

	db.Rows = &database.MockRowsEmpty{} // No existing visit
	db.Result = &database.MockResult{}

	err := visits.RecordVisit("/new-page")

	if err != nil {
		t.Errorf("Expected no error recording new visit, got: %s", err)
	}

	db.Rows = &database.MockVisit_SingleRow{} // Existing visit
	db.Result = &database.MockResult{}

	err = visits.RecordVisit("/test")

	if err != nil {
		t.Errorf("Expected no error updating existing visit, got: %s", err)
	}
}
