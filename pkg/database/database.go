package database

// Database Define a common interface for all database drivers
type Database interface {
	Close()
	Connect(dbFile string) error
}

// Dynamodb denotes a DynamoDB database type
const Dynamodb string = "dynamodb"

// Sqlite denotes a SQLite database type
const Sqlite string = "sqlite"
