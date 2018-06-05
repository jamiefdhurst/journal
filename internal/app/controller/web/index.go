package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Index Handle displaying all blog entries
type Index struct {
	controller.Super
	Journals []model.Journal
	Saved    bool
}

// Run Index action
func (c *Index) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Db: c.Super.Db}
	c.Journals = js.FetchAll()
	query := request.URL.Query()
	if query["saved"] != nil {
		c.Saved = true
	}

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.tmpl",
		"./web/templates/index.tmpl")
	template.ExecuteTemplate(response, "layout", c)
}
