package adapter

import (
	"errors"
)

// GiphyAPIResponse Response holder for GIPHY API call
type GiphyAPIResponse struct {
	Data GiphyAPIResponseData `json:"data"`
}

// GiphyAPIResponseData Data object within API response
type GiphyAPIResponseData struct {
	ID string `json:"id"`
}

// GiphyAdapter Interface for API
type GiphyAdapter interface {
	SearchForID(s string) (string, error)
}

// GiphyClient Actual API client
type GiphyClient struct {
	GiphyAdapter
	APIKey string
	Client JSONClient
}

// GiphySearchAPI Search the Giphy API for a given tag and return the resulting ID
func (c GiphyClient) GiphySearchAPI(s string) (string, error) {
	if c.APIKey == "" {
		return "", errors.New("No API key was found for GIPHY")
	}

	// Perform search
	url := "https://api.giphy.com/v1/gifs/random?api_key=" + c.APIKey + "&tag=" + s + "&rating=G"
	response := GiphyAPIResponse{}
	err := c.Client.Get(url, &response)
	if err != nil {
		return "", err
	}

	if response.Data.ID != "" {
		return response.Data.ID, nil
	}

	return "", errors.New("No response was provided")
}
