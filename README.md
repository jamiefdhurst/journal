# Journal

Written as a first attempt of a web-based project in Go, this is a simple web-
server driven blog with the ability to post new entries.

It makes use of a SQLite database to store the journal entries.

## Installation

* Clone the repository to `$GOPATH/src/journal`.
* Make sure the `$GOPATH/data` directory exists, and initialise a `journal.db`
SQLite file there.
* Run `go build journal` to create the executable.

## Options

`-mode=create` - Use to create the database and initialise the table.
`-port=80` - Use to set the port, defaults to 3000.
