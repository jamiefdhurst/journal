package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// List Display all blog entries as JSON
type List struct {
	controller.Super
}

// Run List action
func (c *List) Run(response http.ResponseWriter, request *http.Request) {

	js := model.NewJournalStore(c.Super.Container.(*app.Container), model.GiphyAdapter(c.Super.Container.(*app.Container)))
	journals := js.FetchAll()
	response.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	encoder.Encode(journals)
}
