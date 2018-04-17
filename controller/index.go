package controller

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/model"
)

// Index Handle displaying all blog entries
type Index struct {
	Controller
}

type indexData struct {
	Journals []model.Journal
	Saved    bool
}

// Run Index
func (c *Index) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	js.FetchAll()
	data := indexData{js.Journals, false}
	query := r.URL.Query()
	if query["saved"] != nil {
		data.Saved = true
	}

	t, _ := template.ParseFiles(
		"./src/journal/views/_layout/default.tmpl",
		"./src/journal/views/index.tmpl")
	t.ExecuteTemplate(w, "layout", data)
}
