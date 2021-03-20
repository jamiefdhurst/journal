package app

import (
	"database/sql"
	"os"

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
	Configuration Configuration
	Db            Database
	Giphy         GiphyAdapter
	Version       string
}

// Configuration can be modified through environment variables
type Configuration struct {
	ArticlesPerPage int
	DatabasePath    string
	EnableCreate    bool
	EnableDelete    bool
	EnableEdit      bool
	Port            string
	Title           string
}

// DefaultConfiguration returns the default settings for the app
func DefaultConfiguration() Configuration {
	return Configuration{
		ArticlesPerPage: 20,
		DatabasePath:    os.Getenv("GOPATH") + "/data/journal.db",
		EnableCreate:    true,
		EnableDelete:    true,
		EnableEdit:      true,
		Port:            "3000",
		Title:           "Jamie's Journal",
	}
}

// ApplyEnvConfiguration applys the env variables on top of existing config
func ApplyEnvConfiguration(config *Configuration) {
	port := os.Getenv("J_PORT")
	if port != "" {
		config.Port = port
	}
}
