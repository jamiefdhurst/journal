package model

import (
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"

	"github.com/jamiefdhurst/journal/test/mocks/adapter"
)

func TestGiphys_ConvertIDsToIframes(t *testing.T) {
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch"
	gs := Giphys{}
	newString := gs.ConvertIDsToIframes(testString)
	if newString != "Hello\n<iframe src=\"https://giphy.com/embed/1234567\"></iframe>\n:gif:testsearch" {
		t.Errorf("Expected iframe substitution did not occur")
	}
}

func TestGiphys_ExtractContentsAndSearchAPI(t *testing.T) {

	// Test without error
	testString := "Hello\n:gif:id:1234567\n:gif:testsearch\n"
	client := &adapter.MockGiphyAdapter{}
	container := &app.Container{Giphy: client}
	gs := Giphys{Container: container}
	newString := gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n:gif:id:9991234\n" {
		t.Errorf("Expected search to have been converted")
	}

	// Test with error
	client.ErrorMode = true
	gs = Giphys{Container: container}
	newString = gs.ExtractContentsAndSearchAPI(testString)
	if newString != "Hello\n:gif:id:1234567\n\n" {
		t.Errorf("Expected search to have been converted and error to have been handled")
	}
}
