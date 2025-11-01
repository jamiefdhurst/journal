package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
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

// MarkdownProcessor defines an interface for markdown processing
type MarkdownProcessor interface {
	ToHTML(input string) string
	FromHTML(input string) string
}

// Container Define the main container for the application
type Container struct {
	Configuration     Configuration
	Db                Database
	Version           string
	MarkdownProcessor MarkdownProcessor
}

// Configuration can be modified through environment variables
type Configuration struct {
	ArticlesPerPage     int
	DatabasePath        string
	Description         string
	EnableCreate        bool
	EnableEdit          bool
	ExcerptWords        int
	GoogleAnalyticsCode string
	Port                string
	SSLCertificate      string
	SSLKey              string
	StaticPath          string
	Theme               string
	ThemePath           string
	Title               string
	SessionKey          string
	SessionName         string
	CookieDomain        string
	CookieMaxAge        int
	CookieSecure        bool
	CookieHTTPOnly      bool
}

// DefaultConfiguration returns the default settings for the app
func DefaultConfiguration() Configuration {
	return Configuration{
		ArticlesPerPage:     20,
		DatabasePath:        os.Getenv("GOPATH") + "/data/journal.db",
		Description:         "A private journal containing Jamie's innermost thoughts",
		EnableCreate:        true,
		EnableEdit:          true,
		ExcerptWords:        50,
		GoogleAnalyticsCode: "",
		Port:                "3000",
		SSLCertificate:      "",
		SSLKey:              "",
		StaticPath:          "web/static",
		Theme:               "default",
		ThemePath:           "web/themes",
		Title:               "Jamie's Journal",
		SessionKey:          "",
		SessionName:         "journal-session",
		CookieDomain:        "",
		CookieMaxAge:        2592000,
		CookieSecure:        false,
		CookieHTTPOnly:      true,
	}
}

// ApplyEnvConfiguration applies the env variables on top of existing config
func ApplyEnvConfiguration(config *Configuration) {
	articles, _ := strconv.Atoi(os.Getenv("J_ARTICLES_PER_PAGE"))
	if articles > 0 {
		config.ArticlesPerPage = articles
	}
	database := os.Getenv("J_DB_PATH")
	if database != "" {
		config.DatabasePath = database
	}
	description := os.Getenv("J_DESCRIPTION")
	if description != "" {
		config.Description = description
	}
	enableCreate := os.Getenv("J_CREATE")
	if enableCreate == "0" {
		config.EnableCreate = false
	}
	enableEdit := os.Getenv("J_EDIT")
	if enableEdit == "0" {
		config.EnableEdit = false
	}
	excerptWords, _ := strconv.Atoi(os.Getenv("J_EXCERPT_WORDS"))
	if excerptWords > 0 {
		config.ExcerptWords = excerptWords
	}
	config.GoogleAnalyticsCode = os.Getenv("J_GA_CODE")
	port := os.Getenv("J_PORT")
	if port != "" {
		config.Port = port
	}
	config.SSLCertificate = os.Getenv("J_SSL_CERT")
	config.SSLKey = os.Getenv("J_SSL_KEY")

	sessionKey := os.Getenv("J_SESSION_KEY")
	if sessionKey != "" {
		if len(sessionKey) != 32 {
			log.Println("WARNING: J_SESSION_KEY must be exactly 32 bytes. Using auto-generated key instead.")
			sessionKey = ""
		}
	}
	if sessionKey == "" {
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err == nil {
			sessionKey = hex.EncodeToString(bytes)
			log.Println("WARNING: J_SESSION_KEY not set or invalid. Using auto-generated key. Sessions will not persist across restarts.")
		}
	}
	config.SessionKey = sessionKey

	sessionName := os.Getenv("J_SESSION_NAME")
	if sessionName != "" {
		config.SessionName = sessionName
	}

	cookieDomain := os.Getenv("J_COOKIE_DOMAIN")
	if cookieDomain != "" {
		config.CookieDomain = cookieDomain
	}

	cookieMaxAge, _ := strconv.Atoi(os.Getenv("J_COOKIE_MAX_AGE"))
	if cookieMaxAge > 0 {
		config.CookieMaxAge = cookieMaxAge
	}

	cookieHTTPOnly := os.Getenv("J_COOKIE_HTTPONLY")
	if cookieHTTPOnly == "0" || cookieHTTPOnly == "false" {
		config.CookieHTTPOnly = false
	}

	if config.SSLCertificate != "" {
		config.CookieSecure = true
	}

	staticPath := os.Getenv("J_STATIC_PATH")
	if staticPath != "" {
		config.StaticPath = staticPath
	}
	theme := os.Getenv("J_THEME")
	if theme != "" {
		config.Theme = theme
	}
	themePath := os.Getenv("J_THEME_PATH")
	if themePath != "" {
		config.ThemePath = themePath
	}
	title := os.Getenv("J_TITLE")
	if title != "" {
		config.Title = title
	}
}
