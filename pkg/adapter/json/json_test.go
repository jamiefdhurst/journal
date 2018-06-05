package json

import "testing"

type TestResponse struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   int    `json:"body"`
}

const testURL = "https://jsonplaceholder.typicode.com/posts/1"

func TestGet(t *testing.T) {
	response := &TestResponse{}
	jsonAdapter := Client{}
	err := jsonAdapter.Get(testURL, response)
	if err != nil {
		t.Error("Expected no error from test API call")
	}
	if response.ID != 1 && response.Title != "sunt aut facere repellat provident occaecati excepturi optio reprehenderit" {
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
