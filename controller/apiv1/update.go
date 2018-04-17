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

// Run Update
func (c *Update) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	j := js.FindBySlug(c.Params[1])

	w.Header().Add("Content-Type", "application/json")
	if j.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		var j2 = journalFromJSON{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&j2)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			// Update only fields that are present
			if j2.Title != "" {
				j.Title = j2.Title
			}
			if j2.Date != "" {
				j.Date = j2.Date
			}
			if j2.Content != "" {
				j.Content = j2.Content
			}
			js.Update(j)
			encoder := json.NewEncoder(w)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(j); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
