package model

import (
	"bytes"
	"errors"
	"testing"
)

type FakeClient struct {
	ErrorMode bool
}

func (f FakeClient) SearchForID(s string) (string, error) {
	if f.ErrorMode {
		return "", errors.New("Simulating error")
	}
	if s == "testsearch" {
		return "9991234", nil
	}

	return "0000000", nil
}

type MockEmptyRows struct{}

func (m *MockEmptyRows) Close() error {
	return nil
}

func (m *MockEmptyRows) Columns() ([]string, error) {
	return []string{}, nil
}

func (m *MockEmptyRows) Next() bool {
	return false
}

func (m *MockEmptyRows) Scan(dest ...interface{}) error {
	return nil
}

type MockAPIReturnedRow struct {
	MockEmptyRows
	RowNumber int
}

func (m *MockAPIReturnedRow) Next() bool {
	m.RowNumber++
	if m.RowNumber < 2 {
		return true
	}
	return false
}

func (m *MockAPIReturnedRow) Scan(dest ...interface{}) error {
	if m.RowNumber == 1 {
		*dest[0].(*int) = 1
		*dest[1].(*string) = "API123456"
	}
	return nil
}

func TestGiphys_ConvertIDsToIframes(t *testing.T) {
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch"
	gs := Giphys{Db: &FakeSqlite{}}
	newString := gs.ConvertIDsToIframes(testString)
	if newString != "Hello\n<iframe src=\"https://giphy.com/embed/1234567\"></iframe>\n:gif:testsearch" {
		t.Errorf("Expected iframe substitution did not occur")
	}
}

func TestGiphys_CreateTable(t *testing.T) {
	database := &FakeSqlite{}
	gs := Giphys{Db: database}
	gs.CreateTable()
	if database.Queries != 1 {
		t.Errorf("Expected 1 query to have been run")
	}
}

func TestGiphys_ExtractContentsAndSearchAPI(t *testing.T) {

	// Test without error
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch\n"
	client := FakeClient{}
	gs := Giphys{Client: client, Db: &FakeSqlite{}}
	newString := gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n:gif:id:9991234\n" {
		t.Errorf("Expected search to have been converted")
	}

	// Test with error
	client.ErrorMode = true
	gs = Giphys{Client: client, Db: &FakeSqlite{}}
	newString = gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n\n" {
		t.Errorf("Expected search to have been converted and error to have been handled")
	}
}

func TestGiphys_GetAPIKey(t *testing.T) {

	// Test error
	database := &FakeSqlite{}
	database.ErrorMode = true
	gs := Giphys{Db: database}
	apiKey := gs.GetAPIKey()
	if apiKey != "" {
		t.Errorf("Expected blank string if error returned")
	}

	database.ErrorMode = false
	database.Rows = &MockEmptyRows{}
	apiKey = gs.GetAPIKey()
	if apiKey != "" {
		t.Errorf("Expected blank string if empty query result")
	}

	// Test successful return
	database.Rows = &MockAPIReturnedRow{}
	apiKey = gs.GetAPIKey()
	if apiKey != "API123456" {
		t.Errorf("Expected correct API key but received blank")
	}
}

func TestGiphys_InputNewAPIKey(t *testing.T) {
	testKey := bytes.NewBufferString("API123456\n")
	database := &FakeSqlite{}
	gs := Giphys{Db: database}
	database.ExpectedArgument = "API123456"
	err := gs.InputNewAPIKey(testKey)
	if err != nil {
		t.Errorf("Expected no error from inputting API key, got %s", err)
	}
	if database.Queries != 1 {
		t.Errorf("Expected 1 query to have run")
	}
}
