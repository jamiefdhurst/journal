package giphy

import (
	"errors"

	"github.com/jamiefdhurst/journal/pkg/adapter/json"
)

// APIResponse Response holder for GIPHY API call
type APIResponse struct {
	Data APIResponseData `json:"data"`
}

// APIResponseData Data object within API response
type APIResponseData struct {
	ID string `json:"id"`
}

// Adapter Interface for API
type Adapter interface {
	SearchForID(s string) (string, error)
}

// Client Actual API client
type Client struct {
	Adapter
	APIKey string
	Client json.Adapter
}

// SearchForID Search the Giphy API for a given tag and return the resulting ID
func (c Client) SearchForID(s string) (string, error) {
	if c.APIKey == "" {
		return "", errors.New("No API key was found for GIPHY")
	}

	// Perform search
	url := "https://api.giphy.com/v1/gifs/random?api_key=" + c.APIKey + "&tag=" + s + "&rating=G"
	response := APIResponse{}
	err := c.Client.Get(url, &response)
	if err != nil {
		return "", err
	}

	if response.Data.ID != "" {
		return response.Data.ID, nil
	}

	return "", errors.New("No response was provided")
}
