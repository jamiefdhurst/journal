package dynamodb

import (
	"github.com/jamiefdhurst/journal/pkg/database"
)

// Dynamodb Handle a DynamoDB connection
type Dynamodb struct {
	database.Database
}

// Close Close open database
func (d *Dynamodb) Close() {

}

// Connect Connect/open the database
func (d *Dynamodb) Connect(dbFile string) error {
	return nil
}
