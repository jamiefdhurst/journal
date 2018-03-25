package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

const giphyTable = "giphy"

// Giphys Collection of Giphy API key entries
type Giphys struct {
	Giphys []Giphy
}

// Giphy model
type Giphy struct {
	ID     int    `json:"id"`
	APIKey string `json:"apiKey"`
}

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

// ConvertIDsForDisplay Convert any IDs in the content into <iframe> embeds
func ConvertIDsForDisplay(s string) string {
	content := extractGiphy(s)
	if len(content.IDs) > 0 {
		for _, i := range content.IDs {
			s = strings.Replace(s, ":gif:id:"+i, "<iframe src=\"https://giphy.com/embed/"+i+"\"></iframe>", 1)
		}
	}

	return s
}

// ConvertSearchesToIDs Convert any searches, connecting to GIPHY where required
func ConvertSearchesToIDs(s string) string {
	content := extractGiphy(s)
	if len(content.searches) > 0 {
		for _, i := range content.searches {
			id, err := searchGiphy(i)
			if err == nil {
				s = strings.Replace(s, ":gif:"+i, ":gif:id:"+id, 1)
			} else {
				s = strings.Replace(s, ":gif:"+i, "", 1)
			}
		}
	}

	return s
}

func extractGiphy(s string) GiphyContent {

	// Extract IDs
	reIDs := regexp.MustCompile(":gif:id:(\\w+)")
	onlyIDs := []string{}
	IDs := reIDs.FindAllStringSubmatch(s, -1)
	for _, i := range IDs {
		onlyIDs = append(onlyIDs, i[1])
	}

	// Extract searches
	reSearches := regexp.MustCompile("gif:([\\w\\-]+)")
	onlySearches := []string{}
	searches := reSearches.FindAllStringSubmatch(s, -1)
	for _, j := range searches {
		if j[1] != "id" {
			onlySearches = append(onlySearches, j[1])
		}
	}

	return GiphyContent{onlyIDs, onlySearches}
}

// GetKey Get the current API key
func (gs *Giphys) GetKey() string {
	// Attempt to find the entry
	rows, _ := db.Query("SELECT * FROM `" + giphyTable + "`")

	defer rows.Close()
	for rows.Next() {
		gs.load(rows)
	}

	if len(gs.Giphys) == 1 {
		return gs.Giphys[0].APIKey
	}

	return ""
}

func (gs *Giphys) load(rows *sql.Rows) {
	var id int
	var apiKey string

	rows.Scan(&id, &apiKey)
	g := Giphy{id, apiKey}
	gs.Giphys = append(gs.Giphys, g)
}

func (gs *Giphys) save(g Giphy) Giphy {
	var stmt *sql.Stmt
	stmt, _ = db.Prepare("REPLACE INTO `" + giphyTable + "` (`id`, `apiKey`) VALUES(?,?)")
	stmt.Exec(g.ID, g.APIKey)

	return g
}

func searchGiphy(s string) (string, error) {
	// Get API key
	gs := Giphys{}
	apiKey := gs.GetKey()
	if apiKey == "" {
		return "", errors.New("No API key was found for GIPHY")
	}

	// Perform search
	var url string
	url = "https://api.giphy.com/v1/gifs/random?api_key=" + apiKey + "&tag=" + s + "&rating=G"
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

// Update Update/save the API key
func (gs *Giphys) Update(apiKey string) Giphy {
	g := Giphy{1, apiKey}
	g = gs.save(g)

	gs.Giphys = append(gs.Giphys, g)

	return g
}

// GiphyCreateTable Create the actual table
func GiphyCreateTable() error {
	_, err := db.Exec("CREATE TABLE `" + giphyTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`apiKey` VARCHAR(64) NOT NULL" +
		")")

	return err
}
