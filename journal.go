package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jamiefdhurst/journal/internal/app/controller/apiv1"
	"github.com/jamiefdhurst/journal/internal/app/controller/web"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/database"
	"github.com/jamiefdhurst/journal/pkg/router"
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
	if err := db.Connect(); err != nil {
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

		router := &router.Router{Db: db, ErrorController: &web.Error{}}
		server := &http.Server{Addr: ":" + *serverPort, Handler: router}

		router.Get("/new", &web.New{})
		router.Post("/new", &web.New{})
		router.Get("/api/v1/post", &apiv1.List{})
		router.Post("/api/v1/post", &apiv1.Create{})
		router.Get("/api/v1/post/[%s]", &apiv1.Single{})
		router.Put("/api/v1/post/[%s]", &apiv1.Update{})
		router.Get("/[%s]/edit", &web.Edit{})
		router.Post("/[%s]/edit", &web.Edit{})
		router.Get("/[%s]", &web.View{})
		router.Get("/", &web.Index{})

		log.Printf("Listening on port %s\n", *serverPort)
		err = router.StartAndServe(server)

	}

	// Close cleanly
	db.Close()
	if err != nil {
		log.Fatal("Error reported: ", err)
	}
}
