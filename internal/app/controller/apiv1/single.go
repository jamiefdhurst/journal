package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Single Find and display single blog entry
type Single struct {
	controller.Super
}

// Run Single action
func (c *Single) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Container: c.Super.Container.(*app.Container)}
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
