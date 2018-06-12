package database

import (
	"database/sql"

	"github.com/jamiefdhurst/journal/pkg/database/rows"
	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

// Database Define a common interface for all database drivers
type Database interface {
	Close()
	Connect(dbFile string) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (rows.Rows, error)
}

// Sqlite Handle an Sqlite connection
type Sqlite struct {
	Database
	db *sql.DB
}

// Close Close open database
func (s *Sqlite) Close() {
	s.db.Close()
}

// Connect Connect/open the database
func (s *Sqlite) Connect(dbFile string) error {
	s.db, _ = sql.Open("sqlite3", dbFile)
	return s.db.Ping()
}

// Exec Execute a query on the database, returning a simple result
func (s *Sqlite) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return s.db.Exec(sql, args...)
}

// Query Query the database
func (s *Sqlite) Query(sql string, args ...interface{}) (rows.Rows, error) {
	return s.db.Query(sql, args...)
}
