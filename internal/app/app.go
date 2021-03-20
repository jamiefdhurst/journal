package app

import (
	"database/sql"
	"os"
	"strconv"

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
		EnableEdit:      true,
		Port:            "3000",
		Title:           "Jamie's Journal",
	}
}

// ApplyEnvConfiguration applys the env variables on top of existing config
func ApplyEnvConfiguration(config *Configuration) {
	articles, _ := strconv.Atoi(os.Getenv("J_ARTICLES_PER_PAGE"))
	if articles > 0 {
		config.ArticlesPerPage = articles
	}
	database := os.Getenv("J_DB_PATH")
	if database != "" {
		config.DatabasePath = database
	}
	enableCreate := os.Getenv("J_CREATE")
	if enableCreate == "0" {
		config.EnableCreate = false
	}
	enableEdit := os.Getenv("J_EDIT")
	if enableEdit == "0" {
		config.EnableEdit = false
	}
	port := os.Getenv("J_PORT")
	if port != "" {
		config.Port = port
	}
	title := os.Getenv("J_TITLE")
	if title != "" {
		config.Title = title
	}
}
