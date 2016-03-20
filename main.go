package main

import (
	"flag"
	"fmt"
	"journal/lib"
	"os"
)

var server lib.Server

func main() {
	const version = "0.1"

	// Command line flags
	var (
		mode = flag.String("mode", "run", "Run or create database file")
		port = flag.String("port", "3000", "Port to run web server on")
	)
	flag.Parse()

	// Set CWD
	os.Chdir(os.Getenv("GOPATH"))
	fmt.Printf("Journal v%s...\n-------------------\n\n", version)

	// Create the server
	server = lib.NewServer()
	server.Run(*mode, *port)
}
