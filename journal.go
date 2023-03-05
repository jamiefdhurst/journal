package main

import (
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
	const version = "0.3.0.3"

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s\n-------------------\n\n", version)

	// Define default configuration
	configuration := app.DefaultConfiguration()
	app.ApplyEnvConfiguration(&configuration)

	// Create/define container
	container := &app.Container{
		Configuration: configuration,
		Version:       version,
	}

	// Open database
	db := &database.Sqlite{}
	log.Printf("Loading DB from %s...\n", configuration.DatabasePath)
	if err := db.Connect(configuration.DatabasePath); err != nil {
		log.Printf("Database error - please verify that the %s path is available and writable.\n", configuration.DatabasePath)
		os.Exit(1)
	}

	// Create Giphy adapter
	giphyAPIKey := os.Getenv("J_GIPHY_API_KEY")
	if giphyAPIKey != "" {
		log.Println("Enabling GIPHY client...")
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
	server := &http.Server{Addr: ":" + configuration.Port, Handler: router}

	if !configuration.EnableCreate {
		log.Println("Article creating is disabled...")
	}
	if !configuration.EnableEdit {
		log.Println("Article editing is disabled...")
	}

	log.Printf("Ready and listening on port %s...\n", configuration.Port)
	err = router.StartAndServe(server)

	// Close cleanly
	db.Close()
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
