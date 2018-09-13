package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
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
			"./web/templates/_layout/default.tmpl",
			"./web/templates/new.tmpl",
			"./web/templates/_partial/form.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	} else {
		if request.FormValue("title") == "" || request.FormValue("date") == "" || request.FormValue("content") == "" {
			http.Redirect(response, request, "/new?error=1", 302)
			return
		}

		js := model.Journals{Container: c.Super.Container.(*app.Container), Gs: &model.Giphys{Container: c.Super.Container.(*app.Container)}}
		journal := model.Journal{ID: 0, Slug: model.Slugify(request.FormValue("title")), Title: request.FormValue("title"), Date: request.FormValue("date"), Content: request.FormValue("content")}
		js.Save(journal)

		http.Redirect(response, request, "/?saved=1", 302)
	}
}
