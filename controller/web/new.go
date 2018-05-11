package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// New Handle creating a new entry
type New struct {
	controller.Super
	Error   bool
	Journal model.Journal
}

// Run New action
func (c *New) Run(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		query := request.URL.Query()
		if query["error"] != nil {
			c.Error = true
		}

		template, _ := template.ParseFiles(
			"./views/_layout/default.tmpl",
			"./views/new.tmpl",
			"./views/_partial/form.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	} else {
		if request.FormValue("title") == "" || request.FormValue("date") == "" || request.FormValue("content") == "" {
			http.Redirect(response, request, "/new?error=1", 302)
			return
		}

		js := model.Journals{Db: c.Super.Db, Gs: &model.Giphys{Db: c.Super.Db}}
		journal := model.Journal{ID: 0, Slug: model.Slugify(request.FormValue("title")), Title: model.Slugify(request.FormValue("title")), Date: request.FormValue("date"), Content: request.FormValue("content")}
		js.Save(journal)

		http.Redirect(response, request, "/?saved=1", 302)
	}
}
