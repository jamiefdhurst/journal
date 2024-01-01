package sqlite

import (
	"os"
	"testing"
)

func TestSqliteClose(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect("test/data/test.db")
	sqlite.Close()
}

func TestSqliteConnect(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	err := sqlite.Connect("test/data/test.db")
	if err != nil {
		t.Errorf("Expected database to have been connected and no error to have been returned")
	}
}

func TestSqliteExec(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect("test/data/test.db")
	result, err := sqlite.Exec("SELECT 1")
	rows, _ := result.RowsAffected()
	if err != nil || rows > 0 {
		t.Errorf("Expected query to have been executed and no rows to have been affected")
	}
}

func TestSqliteQuery(t *testing.T) {
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	sqlite := &Sqlite{}
	_ = sqlite.Connect("test/data/test.db")
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
