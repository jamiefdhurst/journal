package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jamiefdhurst/journal/controller/web"
	"github.com/jamiefdhurst/journal/lib"
	"github.com/jamiefdhurst/journal/model"
)

func main() {
	const version = "0.1"

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
	db := &model.Sqlite{}
	if err := db.Connect(); err != nil {
		log.Println("Database error - please verify that the $GOPATH/data folder is available.")
		os.Exit(1)
	}

	// Handle mode
	var err error
	if *mode == "createdb" {

		err = model.CreateTables(db)
		log.Println("Database created")

	} else if *mode == "giphy" {

		gs := model.Giphys{Db: db}
		err = gs.InputNewAPIKey(os.Stdin)
		log.Println("API key saved")

	} else {

		router := &lib.Router{Db: db, ErrorController: &web.Error{}}
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
