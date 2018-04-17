package controller

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/model"
)

// Edit Handle updating an existing entry
type Edit struct {
	Controller
}

type editData struct {
	Error   bool
	Journal model.Journal
}

// Run Edit
func (c *Edit) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	j := js.FindBySlug(c.Params[1])

	if j.ID == 0 {
		e := Error{}
		e.Run(w, r)
	} else {

		if r.Method == "GET" {
			data := editData{false, j}
			query := r.URL.Query()
			if query["error"] != nil {
				data.Error = true
			}
			t, _ := template.ParseFiles(
				"./src/github.com/jamiefdhurst/journal/views/_layout/default.tmpl",
				"./src/github.com/jamiefdhurst/journal/views/edit.tmpl",
				"./src/github.com/jamiefdhurst/journal/views/_partial/form.tmpl")
			t.ExecuteTemplate(w, "layout", data)
		} else {
			if r.FormValue("title") == "" || r.FormValue("date") == "" || r.FormValue("content") == "" {
				http.Redirect(w, r, "/"+j.Slug+"/edit?error=1", 302)
			}

			j.Title = r.FormValue("title")
			j.Date = r.FormValue("date")
			j.Content = r.FormValue("content")
			js.Update(j)

			http.Redirect(w, r, "/?saved=1", 302)
		}
	}

}
