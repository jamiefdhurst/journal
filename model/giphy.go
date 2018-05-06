package model

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const giphyTable = "giphy"

// GiphyAPIResponse Response holder for GIPHY API call
type GiphyAPIResponse struct {
	Data GiphyAPIResponseData `json:"data"`
}

// GiphyAPIResponseData Data object within API response
type GiphyAPIResponseData struct {
	ID string `json:"id"`
}

// GiphyContent Extracted details about Giphy lookups in content
type GiphyContent struct {
	IDs      []string
	searches []string
}

// Giphy model
type Giphy struct {
	ID     int    `json:"id"`
	APIKey string `json:"apiKey"`
}

// Giphys Common database resource link for Giphy actions
type Giphys struct {
	Db Database
}

// ConvertIDsToIframes Convert any IDs in the content into <iframe> embeds
func (gs *Giphys) ConvertIDsToIframes(s string) string {
	content := gs.findTags(s)
	if len(content.IDs) > 0 {
		for _, i := range content.IDs {
			s = strings.Replace(s, ":gif:id:"+i, "<iframe src=\"https://giphy.com/embed/"+i+"\"></iframe>", 1)
		}
	}

	return s
}

// CreateTable Create the actual table
func (gs *Giphys) CreateTable() error {
	_, err := gs.Db.Exec("CREATE TABLE `" + giphyTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`apiKey` VARCHAR(64) NOT NULL" +
		")")

	return err
}

// ExtractContentsAndSearchAPI Convert any searches, connecting to Giphy where required
func (gs *Giphys) ExtractContentsAndSearchAPI(s string) string {
	content := gs.findTags(s)
	if len(content.searches) > 0 {
		for _, i := range content.searches {
			id, err := gs.searchAPI(i)
			if err == nil {
				s = strings.Replace(s, ":gif:"+i, ":gif:id:"+id, 1)
			} else {
				s = strings.Replace(s, ":gif:"+i, "", 1)
			}
		}
	}

	return s
}

// GetAPIKey Get the current API key
func (gs *Giphys) GetAPIKey() string {
	rows, _ := gs.Db.Query("SELECT * FROM `" + giphyTable + "`")
	giphys := gs.loadFromRows(rows)

	if len(giphys) == 1 {
		return giphys[0].APIKey
	}

	return ""
}

// InputNewAPIKey Load a new API key into the Giphy table
func (gs *Giphys) InputNewAPIKey(reader io.Reader) error {
	bufferReader := bufio.NewReader(reader)
	fmt.Print("Enter GIPHY API key: ")
	apiKey, _ := bufferReader.ReadString('\n')
	gs.UpdateAPIKey(strings.Replace(apiKey, "\n", "", -1))

	return nil
}

// UpdateAPIKey Update/save the API key
func (gs *Giphys) UpdateAPIKey(apiKey string) Giphy {
	g := Giphy{1, apiKey}
	gs.save(g)

	return g
}

func (gs Giphys) findIds(s string) []string {
	reIDs := regexp.MustCompile(":gif:id:(\\w+)")
	onlyIDs := []string{}
	IDs := reIDs.FindAllStringSubmatch(s, -1)
	for _, i := range IDs {
		onlyIDs = append(onlyIDs, i[1])
	}

	return onlyIDs
}

func (gs Giphys) findSearches(s string) []string {
	reSearches := regexp.MustCompile("gif:([\\w\\-]+)")
	onlySearches := []string{}
	searches := reSearches.FindAllStringSubmatch(s, -1)
	for _, j := range searches {
		if j[1] != "id" {
			onlySearches = append(onlySearches, j[1])
		}
	}

	return onlySearches
}

func (gs Giphys) findTags(s string) GiphyContent {
	return GiphyContent{gs.findIds(s), gs.findSearches(s)}
}

func (gs Giphys) loadFromRows(rows *sql.Rows) []Giphy {
	var id int
	var apiKey string

	defer rows.Close()
	giphys := []Giphy{}
	for rows.Next() {
		rows.Scan(&id, &apiKey)
		giphys = append(giphys, Giphy{id, apiKey})
	}

	return giphys
}

func (gs *Giphys) save(g Giphy) Giphy {
	gs.Db.Exec("REPLACE INTO `"+giphyTable+"` (`id`, `apiKey`) VALUES(?,?)", strconv.Itoa(g.ID), g.APIKey)

	return g
}

func (gs *Giphys) searchAPI(s string) (string, error) {
	apiKey := gs.GetAPIKey()
	if apiKey == "" {
		return "", errors.New("No API key was found for GIPHY")
	}

	// Perform search
	url := "https://api.giphy.com/v1/gifs/random?api_key=" + apiKey + "&tag=" + s + "&rating=G"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "Journal")
	client := &http.Client{}
	rs, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rs.Body.Close()

	response := GiphyAPIResponse{}
	json.NewDecoder(rs.Body).Decode(&response)

	if response.Data.ID != "" {
		return response.Data.ID, nil
	}

	return "", errors.New("No response was provided")
}
