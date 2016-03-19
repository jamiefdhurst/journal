package main

import (
	"database/sql"
	"flag"
	"fmt"
	"journal/controller"
	"journal/lib"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	const version = "0.1"

	// Command line flags
	var (
		mode = flag.String("mode", "run", "Run or create database file")
		port = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Set CWD
	os.Chdir(os.Getenv("GOPATH"))

	// Load database
	newdb, err := sql.Open("sqlite3", "./data/journal.db")
	db = newdb
	lib.CheckErr(err)
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	if *mode == "create" {

		_, err := db.Exec("CREATE TABLE `journal` (" +
			"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"`slug` VARCHAR(255) NOT NULL, " +
			"`title` VARCHAR(255) NOT NULL, " +
			"`date` DATE NOT NULL, " +
			"`content` TEXT NOT NULL" +
			")")
		lib.CheckErr(err)
		db.Close()
		log.Println("Database created")

	} else {

		m := &lib.Router{}
		m.SetDb(db)
		m.SetErr(&controller.Error{})
		m.Add("GET", "/", false, &controller.Index{})
		m.Add("GET", "/new", false, &controller.New{})
		m.Add("POST", "/new", false, &controller.New{})
		m.Add("GET", "\\/([\\w\\-]+)", true, &controller.View{})

		log.Printf("Listening on port %s\n", *port)
		log.Fatal(http.ListenAndServe(":"+*port, m))

		db.Close()

	}
}
