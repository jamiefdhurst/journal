package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/internal/app/router"
	"github.com/jamiefdhurst/journal/pkg/database"
)

func main() {
	const version = "0.2"

	// Command line flags
	var (
		mode       = flag.String("mode", "run", "Run or perform a maintenance action (e.g. createdb for creating the database)")
		serverPort = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	// Open database
	db := &database.Sqlite{}
	if err := db.Connect(os.Getenv("GOPATH") + "/data/journal.db"); err != nil {
		log.Println("Database error - please verify that the $GOPATH/data folder is available.")
		os.Exit(1)
	}

	// Handle mode
	var err error
	if *mode == "createdb" {

		gs := model.Giphys{Db: db}
		if err := gs.CreateTable(); err != nil {
			log.Panicln(err)
		}
		js := model.Journals{Db: db}
		if err := js.CreateTable(); err != nil {
			log.Panicln(err)
		}

		log.Println("Database created")

	} else if *mode == "giphy" {

		gs := model.Giphys{Db: db}
		err = gs.InputNewAPIKey(os.Stdin)
		log.Println("API key saved")

	} else {

		router := router.NewRouter(db)
		server := &http.Server{Addr: ":" + *serverPort, Handler: router}

		log.Printf("Listening on port %s\n", *serverPort)
		err = router.StartAndServe(server)

	}

	// Close cleanly
	db.Close()
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
