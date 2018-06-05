package json

import (
	"encoding/json"
	"net/http"
)

// Adapter Common interface for a JSON client
type Adapter interface {
	Get(url string, destination interface{}) error
}

// Client for interacting with JSON over HTTP
type Client struct{}

// Get Perform a GET request to retrieve JSON
func (j Client) Get(url string, destination interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "Journal")
	client := &http.Client{}
	rs, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rs.Body.Close()
	json.NewDecoder(rs.Body).Decode(&destination)

	return nil
}
