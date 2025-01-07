# Journal

![License](https://img.shields.io/github/license/jamiefdhurst/journal.svg)
[![Build](https://github.com/jamiefdhurst/journal/actions/workflows/build.yml/badge.svg)](https://github.com/jamiefdhurst/journal/actions/workflows/build.yml)
[![Latest Version](https://img.shields.io/github/release/jamiefdhurst/journal.svg)](https://github.com/jamiefdhurst/journal/releases)

A simple web-based journal written in Go. You can post, edit and view entries,
with the addition of an API.

It makes use of a SQLite database to store the journal entries.

[API Documentation](api/README.md)

## Purpose

Journal serves as an easy-to-read and simple Golang program for new developers 
to try their hand at modifying, extending and playing with. It deliberately has 
only one dependency to ensure that the full end-to-end flow of the system can 
be understood through standard Golang libraries.

It's also a nice little Journal that you can use to keep your thoughts in, or 
as a basic blog platform.

## Installation and Setup (local method)

1. Clone the repository.
2. Make sure the `$GOPATH/data` directory exists.
3. Run `go get ./...` to install dependencies
4. Run `go build journal.go` to create the executable.
5. Run `./journal` to load the application on port 3000. You should now be able
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
* `J_CREATE` - Set to `0` to disable article creation
* `J_DB_PATH` - Path to SQLite DB - default is `$GOPATH/data/journal.db`
* `J_DESCRIPTION` - Set the HTML description of the Journal
* `J_EDIT` - Set to `0` to disable article modification
* `J_GA_CODE` - Google Analytics tag value, starts with `UA-`, or ignore to disable Google Analytics
* `J_GIPHY_API_KEY` - Set to a GIPHY API key to use, or ignore to disable GIPHY
* `J_PORT` - Port to expose over HTTP, default is `3000`
* `J_THEME` - Theme to use from within the _web/themes_ folder, defaults to `default`
* `J_TITLE` - Set the title of the Journal

To use the API key within your Docker setup, include it as follows:

```bash
docker run --rm -e J_GIPHY_API_KEY=... -v ./data:/go/data -p 3000:3000 -it journal:latest
```

## Layout

The project layout follows the standard set out in the following document:
[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

* `/api` - API documentation
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
* `/web/static` - Compiled static public assets
* `/web/templates` - View templates
* `/web/themes` - Front-end themes, a default theme is included

## Development

### Back-end

The back-end can be extended and modified following the folder structure above. 
Tests for each file live alongside and are designed to be easy to read and as 
functionally complete as possible.

The easiest way to develop incrementally is to use a local go installation and 
run your Journal as follows:

```bash
go run journal.go
```

Naturally, any changes to the logic or functionality will require a restart of 
the binary itself.

#### Dependencies

The application currently only has one dependency:

* [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)

This can be installed using the following commands from the journal folder:

```bash
go get -v ./...
go install -v ./...
```

#### Templates

The templates are in `html/template` format in _web/templates_ and are used 
within each of the controllers. These can be modified while the binary stays 
loaded, as they are loaded on the fly by the application as it runs and serves 
content.

### Front-end

The front-end source files are intended to be divided into themes within the
_web/themes_ folder. Each theme can include icons and a CSS stylesheet.

A simple, basic and minimalist "default" theme is included, but any other 
themes can be built and modified.

### Building/Testing

All pushed code is currently built using GitHub Actions to test PRs, build 
packages and create releases.

To test locally, simply use:

```bash
go test -v ./...
```

### Building for Lambda

The application is designed to run as a Lambda connected to an EFS for SQLite 
storage. This requires a different method of building to ensure it includes the 
appropriate libraries and is built for the correct architecture.

To build for Lambda, you will need the x86_64-unknown-linux-gnu cross compiler
(if you're on a Mac):

```bash
brew tap SergioBenitez/osxct
brew install x86_64-unknown-linux-gnu
```

To build, simply run:

```bash
make build
```

This will produce a Lambda output: `lambda.zip`, that you can upload onto an 
Amazon Linux 2023 (al2023) runtime Lambda, if you configure the appropriate 
environment variables.
