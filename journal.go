package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

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
		log.Println("Post creating is disabled...")
	}
	if !configuration.EnableEdit {
		log.Println("Post editing is disabled...")
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

	// Create needed tables
	js := model.Journals{Container: container}
	if err := js.CreateTable(); err != nil {
		log.Printf("Error creating journal table: %s\n", err)
		log.Panicln(err)
	}
	ms := model.Migrations{Container: container}
	if err := ms.CreateTable(); err != nil {
		log.Printf("Error creating migrations table: %s\n", err)
		log.Panicln(err)
	}
	vs := model.Visits{Container: container}
	if err := vs.CreateTable(); err != nil {
		log.Printf("Error creating visits table: %s\n", err)
		log.Panicln(err)
	}

	// Run migrations
	if err := ms.MigrateHTMLToMarkdown(); err != nil {
		log.Printf("Error during HTML to Markdown migration: %s\n", err)
		log.Panicln(err)
	}
	if err := ms.MigrateRandomSlugs(); err != nil {
		log.Printf("Error during random slug migration: %s\n", err)
		log.Panicln(err)
	}
	if err := ms.MigrateAddTimestamps(); err != nil {
		log.Printf("Error during add timestamps migration: %s\n", err)
		log.Panicln(err)
	}

	return func() {
		container.Db.Close()
	}
}

func main() {
	const version = "0.9.6"
	fmt.Printf("Journal v%s\n-------------------\n\n", version)

	configuration := config()

	// Create/define container
	container.Configuration = configuration
	container.Version = version

	closeFunc := loadDatabase()
	defer closeFunc()

	router := router.NewRouter(container)

	var err error
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)
	server := &http.Server{
		Addr:      ":" + configuration.Port,
		Handler:   router,
		Protocols: &protocols,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
	log.Printf("Ready and listening on port %s...\n", configuration.Port)
	if configuration.SSLCertificate == "" {
		err = router.StartAndServe(server)
	} else {
		log.Printf("Certificate: %s\n", configuration.SSLCertificate)
		log.Printf("Certificate Key: %s\n", configuration.SSLKey)
		log.Println("Serving with SSL enabled...")
		err = router.StartAndServeTLS(server, configuration.SSLCertificate, configuration.SSLKey)
	}

	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
