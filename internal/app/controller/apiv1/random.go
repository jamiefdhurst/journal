package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Random Controller to handle returning a random journal entry via API
type Random struct {
	controller.Super
}

// Run Random controller action
func (c *Random) Run(response http.ResponseWriter, request *http.Request) {
	container := c.Super.Container().(*app.Container)
	js := model.Journals{Container: container}
	
	// Find a random journal entry
	randomJournal := js.FindRandom()
	
	// Set content type to JSON
	response.Header().Set("Content-Type", "application/json")
	
	// Return 404 if no journal was found
	if randomJournal.ID == 0 {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	
	// Encode and return the journal
	encoder := json.NewEncoder(response)
	encoder.Encode(randomJournal)
}