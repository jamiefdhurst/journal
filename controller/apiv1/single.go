package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// Single Find and display single blog entry
type Single struct {
	controller.Super
}

// Run Single action
func (c *Single) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Db: c.Super.Db, Gs: &model.Giphys{Db: c.Super.Db}}
	journal := js.FindBySlug(c.Params[1])

	response.Header().Add("Content-Type", "application/json")
	if journal.ID == 0 {
		response.WriteHeader(http.StatusNotFound)
	} else {
		encoder := json.NewEncoder(response)
		encoder.SetEscapeHTML(false)
		encoder.Encode(journal)
	}

}
