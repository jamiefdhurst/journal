package dynamodb

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jamiefdhurst/journal/pkg/database"
)

type DynamodbLike interface {
	Close()
	Connect(string) error
	PutItem(map[string]types.AttributeValue) error
	Query(expression.Expression, bool, interface{}) error
	Scan(expression.Expression, interface{}) error
	ScanCount(expression.Expression) (int32, error)
	ScanLimit(expression.Expression, int, int, interface{}) error
}

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

	cfg, _ := config.LoadDefaultConfig(context.TODO())
	d.db = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &endpoint
	})

	return nil
}

func (d *Dynamodb) CreateTable(table string) error {
	_, err := d.db.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName:   &d.tableName,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
	})

	return err
}

// PutItem save or replace an existing entry
func (d *Dynamodb) PutItem(item map[string]types.AttributeValue) error {
	_, err := d.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})

	return err
}

// Query a table and return a single results if available, with expression
func (d *Dynamodb) Query(expr expression.Expression, sortForward bool, result interface{}) error {
	var err error
	var response *dynamodb.QueryOutput
	response, err = d.db.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 &d.tableName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ScanIndexForward:          &sortForward,
	})

	if err != nil {
		return err
	}

	if response.Count == 0 {
		return nil
	}

	return attributevalue.UnmarshalMap(response.Items[0], &result)
}

// Scan a table to return all information, with expression
func (d *Dynamodb) Scan(expr expression.Expression, result interface{}) error {
	var err error
	var response *dynamodb.ScanOutput
	scanPaginator := dynamodb.NewScanPaginator(d.db, &dynamodb.ScanInput{
		TableName:                 &d.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	var items []map[string]types.AttributeValue
	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			return err
		}
		items = append(items, response.Items...)
	}

	return attributevalue.UnmarshalListOfMaps(items, &result)
}

// ScanCount scan a table and return a count, with expression
func (d *Dynamodb) ScanCount(expr expression.Expression) (int32, error) {
	var err error
	var response *dynamodb.ScanOutput
	response, err = d.db.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 &d.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		Select:                    types.SelectCount,
	})

	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

// ScanLimit scan a table to return all information with a limited set of results, with expression
func (d *Dynamodb) ScanLimit(expr expression.Expression, offset int, limit int, result interface{}) error {
	var err error
	var response *dynamodb.ScanOutput
	scanPaginator := dynamodb.NewScanPaginator(d.db, &dynamodb.ScanInput{
		TableName:                 &d.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})
	var items []map[string]types.AttributeValue
	for scanPaginator.HasMorePages() {
		response, err = scanPaginator.NextPage(context.TODO())
		if err != nil {
			return err
		}
		items = append(items, response.Items...)
	}

	return attributevalue.UnmarshalListOfMaps(items[offset:offset+limit], &result)
}
