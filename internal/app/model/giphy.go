package model

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jamiefdhurst/journal/pkg/adapter/giphy"
	"github.com/jamiefdhurst/journal/pkg/database"
	"github.com/jamiefdhurst/journal/pkg/database/rows"
)

const giphyTable = "giphy"

type giphyContent struct {
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
	Client giphy.Adapter
	Db     database.Database
}

// GiphysExtractor Interface for extracting a Giphy search
type GiphysExtractor interface {
	ExtractContentsAndSearchAPI(s string) string
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
			id, err := gs.Client.SearchForID(i)
			if err == nil {
				s = strings.Replace(s, ":gif:"+i, ":gif:id:"+id, 1)
			} else {
				s = strings.Replace(s, ":gif:"+i, "", 1)
			}
		}
	}

	return s
}

// GetAPIKey Get the API key from Giphy to be used in client
func (gs *Giphys) GetAPIKey() string {
	rows, err := gs.Db.Query("SELECT * FROM `" + giphyTable + "`")
	if err != nil {
		return ""
	}
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
	_, err := gs.updateAPIKey(strings.Replace(apiKey, "\n", "", -1))

	return err
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

func (gs Giphys) findTags(s string) giphyContent {
	return giphyContent{gs.findIds(s), gs.findSearches(s)}
}

func (gs Giphys) loadFromRows(rows rows.Rows) []Giphy {
	defer rows.Close()
	giphys := []Giphy{}
	for rows.Next() {
		g := Giphy{}
		rows.Scan(&g.ID, &g.APIKey)
		log.Println(g.APIKey)
		giphys = append(giphys, g)
	}

	return giphys
}

func (gs *Giphys) save(g Giphy) (Giphy, error) {
	_, err := gs.Db.Exec("REPLACE INTO `"+giphyTable+"` (`id`, `apiKey`) VALUES(?,?)", strconv.Itoa(g.ID), g.APIKey)

	return g, err
}

func (gs *Giphys) updateAPIKey(apiKey string) (Giphy, error) {
	g := Giphy{1, apiKey}

	return gs.save(g)
}
