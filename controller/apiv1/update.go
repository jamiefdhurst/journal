package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// Update Update an existing entry via API
type Update struct {
	controller.Controller
}

// Run Update action
func (c *Update) Run(response http.ResponseWriter, request *http.Request) {

	journal := model.FindJournalBySlug(c.Params[1])

	response.Header().Add("Content-Type", "application/json")
	if journal.ID == 0 {
		response.WriteHeader(http.StatusNotFound)
	} else {
		var journalRequest = journalFromJSON{}
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&journalRequest)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
		} else {
			// Update only fields that are present
			if journalRequest.Title != "" {
				journal.Title = journalRequest.Title
			}
			if journalRequest.Date != "" {
				journal.Date = journalRequest.Date
			}
			if journalRequest.Content != "" {
				journal.Content = journalRequest.Content
			}
			journal.Save()
			encoder := json.NewEncoder(response)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(journal); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
