package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jamiefdhurst/journal/lib"
)

var app lib.App

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

	app = lib.App{}
	app.Run(*mode, *serverPort)
}
