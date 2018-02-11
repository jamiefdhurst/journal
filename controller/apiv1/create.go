package apiv1

import (
	"encoding/json"
	"journal/controller"
	"journal/model"
	"net/http"
)

// Create Create a new entry via API
type Create struct {
	controller.Controller
}

type journalFromJSON struct {
	Title   string
	Date    string
	Content string
}

// Run Create
func (c *Create) Run(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var j = journalFromJSON{}
	err := decoder.Decode(&j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if j.Title == "" || j.Content == "" || j.Date == "" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			js := model.Journals{}
			journal := js.Create(0, model.Slugify(j.Title), j.Title, j.Date, j.Content)
			encoder := json.NewEncoder(w)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(journal); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}
