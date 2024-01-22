package dynamodb

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jamiefdhurst/journal/pkg/database"
)

const defaultEndpoint = "https://dynamodb.eu-west-1.amazonaws.com"

// Dynamodb Handle a DynamoDB connection
type Dynamodb struct {
	database.Database
	db        *dynamodb.Client
	tableName string
}

// Close Close open database
func (d *Dynamodb) Close() {}

// Connect Connect/open the database
func (d *Dynamodb) Connect(dbFile string) error {
	// Split config with pipe if present for endpoint and table name
	endpoint := defaultEndpoint
	dbFileParts := strings.Split(dbFile, "|")
	d.tableName = dbFileParts[0]
	if len(dbFileParts) > 1 {
		endpoint = dbFileParts[1]
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	d.db = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})

	return nil
}
