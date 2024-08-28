package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/jamiefdhurst/journal/pkg/adapter/giphy"
	"github.com/jamiefdhurst/journal/pkg/adapter/json"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/internal/app/router"
	"github.com/jamiefdhurst/journal/pkg/database"
)

var container *app.Container = &app.Container{}

func config() app.Configuration {
	// Define default configuration
	configuration := app.DefaultConfiguration()
	app.ApplyEnvConfiguration(&configuration)

	if !configuration.EnableCreate {
		log.Println("Article creating is disabled...")
	}
	if !configuration.EnableEdit {
		log.Println("Article editing is disabled...")
	}

	return configuration
}

func loadDatabase() func() {
	container.Db = &database.Sqlite{}
	log.Printf("Loading DB from %s...\n", container.Configuration.DatabasePath)
	if err := container.Db.Connect(container.Configuration.DatabasePath); err != nil {
		log.Printf("Database error - please verify that the %s path is available and writeable.\nError: %s\n", container.Configuration.DatabasePath, err)
		os.Exit(1)
	}

	js := model.Journals{Container: container}
	if err := js.CreateTable(); err != nil {
		log.Panicln(err)
	}

	return func() {
		container.Db.Close()
	}
}

func loadGiphy() {
	giphyAPIKey := os.Getenv("J_GIPHY_API_KEY")
	if giphyAPIKey != "" {
		log.Println("Enabling GIPHY client...")
		container.Giphy = &giphy.Client{APIKey: giphyAPIKey, Client: &json.Client{}}
	}
}

func main() {
	const version = "0.9.4"

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s\n-------------------\n\n", version)

	configuration := config()

	// Create/define container
	container.Configuration = configuration
	container.Version = version

	closeFunc := loadDatabase()
	defer closeFunc()
	loadGiphy()

	router := router.NewRouter(container)

	var err error
	if lambdaRuntimeApi, _ := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); lambdaRuntimeApi != "" {
		log.Printf("Ready for Lambda payload...\n")
		algnhsa.ListenAndServe(router, nil)
	} else {
		server := &http.Server{Addr: ":" + configuration.Port, Handler: router}
		log.Printf("Ready and listening on port %s...\n", configuration.Port)
		err = router.StartAndServe(server)
	}

	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
