package apiv1

import (
	"encoding/json"
	"journal/controller"
	"journal/model"
	"net/http"
)

// Single Find and display single blog entry
type Single struct {
	controller.Controller
}

type singleData struct {
	Journal model.Journal
}

// Run Single
func (c *Single) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	j := js.FindBySlug(c.Params[1])

	w.Header().Add("Content-Type", "application/json")
	if j.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(j); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}
