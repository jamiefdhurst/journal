package web

import (
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Edit Handle updating an existing entry
type Edit struct {
	controller.Super
}

// Run Edit action
func (c *Edit) Run(response http.ResponseWriter, request *http.Request) {
	data := formTemplateData{}
	container := c.Super.Container().(*app.Container)
	data.Container = container
	if !container.Configuration.EnableEdit {
		RunBadRequest(response, request, container)
		return
	}

	js := model.Journals{Container: container}
	data.Journal = js.FindBySlug(c.Params()[1])

	if data.Journal.ID == 0 {
		RunBadRequest(response, request, container)
		return
	}

	if request.Method == "GET" {
		RenderFromSession(c, data, "edit", response)
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
			http.Redirect(response, request, "/"+data.Journal.Slug+"/edit", http.StatusFound)
			return
		}

		data.Journal.Title = request.FormValue("title")
		data.Journal.Date = request.FormValue("date")
		data.Journal.Content = request.FormValue("content")
		js.Save(data.Journal)

		c.Session().AddFlash("saved")
		c.SaveSession(response)
		http.Redirect(response, request, "/", http.StatusFound)
	}

}
