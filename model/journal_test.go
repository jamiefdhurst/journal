package model

import (
	"testing"
)

type MockGiphyExtractor struct {
	CalledTimes int
}

func (m *MockGiphyExtractor) ExtractContentsAndSearchAPI(s string) string {
	m.CalledTimes++
	return s
}

type MockJournalMultipleRows struct {
	MockEmptyRows
	RowNumber int
}

func (m *MockJournalMultipleRows) Next() bool {
	m.RowNumber++
	if m.RowNumber < 3 {
		return true
	}
	return false
}

func (m *MockJournalMultipleRows) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "slug"
		*dest[2].(*string) = "Title"
		*dest[3].(*string) = "2018-02-01"
		*dest[4].(*string) = "Content"
	} else if m.RowNumber == 2 {
		*dest[0].(*int) = 2
		*dest[1].(*string) = "slug-2"
		*dest[2].(*string) = "Title 2"
		*dest[3].(*string) = "2018-03-01"
		*dest[4].(*string) = "Content 2"
	}
	return nil
}

type MockJournalSingleRow struct {
	MockEmptyRows
	RowNumber int
}

func (m *MockJournalSingleRow) Next() bool {
	m.RowNumber++
	if m.RowNumber < 2 {
		return true
	}
	return false
}

func (m *MockJournalSingleRow) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "slug"
		*dest[2].(*string) = "Title"
		*dest[3].(*string) = "2018-02-01"
		*dest[4].(*string) = "Content"
	}
	return nil
}

type MockJournalSaveResult struct{}

func (m *MockJournalSaveResult) LastInsertId() (int64, error) {
	return 10, nil
}

func (m *MockJournalSaveResult) RowsAffected() (int64, error) {
	return 0, nil
}

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
	database := &FakeSqlite{}
	js := Journals{Db: database}
	js.CreateTable()
	if database.Queries != 1 {
		t.Errorf("Expected 1 query to have been run")
	}
}

func TestJournals_FetchAll(t *testing.T) {

	// Test error
	database := &FakeSqlite{}
	database.ErrorMode = true
	js := Journals{Db: database}
	journals := js.FetchAll()
	if len(journals) > 0 {
		t.Errorf("Expected empty result set returned when error received")
	}

	// Test empty result
	database.ErrorMode = false
	database.Rows = &MockEmptyRows{}
	journals = js.FetchAll()
	if len(journals) > 0 {
		t.Errorf("Expected empty result set returned")
	}

	// Test successful result
	database.Rows = &MockJournalMultipleRows{}
	journals = js.FetchAll()
	if len(journals) < 2 || journals[0].ID != 1 || journals[1].Content != "Content 2" {
		t.Errorf("Expected 2 rows returned and with correct data")
	}
}

func TestJournals_FindBySlug(t *testing.T) {
	// Test error
	database := &FakeSqlite{}
	database.ErrorMode = true
	js := Journals{Db: database}
	journal := js.FindBySlug("example")
	if journal.ID > 0 {
		t.Errorf("Expected empty result set returned when error received")
	}

	// Test empty result
	database.ErrorMode = false
	database.Rows = &MockEmptyRows{}
	journal = js.FindBySlug("example")
	if journal.ID > 0 {
		t.Errorf("Expected empty result set returned")
	}

	// Test successful result
	database.Rows = &MockJournalSingleRow{}
	database.ExpectedArgument = "slug"
	journal = js.FindBySlug("slug")
	if journal.ID != 1 || journal.Content != "Content" {
		t.Errorf("Expected 1 row returned and with correct data")
	}

	// Test unexpected amount of rows
	database.Rows = &MockJournalMultipleRows{}
	journal = js.FindBySlug("slug")
	if journal.ID > 0 {
		t.Errorf("Expected no rows when query returns more than one result")
	}
}

func TestJournals_Save(t *testing.T) {
	database := &FakeSqlite{Result: &MockJournalSaveResult{}}
	gs := &MockGiphyExtractor{}
	js := Journals{Db: database, Gs: gs}

	// Test with new Journal
	journal := js.Save(Journal{ID: 0, Title: "Testing"})
	if journal.ID != 10 || journal.Title != "Testing" {
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
