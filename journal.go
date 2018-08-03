package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jamiefdhurst/journal/pkg/adapter/giphy"
	"github.com/jamiefdhurst/journal/pkg/adapter/json"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/internal/app/router"
	"github.com/jamiefdhurst/journal/pkg/database"
)

func main() {
	const version = "0.2.1"

	// Command line flags
	var (
		mode       = flag.String("mode", "run", "Run or perform a maintenance action (e.g. createdb for creating the database)")
		serverPort = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	// Create/define container
	container := &app.Container{}

	// Open database
	db := &database.Sqlite{}
	if err := db.Connect(os.Getenv("GOPATH") + "/data/journal.db"); err != nil {
		log.Println("Database error - please verify that the $GOPATH/data folder is available.")
		os.Exit(1)
	}

	// Create Giphy adapter
	adapter := &giphy.Client{APIKey: os.Getenv("GIPHY_API_KEY"), Client: &json.Client{}}

	container.Db = db
	container.Giphy = adapter

	// Handle mode
	var err error
	if *mode == "createdb" {

		js := model.Journals{Container: container}
		if err := js.CreateTable(); err != nil {
			log.Panicln(err)
		}

		log.Println("Database created")

	} else {

		router := router.NewRouter(container)
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
