package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// List Display all blog entries as JSON
type List struct {
	controller.Super
}

// Run List action
func (c *List) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Db: c.Super.Db, Gs: &model.Giphys{Db: c.Super.Db}}
	journals := js.FetchAll()
	response.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	encoder.Encode(journals)
}
