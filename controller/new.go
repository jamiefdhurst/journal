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
		data := map[string]bool{
			"Error": false,
		}
		query := r.URL.Query()
		if query["error"] != nil {
			data["Error"] = true
		}

		t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/new.tmpl")
		t.ExecuteTemplate(w, "header", nil)
		t.ExecuteTemplate(w, "content", data)
		t.ExecuteTemplate(w, "footer", nil)
		t.Execute(w, nil)
	} else {
		if r.FormValue("title") == "" || r.FormValue("date") == "" || r.FormValue("content") == "" {
			http.Redirect(w, r, "/new?error=1", 302)
		}

		js := model.Journals{}
		js.SetDb(c.Db)
		js.Create(0, model.Slugify(r.FormValue("title")), r.FormValue("title"), r.FormValue("date"), r.FormValue("content"))

		http.Redirect(w, r, "/?saved=1", 302)
	}
}
