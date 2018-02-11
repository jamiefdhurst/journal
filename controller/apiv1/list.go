package apiv1

import (
	"encoding/json"
	"journal/controller"
	"journal/model"
	"net/http"
)

// List Display all blog entries as JSON
type List struct {
	controller.Controller
}

// Run List
func (c *List) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	js.FetchAll()
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(js.Journals); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
