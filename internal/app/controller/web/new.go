package web

import (
	"net/http"
	"time"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// New Handle creating a new entry
type New struct {
	controller.Super
}

// Run New action
func (c *New) Run(response http.ResponseWriter, request *http.Request) {

	data := formTemplateData{}
	container := c.Super.Container().(*app.Container)
	if !container.Configuration.EnableCreate {
		RunBadRequest(response, request, c.Super.Container)
		return
	}

	if request.Method == "GET" {
		data.Journal.Date = time.Now().Format("2006-01-02")
		RenderFromSession(c, data, "new", response)
	} else {
		if !model.Validate(request.FormValue("title"), request.FormValue("date"), request.FormValue("content")) {
			// Save form data to session before redirecting
			c.Session().Set("form_data", map[string]string{
				"title":   request.FormValue("title"),
				"date":    request.FormValue("date"),
				"content": request.FormValue("content"),
			})
			c.Session().AddFlash("error")
			c.SaveSession(response)
			http.Redirect(response, request, "/new", http.StatusFound)
			return
		}

		js := model.Journals{Container: container}
		journal := model.Journal{ID: 0, Slug: model.Slugify(request.FormValue("title")), Title: request.FormValue("title"), Date: request.FormValue("date"), Content: request.FormValue("content")}
		js.Save(journal)

		c.Session().AddFlash("saved")
		c.SaveSession(response)
		http.Redirect(response, request, "/", http.StatusFound)
	}
}
