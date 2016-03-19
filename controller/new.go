package controller

import (
	"journal/lib"
	"journal/model"
	"net/http"
	"text/template"
)

// New Handle creating a new entry
type New struct {
	lib.Controller
}

// Run NewC
func (c *New) Run(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/new.tmpl")
		t.ExecuteTemplate(w, "header", nil)
		t.ExecuteTemplate(w, "content", nil)
		t.ExecuteTemplate(w, "footer", nil)
		t.Execute(w, nil)
	} else {

		js := model.Journals{}
		js.SetDb(c.Db)
		js.Create(0, model.Slugify(r.FormValue("title")), r.FormValue("title"), r.FormValue("date"), r.FormValue("content"))

		http.Redirect(w, r, "/", 302)
	}
}
