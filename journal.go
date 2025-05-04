package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/internal/app/router"
	"github.com/jamiefdhurst/journal/pkg/database"
	"github.com/jamiefdhurst/journal/pkg/markdown"
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
	
	// Set up the markdown processor
	container.MarkdownProcessor = &markdown.Markdown{}
	
	log.Printf("Loading DB from %s...\n", container.Configuration.DatabasePath)
	if err := container.Db.Connect(container.Configuration.DatabasePath); err != nil {
		log.Printf("Database error - please verify that the %s path is available and writeable.\nError: %s\n", container.Configuration.DatabasePath, err)
		os.Exit(1)
	}

	// Initialize journal table
	js := model.Journals{Container: container}
	if err := js.CreateTable(); err != nil {
		log.Panicln(err)
	}

	// Initialize and run migrations
	migrations := model.Migrations{Container: container}
	if err := migrations.CreateTable(); err != nil {
		log.Printf("Error creating migrations table: %s\n", err)
		log.Panicln(err)
	}

	// Run HTML to Markdown migration if needed
	if err := migrations.MigrateHTMLToMarkdown(); err != nil {
		log.Printf("Error during HTML to Markdown migration: %s\n", err)
		log.Panicln(err)
	}
	
	// Run random slug migration if needed
	if err := migrations.MigrateRandomSlugs(); err != nil {
		log.Printf("Error during random slug migration: %s\n", err)
		log.Panicln(err)
	}

	return func() {
		container.Db.Close()
	}
}

func main() {
	const version = "0.9.6"

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")
	fmt.Printf("Journal v%s\n-------------------\n\n", version)

	configuration := config()

	// Create/define container
	container.Configuration = configuration
	container.Version = version

	closeFunc := loadDatabase()
	defer closeFunc()

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
