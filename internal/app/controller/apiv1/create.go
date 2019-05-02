package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Create Create a new entry via API
type Create struct {
	controller.Super
}

// Run Create action
func (c *Create) Run(response http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	var journalRequest = journalFromJSON{}
	err := decoder.Decode(&journalRequest)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
	} else {
		if journalRequest.Title == "" || journalRequest.Content == "" || journalRequest.Date == "" {
			response.WriteHeader(http.StatusBadRequest)
		} else {
			journal := model.Journal{ID: 0, Slug: model.Slugify(journalRequest.Title), Title: journalRequest.Title, Date: journalRequest.Date, Content: journalRequest.Content}
			js := model.Journals{Container: c.Super.Container.(*app.Container), Gs: model.GiphyAdapter(c.Super.Container.(*app.Container))}
			journal = js.Save(journal)
			response.WriteHeader(http.StatusCreated)
			encoder := json.NewEncoder(response)
			encoder.SetEscapeHTML(false)
			encoder.Encode(journal)
		}
	}
}
