# Journal

![License](https://img.shields.io/github/license/jamiefdhurst/journal.svg)
[![Build Status](https://ci.jamiehurst.co.uk/buildStatus/icon?job=github%2Fjournal%2Fmaster)](https://ci.jamiehurst.co.uk/job/github/job/journal/job/master/)
[![Latest Version](https://img.shields.io/github/release/jamiefdhurst/journal.svg)](https://github.com/jamiefdhurst/journal/releases)

A simple web-based journal written in Go. You can post, edit and view entries,
witht he addition of an API.

It makes use of a SQLite database to store the journal entries.

[API Documentation](api/README.md)

## Installation and Setup (local method)

1. Clone the repository to `$GOPATH/src/github.com/jamiefdhurst/journal`.
2. Make sure the `$GOPATH/data` directory exists.
3. Change directory to `$GOPATH/src/github.com/jamiefdhurst/journal`.
4. Run `go get` to install dependencies
5. Run `go build journal` to create the executable.
6. Run `./journal` to load the application on port 3000. You should now be able
    to fully access it at [http://localhost:3000](http://localhost:3000)

## Installation and Setup (Docker method)

_Please note: you will need Docker installed on your local machine._

1. Clone the repository to your chosen folder.
2. Build the container with `docker build -t journal:latest .`
3. Run the following to load the application and serve it on port 3000. You
    should now be able to fully access it at [http://localhost:3000](http://localhost:3000)

    ```bash
    docker run --rm -v ./data:/go/data -p 3000:3000 -it journal:latest
    ```

## Environment Variables

* `J_ARTICLES_PER_PAGE` - Articles to display per page, default `20`
* `J_DB_PATH` - Path to SQLite DB - default is `$GOPATH/data/journal.db`
* `J_GIPHY_API_KEY` - Set to a GIPHY API key to use, or ignore to disable GIPHY
* `J_PORT` - Port to expose over HTTP, default is `3000`
* `J_TITLE` - Set the title of the Journal

To use the API key within your Docker setup, include it as follows:

```bash
docker run --rm -e J_GIPHY_API_KEY=... -v ./data:/go/data -p 3000:3000 -it journal:latest
```

## Layout

The project layout follows the standard set out in the following document:
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

* `/api` - API documentation
* `/cmd/journal` - Main Journal executable folder
* `/internal/app/controller` - Controllers for the main application
* `/internal/app/model` - Models for the main application
* `/internal/app/router` - Implementation of router for given app
* `/pkg/adapter` - Adapters for connecting to external services
* `/pkg/controller` - Controller logic
* `/pkg/database` - Database connection logic
* `/pkg/router` - Router for handling services
* `/test` - API tests
* `/test/data` - Test data
* `/test/mocks` - Mock files for testing
* `/web/app` - CSS/JS source files
* `/web/static` - Compiled static public assets
* `/web/templates` - View templates

## Development

### Back-end

...

### Front-end

...