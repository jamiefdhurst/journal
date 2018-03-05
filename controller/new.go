package controller

import (
	"journal/model"
	"net/http"
	"text/template"
)

// New Handle creating a new entry
type New struct {
	Controller
}

type newData struct {
	Error   bool
	Journal model.Journal
}

// Run NewC
func (c *New) Run(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := newData{}
		query := r.URL.Query()
		if query["error"] != nil {
			data.Error = true
		}

		t, _ := template.ParseFiles(
			"./src/journal/views/_layout/default.tmpl",
			"./src/journal/views/new.tmpl",
			"./src/journal/views/_partial/form.tmpl")
		t.ExecuteTemplate(w, "layout", data)
	} else {
		if r.FormValue("title") == "" || r.FormValue("date") == "" || r.FormValue("content") == "" {
			http.Redirect(w, r, "/new?error=1", 302)
		}

		js := model.Journals{}
		js.Create(0, model.Slugify(r.FormValue("title")), r.FormValue("title"), r.FormValue("date"), r.FormValue("content"))

		http.Redirect(w, r, "/?saved=1", 302)
	}
}
