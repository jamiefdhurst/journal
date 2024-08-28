package json

import "testing"

type TestResponse struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const testURL = "https://jsonplaceholder.typicode.com/todos/1"

func TestGet(t *testing.T) {
	response := &TestResponse{}
	jsonAdapter := Client{}
	err := jsonAdapter.Get(testURL, response)
	if err != nil {
		t.Error("Expected no error from test API call")
	}
	if response.ID != 1 && response.Title != "delectus aut autem" {
		t.Error("Expected result from JSON decode was not achieved")
	}

	// Create error in request
	err = jsonAdapter.Get("://Not a URL", response)
	if err == nil {
		t.Error("Expected error with blank request was not achieved")
	}

	// Create error in response
	err = jsonAdapter.Get("https://not-a-url.com", response)
	if err == nil {
		t.Error("Expected error with invalid request was not achieved")
	}
}
