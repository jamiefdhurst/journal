package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// List Display all blog entries as JSON
type List struct {
	controller.Controller
}

// Run List action
func (c *List) Run(response http.ResponseWriter, request *http.Request) {

	journals := model.FetchAllJournals()
	response.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(journals); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
}
