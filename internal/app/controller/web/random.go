package web

import (
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Random Controller to handle redirecting to a random journal entry
type Random struct {
	controller.Super
}

// Run Random controller action
func (c *Random) Run(response http.ResponseWriter, request *http.Request) {
	container := c.Super.Container().(*app.Container)
	js := model.Journals{Container: container}

	// Find a random journal entry
	randomJournal := js.FindRandom()

	// Redirect to the entry or home page if none found
	if randomJournal.ID > 0 {
		http.Redirect(response, request, "/"+randomJournal.Slug, http.StatusFound)
	} else {
		http.Redirect(response, request, "/", http.StatusFound)
	}
}
