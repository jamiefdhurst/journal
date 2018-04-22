package controller

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/model"
)

// Edit Handle updating an existing entry
type Edit struct {
	Controller
	Error   bool
	Journal model.Journal
}

// Run Edit action
func (c *Edit) Run(response http.ResponseWriter, request *http.Request) {

	c.Journal = model.FindJournalBySlug(c.Params[1])

	if c.Journal.ID == 0 {
		errorController := Error{}
		errorController.Run(response, request)
	} else {

		if request.Method == "GET" {
			query := request.URL.Query()
			if query["error"] != nil {
				c.Error = true
			}
			template, _ := template.ParseFiles(
				"./views/_layout/default.tmpl",
				"./views/edit.tmpl",
				"./views/_partial/form.tmpl")
			template.ExecuteTemplate(response, "layout", c)
		} else {
			if request.FormValue("title") == "" || request.FormValue("date") == "" || request.FormValue("content") == "" {
				http.Redirect(response, request, "/"+c.Journal.Slug+"/edit?error=1", 302)
			}

			c.Journal.Title = request.FormValue("title")
			c.Journal.Date = request.FormValue("date")
			c.Journal.Content = request.FormValue("content")
			c.Journal.Save()

			http.Redirect(response, request, "/?saved=1", 302)
		}
	}

}
