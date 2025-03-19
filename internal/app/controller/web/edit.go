package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Edit Handle updating an existing entry
type Edit struct {
	controller.Super
	Error   bool
	Journal model.Journal
}

// Run Edit action
func (c *Edit) Run(response http.ResponseWriter, request *http.Request) {
	container := c.Super.Container.(*app.Container)
	if !container.Configuration.EnableEdit {
		RunBadRequest(response, request, c.Super.Container)
		return
	}

	js := model.Journals{Container: c.Super.Container.(*app.Container)}
	c.Journal = js.FindBySlug(c.Params[1])

	if c.Journal.ID == 0 {
		RunBadRequest(response, request, c.Super.Container)
		return
	}

	if request.Method == "GET" {
		c.Error = false
		flash := c.Session.GetFlash()
		if flash != nil && flash[0] == "error" {
			c.Error = true
			
			// Retrieve saved form data from session
			formData := c.Session.Get("form_data")
			if formData != nil {
				formMap := formData.(map[string]string)
				// Only update the Journal object if we have data in the session
				c.Journal.Title = formMap["title"]
				c.Journal.Date = formMap["date"]
				c.Journal.Content = formMap["content"]
				// Clear the form data from session after retrieving
				c.Session.Delete("form_data")
			}
		}

		c.SessionStore.Save(response)
		template, _ := template.ParseFiles(
			"./web/templates/_layout/default.html.tmpl",
			"./web/templates/edit.html.tmpl",
			"./web/templates/_partial/form.html.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	} else {
		if request.FormValue("title") == "" || request.FormValue("date") == "" || request.FormValue("content") == "" {
			// Save form data to session before redirecting
			formData := map[string]string{
				"title":   request.FormValue("title"),
				"date":    request.FormValue("date"),
				"content": request.FormValue("content"),
			}
			c.Session.Set("form_data", formData)
			c.Session.AddFlash("error")
			c.SessionStore.Save(response)
			http.Redirect(response, request, "/"+c.Journal.Slug+"/edit", http.StatusFound)
			return
		}

		c.Journal.Title = request.FormValue("title")
		c.Journal.Date = request.FormValue("date")
		c.Journal.Content = request.FormValue("content")
		js.Save(c.Journal)

		c.Session.AddFlash("saved")
		c.SessionStore.Save(response)
		http.Redirect(response, request, "/", http.StatusFound)
	}

}
