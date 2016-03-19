package controller

import (
	"journal/lib"
	"journal/model"
	"net/http"
	"text/template"
)

// Index Handle displaying all blog entries
type Index struct {
	lib.Controller
}

type indexData struct {
	Journals []model.Journal
}

// Run Index
func (c *Index) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	js.SetDb(c.Db)
	js.FetchAll()
	data := indexData{js.Journals}

	t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/index.tmpl")
	t.ExecuteTemplate(w, "header", nil)
	t.ExecuteTemplate(w, "content", data)
	t.ExecuteTemplate(w, "footer", nil)
	t.Execute(w, nil)
}
