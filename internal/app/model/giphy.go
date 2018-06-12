package model

import (
	"regexp"
	"strings"

	"github.com/jamiefdhurst/journal/internal/app"
)

type giphyContent struct {
	IDs      []string
	searches []string
}

// Giphys Common resource link for Giphy actions
type Giphys struct {
	Container *app.Container
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

// ExtractContentsAndSearchAPI Convert any searches, connecting to Giphy where required
func (gs *Giphys) ExtractContentsAndSearchAPI(s string) string {
	content := gs.findTags(s)
	if len(content.searches) > 0 {
		for _, i := range content.searches {
			id, err := gs.Container.Giphy.SearchForID(i)
			if err == nil {
				s = strings.Replace(s, ":gif:"+i, ":gif:id:"+id, 1)
			} else {
				s = strings.Replace(s, ":gif:"+i, "", 1)
			}
		}
	}

	return s
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
