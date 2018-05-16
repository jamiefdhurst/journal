package model

import (
	"database/sql"
	"errors"
)

// MockDatabase Mock for model.Database
type MockDatabase struct{}

// Close Mock the close method
func (m *MockDatabase) Close() {}

// Connect Mock the connect method
func (m *MockDatabase) Connect() error {
	return nil
}

// Exec Mock empty exec
func (m *MockDatabase) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

// Query Mock empty query
func (m *MockDatabase) Query(sql string, args ...interface{}) (Rows, error) {
	return nil, nil
}

// MockGiphyExtractor Mock the Giphy Extractor interface
type MockGiphyExtractor struct {
	CalledTimes int
}

// ExtractContentsAndSearchAPI Mock the full call
func (m *MockGiphyExtractor) ExtractContentsAndSearchAPI(s string) string {
	m.CalledTimes++
	return s
}

// MockGiphy_SingleRow Mock single row for the Giphy API
type MockGiphy_SingleRow struct {
	MockRowsEmpty
	RowNumber int
}

// Next Mock 1 row
func (m *MockGiphy_SingleRow) Next() bool {
	m.RowNumber++
	if m.RowNumber < 2 {
		return true
	}
	return false
}

// Scan Return the data
func (m *MockGiphy_SingleRow) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "API123456"
	}
	return nil
}

// MockJournal_MultipleRows Mock multiple rows returned for a Journal
type MockJournal_MultipleRows struct {
	MockRowsEmpty
	RowNumber int
}

// Next Mock 2 rows
func (m *MockJournal_MultipleRows) Next() bool {
	m.RowNumber++
	if m.RowNumber < 3 {
		return true
	}
	return false
}

// Scan Return the data
func (m *MockJournal_MultipleRows) Scan(dest ...interface{}) error {
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

// MockJournal_SingleRow Mock single row returned for a Journal
type MockJournal_SingleRow struct {
	MockRowsEmpty
	RowNumber int
}

// Next Mock 1 row
func (m *MockJournal_SingleRow) Next() bool {
	m.RowNumber++
	if m.RowNumber < 2 {
		return true
	}
	return false
}

// Scan Return the data
func (m *MockJournal_SingleRow) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "slug"
		*dest[2].(*string) = "Title"
		*dest[3].(*string) = "2018-02-01"
		*dest[4].(*string) = "Content"
	}
	return nil
}

// MockResult Mock the result for a saved Journal
type MockResult struct{}

// LastInsertId Mock the last inserted ID
func (m *MockResult) LastInsertId() (int64, error) {
	return 1, nil
}

// RowsAffected Mock the rows affected
func (m *MockResult) RowsAffected() (int64, error) {
	return 0, nil
}

// MockRowsEmpty An empty row set
type MockRowsEmpty struct{}

// Close Mock close method
func (m *MockRowsEmpty) Close() error {
	return nil
}

// Columns Mock columns method
func (m *MockRowsEmpty) Columns() ([]string, error) {
	return []string{}, nil
}

// Next No rows
func (m *MockRowsEmpty) Next() bool {
	return false
}

// Scan No rows
func (m *MockRowsEmpty) Scan(dest ...interface{}) error {
	return nil
}

// MockSqlite Mock model.Sqlite allowing injected results, rows and errors
type MockSqlite struct {
	Closed           bool
	Connected        bool
	ErrorAtQuery     int
	ErrorMode        bool
	ExpectedArgument string
	Queries          int
	Result           sql.Result
	Rows             Rows
}

// Close Mark as closed
func (m *MockSqlite) Close() {
	m.Closed = true
}

// Connect Mark as connected
func (m *MockSqlite) Connect() error {
	m.Connected = true
	return nil
}

// Exec Test arguments and errors
func (m *MockSqlite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	m.Queries++
	if m.ErrorMode || m.ErrorAtQuery == m.Queries {
		return nil, errors.New("Simulating error")
	}
	if m.ExpectedArgument != "" && !m.inArgs(args) {
		return nil, errors.New("Expected " + m.ExpectedArgument + " in query")
	}
	return m.Result, nil
}

// Query Test arguments and errors
func (m *MockSqlite) Query(sql string, args ...interface{}) (Rows, error) {
	m.Queries++
	if m.ErrorMode || m.ErrorAtQuery == m.Queries {
		return nil, errors.New("Simulating error")
	}
	if m.ExpectedArgument != "" && !m.inArgs(args) {
		return nil, errors.New("Expected " + m.ExpectedArgument + " in query")
	}
	return m.Rows, nil
}

func (m *MockSqlite) inArgs(slice []interface{}) bool {
	for _, v := range slice {
		if v.(string) == m.ExpectedArgument {
			return true
		}
	}
	return false
}
