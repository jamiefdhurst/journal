package model

import (
	"database/sql"
	"errors"
	"os"
	"testing"
)

type FakeSqlite struct {
	Closed           bool
	Connected        bool
	ErrorAtQuery     int
	ErrorMode        bool
	ExpectedArgument string
	Queries          int
	Result           sql.Result
	Rows             Rows
}

func (f *FakeSqlite) Close() {
	f.Closed = true
}

func (f *FakeSqlite) Connect() error {
	f.Connected = true
	return nil
}

func (f *FakeSqlite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	f.Queries++
	if f.ErrorMode || f.ErrorAtQuery == f.Queries {
		return nil, errors.New("Simulating error")
	}
	if f.ExpectedArgument != "" && !f.inArgs(args) {
		return nil, errors.New("Expected " + f.ExpectedArgument + " in query")
	}
	return f.Result, nil
}

func (f *FakeSqlite) Query(sql string, args ...interface{}) (Rows, error) {
	f.Queries++
	if f.ErrorMode || f.ErrorAtQuery == f.Queries {
		return nil, errors.New("Simulating error")
	}
	if f.ExpectedArgument != "" && !f.inArgs(args) {
		return nil, errors.New("Expected " + f.ExpectedArgument + " in query")
	}
	return f.Rows, nil
}

func (f *FakeSqlite) inArgs(slice []interface{}) bool {
	for _, v := range slice {
		if v.(string) == f.ExpectedArgument {
			return true
		}
	}
	return false
}

func TestCreateTables(t *testing.T) {
	// Test without errors
	database := &FakeSqlite{}
	err := CreateTables(database)
	if err != nil {
		t.Errorf("Expected no error from creating tables")
	}
	if database.Queries != 2 {
		t.Errorf("Expected at least 2 queries from creating tables")
	}

	// Test with errors
	database.ErrorMode = true
	err = CreateTables(database)
	if err == nil {
		t.Errorf("Expected error from creating tables")
	}

	// Simulate error in 2nd query only
	database.ErrorMode = false
	database.ErrorAtQuery = 5
	err = CreateTables(database)
	if err == nil {
		t.Errorf("Expected error from creating tables")
	}
}

func TestSqliteClose(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect()
	sqlite.Close()
}

func TestSqliteConnect(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	err := sqlite.Connect()
	if err != nil {
		t.Errorf("Expected database to have been connected and no error to have been returned")
	}
}

func TestSqliteExec(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect()
	result, err := sqlite.Exec("SELECT 1")
	rows, _ := result.RowsAffected()
	if err != nil || rows > 0 {
		t.Errorf("Expected query to have been executed and no rows to have been affected")
	}
}

func TestSqliteQuery(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect()
	rows, err := sqlite.Query("SELECT 1 AS example")
	if err != nil {
		t.Errorf("Expected query to have been executed")
	}
	columns, _ := rows.Columns()
	if len(columns) != 1 || columns[0] != "example" {
		t.Errorf("Expected column of 'example' to have been returned")
	}
	var test int
	for rows.Next() {
		rows.Scan(&test)
		if test != 1 {
			t.Errorf("Expected row with value of '1' to have been returned")
		}
	}
	rows.Close()
}
