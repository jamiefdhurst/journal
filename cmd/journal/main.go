package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "net/http"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/internal/app/model"
    "github.com/jamiefdhurst/journal/internal/app/router"
    "github.com/jamiefdhurst/journal/pkg/database"
    "github.com/jamiefdhurst/journal/pkg/markdown"
)

var container *app.Container = &app.Container{}

func config() app.Configuration {
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

func bootstrap(c *app.Container, db app.Database, mp app.MarkdownProcessor) (func(), error) {
    c.Db = db
    c.MarkdownProcessor = mp

    log.Printf("Loading DB from %s...\n", c.Configuration.DatabasePath)
    if err := c.Db.Connect(c.Configuration.DatabasePath); err != nil {
        return nil, fmt.Errorf("database connect: %w", err)
    }

    js := model.Journals{Container: c}
    if err := js.CreateTable(); err != nil {
        return nil, fmt.Errorf("journal table: %w", err)
    }
    ms := model.Migrations{Container: c}
    if err := ms.CreateTable(); err != nil {
        return nil, fmt.Errorf("migrations table: %w", err)
    }
    vs := model.Visits{Container: c}
    if err := vs.CreateTable(); err != nil {
        return nil, fmt.Errorf("visits table: %w", err)
    }

    if err := ms.MigrateHTMLToMarkdown(); err != nil {
        return nil, fmt.Errorf("html to markdown migration: %w", err)
    }
    if err := ms.MigrateRandomSlugs(); err != nil {
        return nil, fmt.Errorf("random slug migration: %w", err)
    }
    if err := ms.MigrateAddTimestamps(); err != nil {
        return nil, fmt.Errorf("add timestamps migration: %w", err)
    }

    return func() { c.Db.Close() }, nil
}

func main() {
    const version = "0.9.6"
    fmt.Printf("Journal v%s\n-------------------\n\n", version)

    configuration := config()
    container.Configuration = configuration
    container.Version = version

    closeFunc, err := bootstrap(container, &database.Sqlite{}, &markdown.Markdown{})
    if err != nil {
        log.Fatalf("Setup failed: %s\n", err)
    }
    defer closeFunc()

    rtr := router.NewRouter(container)

    var protocols http.Protocols
    protocols.SetHTTP1(true)
    protocols.SetHTTP2(true)
    protocols.SetUnencryptedHTTP2(true)
    server := &http.Server{
        Addr:    ":" + configuration.Port,
        Handler: rtr,
        Protocols: &protocols,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS13,
        },
    }
    log.Printf("Ready and listening on port %s...\n", configuration.Port)
    if configuration.SSLCertificate == "" {
        err = rtr.StartAndServe(server)
    } else {
        log.Printf("Certificate: %s\n", configuration.SSLCertificate)
        log.Printf("Certificate Key: %s\n", configuration.SSLKey)
        log.Println("Serving with SSL enabled...")
        err = rtr.StartAndServeTLS(server, configuration.SSLCertificate, configuration.SSLKey)
    }

    if err != nil {
        log.Fatal("Error reported: ", err)
    }
}
