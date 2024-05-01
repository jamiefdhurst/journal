package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	containers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type row struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

func setupContainer(t *testing.T) (string, func(t *testing.T)) {
	ctx := context.Background()
	req := containers.ContainerRequest{
		Image:        "amazon/dynamodb-local:latest",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	container, err := containers.GenericContainer(ctx, containers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start DynamoDB: %s", err)
	}
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("Could not get DynamoDB endpoint: %s", err)
	}

	return endpoint, func(t *testing.T) {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not stop DynamoDB: %s", err)
		}
	}
}

func connect(d *Dynamodb, e string, t *testing.T) {
	if err := d.Connect("journal|http://" + e); err != nil {
		t.Errorf("Unable to connect to DynamoDB: %s", err)
	}
	if err := d.CreateTable("journal"); err != nil {
		t.Errorf("Expected to be able to create DynamoDB table, but received: %s", err)
	}
}

func fill(d *Dynamodb, t *testing.T) {
	if err := d.PutItem(map[string]types.AttributeValue{
		"id":    &types.AttributeValueMemberN{Value: "1"},
		"slug":  &types.AttributeValueMemberS{Value: "foobar"},
		"title": &types.AttributeValueMemberS{Value: "Foo Bar"},
	}); err != nil {
		t.Errorf("Expected to be able to put item into DynamoDB, but received: %s", err)
	}
	if err := d.PutItem(map[string]types.AttributeValue{
		"id":    &types.AttributeValueMemberN{Value: "2"},
		"slug":  &types.AttributeValueMemberS{Value: "foobarbaz"},
		"title": &types.AttributeValueMemberS{Value: "Foo Bar Baz"},
	}); err != nil {
		t.Errorf("Expected to be able to put item into DynamoDB, but received: %s", err)
	}
	if err := d.PutItem(map[string]types.AttributeValue{
		"id":    &types.AttributeValueMemberN{Value: "3"},
		"slug":  &types.AttributeValueMemberS{Value: "something-else"},
		"title": &types.AttributeValueMemberS{Value: "Something Else"},
	}); err != nil {
		t.Errorf("Expected to be able to put item into DynamoDB, but received: %s", err)
	}
}

func TestConnect(t *testing.T) {
	db := &Dynamodb{}
	err := db.Connect("table")
	if err != nil {
		t.Errorf("Expected database to have been connected and no error to have been returned but received %s", err)
	}

	err = db.Connect("table|https://dynamodb.us-east-1.amazonaws.com")
	if err != nil {
		t.Errorf("Expected database to have been connected and no error to have been returned but received %s", err)
	}
}

func TestPutItem(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	db := &Dynamodb{}
	connect(db, ep, t)

	if err := db.PutItem(map[string]types.AttributeValue{
		"id":    &types.AttributeValueMemberN{Value: "1234"},
		"slug":  &types.AttributeValueMemberS{Value: "foo"},
		"title": &types.AttributeValueMemberS{Value: "bar"},
	}); err != nil {
		t.Errorf("Expected to be able to put item into DynamoDB, but received: %s", err)
	}
}

func TestQuery(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	db := &Dynamodb{}
	connect(db, ep, t)
	fill(db, t)

	// Error query (wrong query expression)
	expr, _ := expression.NewBuilder().WithKeyCondition(
		expression.Key("id").Equal(expression.Value("1"))).Build()
	row := row{}
	err := db.Query(expr, false, &row)
	if err == nil {
		t.Error("Expected error but received nothing")
	}

	// Empty query
	expr, _ = expression.NewBuilder().WithKeyCondition(
		expression.Key("id").Equal(expression.Value(64))).Build()
	err = db.Query(expr, false, &row)
	if err != nil {
		t.Errorf("Expected no error but received: %s", err)
	}
	if row.ID != 0 {
		t.Error("Expected row to be empty but received information")
	}

	// Successful query
	expr, _ = expression.NewBuilder().WithKeyCondition(
		expression.Key("id").Equal(expression.Value(1))).Build()
	err = db.Query(expr, false, &row)
	if err != nil {
		t.Errorf("Expected row to be returned but received error: %s", err)
	}
	if row.ID != 1 {
		t.Errorf("Expected row to be returned correctly but received ID: %d", row.ID)
	}
	if row.Slug != "foobar" {
		t.Errorf("Expected row to be returned correctly but received Slug: %s", row.Slug)
	}
	if row.Title != "Foo Bar" {
		t.Errorf("Expected row to be returned correctly but received Title: %s", row.Title)
	}
}

func TestScan(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	db := &Dynamodb{}
	connect(db, ep, t)
	fill(db, t)

	// Error scan (wrong unmarshall)
	expr, _ := expression.NewBuilder().Build()
	err := db.Scan(expr, &row{})
	if err == nil {
		t.Error("Expected error but received nothing")
	}

	// Correct scan
	rows := []row{}
	err = db.Scan(expr, &rows)
	if err != nil {
		t.Errorf("Expected rows to be returned but received error: %s", err)
	}
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows to be returned but received: %d", len(rows))
	}
	// if rows[0].ID != 1 {
	// 	t.Errorf("Expected row to be returned correctly but received ID: %d", rows[0].ID)
	// }
	// if rows[0].Slug != "foobar" {
	// 	t.Errorf("Expected row to be returned correctly but received Slug: %s", rows[0].Slug)
	// }
	// if rows[0].Title != "Foo Bar" {
	// 	t.Errorf("Expected row to be returned correctly but received Title: %s", rows[0].Title)
	// }
}

func TestScanCount(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	db := &Dynamodb{}
	connect(db, ep, t)
	fill(db, t)

	expr, _ := expression.NewBuilder().Build()
	rows, err := db.ScanCount(expr)
	if err != nil {
		t.Errorf("Expected rows to be returned but received error: %s", err)
	}
	if rows != 3 {
		t.Errorf("Expected 3 rows to be returned but received: %d", rows)
	}
}

func TestScanLimit(t *testing.T) {
	ep, tearDown := setupContainer(t)
	defer tearDown(t)

	db := &Dynamodb{}
	connect(db, ep, t)
	fill(db, t)

	expr, _ := expression.NewBuilder().WithFilter(
		expression.And(
			expression.AttributeNotExists(expression.Name("something-random")),
			expression.Equal(expression.Name("slug"), expression.Value("foobar")),
		)).Build()
	rows := []row{}
	err := db.ScanLimit(expr, 0, 1, &rows)
	if err != nil {
		t.Errorf("Expected rows to be returned but received error: %s", err)
	}
	if len(rows) != 1 {
		t.Errorf("Expected 1 row to be returned but received: %d", len(rows))
	}
	if rows[0].ID != 1 {
		t.Errorf("Expected row to be returned correctly but received ID: %d", rows[0].ID)
	}
	if rows[0].Slug != "foobar" {
		t.Errorf("Expected row to be returned correctly but received Slug: %s", rows[0].Slug)
	}
	if rows[0].Title != "Foo Bar" {
		t.Errorf("Expected row to be returned correctly but received Title: %s", rows[0].Title)
	}
}
