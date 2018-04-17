# Journal

Written as a first attempt of a web-based project in Go, this is a simple web-
server driven blog with the ability to post new entries.

It makes use of a SQLite database to store the journal entries.

## Installation and Setup

1. Clone the repository to `$GOPATH/src/github.com/jamiefdhurst/journal`.
2. Make sure the `$GOPATH/data` directory exists.
3. Run `go build journal` to create the executable.
4. Run `./journal -mode=createdb` to create the database.
5. Run `./journal` to load the application on port 3000. You should now be able
    to fully access it at [](http://localhost:3000)

## Options

* `-mode=createdb` - Use to create the database within the data directory.
* `-mode=giphy` - Input a GIPHY API key to enable GIF usage within posts.
* `-mode=giphy` - Input a GIPHY API key to enable GIF usage within posts.
* `-port=3000` - Use to set the port to serve on for HTTP, defaults to 3000.
