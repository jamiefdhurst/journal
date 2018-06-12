package app

import (
	"database/sql"

	"github.com/jamiefdhurst/journal/pkg/database/rows"
)

// Database Define same interface as database
type Database interface {
	Close()
	Connect(dbFile string) error
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Query(sql string, args ...interface{}) (rows.Rows, error)
}

// GiphyAdapter Interface for API
type GiphyAdapter interface {
	SearchForID(s string) (string, error)
}

// Container Define the main container for the application
type Container struct {
	Db    Database
	Giphy GiphyAdapter
}
