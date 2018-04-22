package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

const giphyTable = "giphy"

// ConvertGiphyIDsToIframes Convert any IDs in the content into <iframe> embeds
func ConvertGiphyIDsToIframes(s string) string {
	content := findGiphyTags(s)
	if len(content.IDs) > 0 {
		for _, i := range content.IDs {
			s = strings.Replace(s, ":gif:id:"+i, "<iframe src=\"https://giphy.com/embed/"+i+"\"></iframe>", 1)
		}
	}

	return s
}

// CreateGiphyTable Create the actual table
func CreateGiphyTable() error {
	_, err := db.Exec("CREATE TABLE `" + giphyTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`apiKey` VARCHAR(64) NOT NULL" +
		")")

	return err
}

// ExtractContentsAndSearchGiphyAPI Convert any searches, connecting to Giphy where required
func ExtractContentsAndSearchGiphyAPI(s string) string {
	content := findGiphyTags(s)
	if len(content.searches) > 0 {
		for _, i := range content.searches {
			id, err := searchGiphyAPI(i)
			if err == nil {
				s = strings.Replace(s, ":gif:"+i, ":gif:id:"+id, 1)
			} else {
				s = strings.Replace(s, ":gif:"+i, "", 1)
			}
		}
	}

	return s
}

// GetGiphyAPIKey Get the current API key
func GetGiphyAPIKey() string {
	rows, _ := db.Query("SELECT * FROM `" + giphyTable + "`")
	giphys := loadGiphysFromRows(rows)

	if len(giphys) == 1 {
		return giphys[0].APIKey
	}

	return ""
}

// UpdateGiphyAPIKey Update/save the API key
func UpdateGiphyAPIKey(apiKey string) Giphy {
	g := Giphy{1, apiKey}
	g.save()

	return g
}

func findGiphyIds(s string) []string {
	reIDs := regexp.MustCompile(":gif:id:(\\w+)")
	onlyIDs := []string{}
	IDs := reIDs.FindAllStringSubmatch(s, -1)
	for _, i := range IDs {
		onlyIDs = append(onlyIDs, i[1])
	}

	return onlyIDs
}

func findGiphySearches(s string) []string {
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

func findGiphyTags(s string) GiphyContent {
	return GiphyContent{findGiphyIds(s), findGiphySearches(s)}
}

func loadGiphysFromRows(rows *sql.Rows) []Giphy {
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

func searchGiphyAPI(s string) (string, error) {
	apiKey := GetGiphyAPIKey()
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

func (g *Giphy) save() {
	var stmt *sql.Stmt
	stmt, _ = db.Prepare("REPLACE INTO `" + giphyTable + "` (`id`, `apiKey`) VALUES(?,?)")
	stmt.Exec(g.ID, g.APIKey)
}
