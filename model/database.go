package model

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

var db *sql.DB

// Close Close open database
func Close() {
	db.Close()
}

func init() {
	db, _ = sql.Open("sqlite3", os.Getenv("GOPATH")+"/data/journal.db")
	err := db.Ping()
	if err != nil {
		log.Println("Database error - please verify that the $GOPATH/data folder is available.")
		os.Exit(1)
	}
}
