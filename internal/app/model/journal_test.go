package model

import (
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestJournal_GetDate(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"2018-05-10", "10/05/2018"},
		{"200-00-00", ""},
		{"", ""},
		{"0000-00-00", "00/00/0000"},
	}

	for _, table := range tables {
		j := Journal{Date: table.input}
		actual := j.GetDate()
		if actual != table.output {
			t.Errorf("Expected GetDate() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}

func TestJournal_GetEditableDate(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"2018-05-10", "2018-05-10"},
		{"2018-05-10EXTRATHINGS", "2018-05-10"},
		{"200-00-00", ""},
		{"", ""},
		{"0000-00-00", "0000-00-00"},
	}

	for _, table := range tables {
		j := Journal{Date: table.input}
		actual := j.GetEditableDate()
		if actual != table.output {
			t.Errorf("Expected GetEditableDate() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}

func TestJournals_CreateTable(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	js := Journals{Container: container}
	js.CreateTable()
	if db.Queries != 1 {
		t.Errorf("Expected 1 query to have been run")
	}
}

func TestJournals_FetchAll(t *testing.T) {

	// Test error
	db := &database.MockSqlite{}
	db.ErrorMode = true
	container := &app.Container{Db: db}
	js := Journals{Container: container}
	journals := js.FetchAll()
	if len(journals) > 0 {
		t.Errorf("Expected empty result set returned when error received")
	}

	// Test empty result
	db.ErrorMode = false
	db.Rows = &database.MockRowsEmpty{}
	journals = js.FetchAll()
	if len(journals) > 0 {
		t.Errorf("Expected empty result set returned")
	}

	// Test successful result
	db.Rows = &database.MockJournal_MultipleRows{}
	journals = js.FetchAll()
	if len(journals) < 2 || journals[0].ID != 1 || journals[1].Content != "Content 2" {
		t.Errorf("Expected 2 rows returned and with correct data")
	}
}

func TestJournals_FindBySlug(t *testing.T) {
	// Test error
	db := &database.MockSqlite{}
	db.ErrorMode = true
	container := &app.Container{Db: db}
	js := Journals{Container: container}
	journal := js.FindBySlug("example")
	if journal.ID > 0 {
		t.Errorf("Expected empty result set returned when error received")
	}

	// Test empty result
	db.ErrorMode = false
	db.Rows = &database.MockRowsEmpty{}
	journal = js.FindBySlug("example")
	if journal.ID > 0 {
		t.Errorf("Expected empty result set returned")
	}

	// Test successful result
	db.Rows = &database.MockJournal_SingleRow{}
	db.ExpectedArgument = "slug"
	journal = js.FindBySlug("slug")
	if journal.ID != 1 || journal.Content != "Content" {
		t.Errorf("Expected 1 row returned and with correct data")
	}

	// Test unexpected amount of rows
	db.Rows = &database.MockJournal_MultipleRows{}
	journal = js.FindBySlug("slug")
	if journal.ID > 0 {
		t.Errorf("Expected no rows when query returns more than one result")
	}
}

func TestJournals_Save(t *testing.T) {
	db := &database.MockSqlite{Result: &database.MockResult{}}
	gs := &database.MockGiphyExtractor{}
	container := &app.Container{Db: db}
	js := Journals{Container: container, Gs: gs}

	// Test with new Journal
	journal := js.Save(Journal{ID: 0, Title: "Testing"})
	if journal.ID != 1 || journal.Title != "Testing" {
		t.Error("Expected same Journal to have been returned with new ID")
	}

	// Test with same Journal
	journal = js.Save(Journal{ID: 2, Title: "Testing 2"})
	if journal.ID != 2 || journal.Title != "Testing 2" {
		t.Error("Expected same Journal to have been returned with new ID")
	}

	// Check Giphy calls
	if gs.CalledTimes != 2 {
		t.Error("Expected Giphy to have been called 2 times within test scope")
	}
}

func TestSlugify(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"A SIMPLE TITLE", "a-simple-title"},
		{"already-slugified", "already-slugified"},
		{"   ", "---"},
		{"lower cased", "lower-cased"},
		{"Special!!!Characters@$%^&*(", "special---characters-------"},
	}

	for _, table := range tables {
		actual := Slugify(table.input)
		if actual != table.output {
			t.Errorf("Expected Slugify() to produce result of '%s', got '%s'", table.output, actual)
		}
	}
}
