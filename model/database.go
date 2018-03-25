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
	var err error
	db, err = sql.Open("sqlite3", "./data/journal.db")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
