package model

import (
	"bytes"
	"testing"

	"github.com/jamiefdhurst/journal/test/mocks/adapter"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestGiphys_ConvertIDsToIframes(t *testing.T) {
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch"
	gs := Giphys{Db: &database.MockSqlite{}}
	newString := gs.ConvertIDsToIframes(testString)
	if newString != "Hello\n<iframe src=\"https://giphy.com/embed/1234567\"></iframe>\n:gif:testsearch" {
		t.Errorf("Expected iframe substitution did not occur")
	}
}

func TestGiphys_CreateTable(t *testing.T) {
	db := &database.MockSqlite{}
	gs := Giphys{Db: db}
	gs.CreateTable()
	if db.Queries != 1 {
		t.Errorf("Expected 1 query to have been run")
	}
}

func TestGiphys_ExtractContentsAndSearchAPI(t *testing.T) {

	// Test without error
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch\n"
	client := adapter.MockGiphyAdapter{}
	gs := Giphys{Client: client, Db: &database.MockSqlite{}}
	newString := gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n:gif:id:9991234\n" {
		t.Errorf("Expected search to have been converted")
	}

	// Test with error
	client.ErrorMode = true
	gs = Giphys{Client: client, Db: &database.MockSqlite{}}
	newString = gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n\n" {
		t.Errorf("Expected search to have been converted and error to have been handled")
	}
}

func TestGiphys_GetAPIKey(t *testing.T) {

	// Test error
	db := &database.MockSqlite{}
	db.ErrorMode = true
	gs := Giphys{Db: db}
	apiKey := gs.GetAPIKey()
	if apiKey != "" {
		t.Errorf("Expected blank string if error returned")
	}

	db.ErrorMode = false
	db.Rows = &database.MockRowsEmpty{}
	apiKey = gs.GetAPIKey()
	if apiKey != "" {
		t.Errorf("Expected blank string if empty query result")
	}

	// Test successful return
	db.Rows = &database.MockGiphy_SingleRow{}
	apiKey = gs.GetAPIKey()
	if apiKey != "API123456" {
		t.Errorf("Expected correct API key but received blank")
	}
}

func TestGiphys_InputNewAPIKey(t *testing.T) {
	testKey := bytes.NewBufferString("API123456\n")
	db := &database.MockSqlite{}
	gs := Giphys{Db: db}
	db.ExpectedArgument = "API123456"
	err := gs.InputNewAPIKey(testKey)
	if err != nil {
		t.Errorf("Expected no error from inputting API key, got %s", err)
	}
	if db.Queries != 1 {
		t.Errorf("Expected 1 query to have run")
	}
}
