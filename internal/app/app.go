package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"os"
	"strconv"

	"github.com/jamiefdhurst/journal/pkg/database/rows"
	"github.com/jamiefdhurst/journal/pkg/env"
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
	DatabasePath        string
	Description         string
	EnableCreate        bool
	EnableEdit          bool
	ExcerptWords        int
	GoogleAnalyticsCode string
	Port                string
	PostsPerPage        int
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
		DatabasePath:        os.Getenv("GOPATH") + "/data/journal.db",
		Description:         "A fantastic journal containing some thoughts, ideas and reflections",
		EnableCreate:        true,
		EnableEdit:          true,
		ExcerptWords:        50,
		GoogleAnalyticsCode: "",
		Port:                "3000",
		PostsPerPage:        20,
		SSLCertificate:      "",
		SSLKey:              "",
		StaticPath:          "web/static",
		Theme:               "default",
		ThemePath:           "web/themes",
		Title:               "A Fantastic Journal",
		SessionKey:          "",
		SessionName:         "journal-session",
		CookieDomain:        "",
		CookieMaxAge:        2592000,
		CookieSecure:        false,
		CookieHTTPOnly:      true,
	}
}

// ApplyEnvConfiguration applies the env variables on top of existing config
// It first loads values from a .env file (if it exists), then applies any
// environment variables set in the system (which override .env values)
func ApplyEnvConfiguration(config *Configuration) {
	// Parse .env file (if it exists)
	dotenvVars, _ := env.Parse(".env")

	// Helper function to get env var, preferring system env over .env file
	getEnv := func(key string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return dotenvVars[key]
	}

	// J_ARTICLES_PER_PAGE is deprecated, but it's checked first
	articles, _ := strconv.Atoi(getEnv("J_ARTICLES_PER_PAGE"))
	if articles > 0 {
		config.PostsPerPage = articles
	}
	posts, _ := strconv.Atoi(getEnv("J_POSTS_PER_PAGE"))
	if posts > 0 {
		config.PostsPerPage = posts
	}
	database := getEnv("J_DB_PATH")
	if database != "" {
		config.DatabasePath = database
	}
	description := getEnv("J_DESCRIPTION")
	if description != "" {
		config.Description = description
	}
	enableCreate := getEnv("J_CREATE")
	if enableCreate == "0" {
		config.EnableCreate = false
	}
	enableEdit := getEnv("J_EDIT")
	if enableEdit == "0" {
		config.EnableEdit = false
	}
	excerptWords, _ := strconv.Atoi(getEnv("J_EXCERPT_WORDS"))
	if excerptWords > 0 {
		config.ExcerptWords = excerptWords
	}
	config.GoogleAnalyticsCode = getEnv("J_GA_CODE")
	port := getEnv("J_PORT")
	if port != "" {
		config.Port = port
	}

	config.SSLCertificate = getEnv("J_SSL_CERT")
	config.SSLKey = getEnv("J_SSL_KEY")
	staticPath := getEnv("J_STATIC_PATH")
	if staticPath != "" {
		config.StaticPath = staticPath
	}
	theme := getEnv("J_THEME")
	if theme != "" {
		config.Theme = theme
	}
	themePath := getEnv("J_THEME_PATH")
	if themePath != "" {
		config.ThemePath = themePath
	}
	title := getEnv("J_TITLE")
	if title != "" {
		config.Title = title
	}

	sessionKey := getEnv("J_SESSION_KEY")
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

	sessionName := getEnv("J_SESSION_NAME")
	if sessionName != "" {
		config.SessionName = sessionName
	}
	cookieDomain := getEnv("J_COOKIE_DOMAIN")
	if cookieDomain != "" {
		config.CookieDomain = cookieDomain
	}
	cookieMaxAge, _ := strconv.Atoi(getEnv("J_COOKIE_MAX_AGE"))
	if cookieMaxAge > 0 {
		config.CookieMaxAge = cookieMaxAge
	}
	cookieHTTPOnly := getEnv("J_COOKIE_HTTPONLY")
	if cookieHTTPOnly == "0" || cookieHTTPOnly == "false" {
		config.CookieHTTPOnly = false
	}
	if config.SSLCertificate != "" {
		config.CookieSecure = true
	}
}
