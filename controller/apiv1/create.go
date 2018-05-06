package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
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
			js := model.Journals{Db: c.Super.Db}
			journal = js.Save(journal)
			encoder := json.NewEncoder(response)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(journal); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
