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
	const version = "0.3.0"

	// Command line flags
	var (
		serverPort = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	// Create/define container
	container := &app.Container{
		ArticlesPerPage: 20,
		Title:           "Jamie's Journal",
		Version:         version,
	}

	// Open database
	db := &database.Sqlite{}
	if err := db.Connect(os.Getenv("GOPATH") + "/data/journal.db"); err != nil {
		log.Println("Database error - please verify that the $GOPATH/data folder is available.")
		os.Exit(1)
	}

	// Create Giphy adapter
	giphyAPIKey := os.Getenv("GIPHY_API_KEY")
	if giphyAPIKey != "" {
		container.Giphy = &giphy.Client{APIKey: giphyAPIKey, Client: &json.Client{}}
	}

	// Create table if required
	container.Db = db
	var err error
	js := model.Journals{Container: container}
	if err = js.CreateTable(); err != nil {
		log.Panicln(err)
	}

	router := router.NewRouter(container)
	server := &http.Server{Addr: ":" + *serverPort, Handler: router}

	log.Printf("Listening on port %s\n", *serverPort)
	err = router.StartAndServe(server)

	// Close cleanly
	db.Close()
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
