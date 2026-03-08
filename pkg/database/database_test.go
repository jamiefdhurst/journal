package database

import (
    "os"
    "path/filepath"
    "testing"
)

func TestSqliteConnect_NewFile(t *testing.T) {
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "new.db")

    sqlite := &Sqlite{}
    err := sqlite.Connect(dbPath)
    if err != nil {
        t.Errorf("Expected successful connect to new file, got: %s", err)
    }
    sqlite.Close()

    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        t.Error("Expected database file to have been created")
    }
}

func TestSqliteConnect_Error(t *testing.T) {
    sqlite := &Sqlite{}
    err := sqlite.Connect("/nonexistent/directory/test.db")
    if err == nil {
        t.Error("Expected error when connecting to uncreateable path")
    }
}

func TestSqliteClose(t *testing.T) {
    sqlite := &Sqlite{}
    _ = sqlite.Connect("../../test/data/test.db")
    sqlite.Close()
}

func TestSqliteConnect(t *testing.T) {
    sqlite := &Sqlite{}
    err := sqlite.Connect("../../test/data/test.db")
    if err != nil {
        t.Errorf("Expected database to have been connected and no error to have been returned, got %s", err)
    }
}

func TestSqliteExec(t *testing.T) {
    sqlite := &Sqlite{}
    _ = sqlite.Connect("../../test/data/test.db")
    result, err := sqlite.Exec("SELECT 1")
    rows, _ := result.RowsAffected()
    if err != nil || rows > 0 {
        t.Errorf("Expected query to have been executed and no rows to have been affected")
    }
}

func TestSqliteQuery(t *testing.T) {
    sqlite := &Sqlite{}
    _ = sqlite.Connect("../../test/data/test.db")
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
