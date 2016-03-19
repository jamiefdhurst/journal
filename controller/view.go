package controller

import (
	"journal/lib"
	"journal/model"
	"net/http"
	"text/template"
)

// View Handle displaying individual entry
type View struct {
	lib.Controller
}

type viewData struct {
	Journal model.Journal
}

// Run View
func (c *View) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	js.SetDb(c.Db)
	j := js.FindBySlug(c.Params[0])

	if j.ID == 0 {
		e := Error{}
		e.Run(w, r)
	} else {
		data := viewData{j}
		t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/view.tmpl")
		t.ExecuteTemplate(w, "header", nil)
		t.ExecuteTemplate(w, "content", data)
		t.ExecuteTemplate(w, "footer", nil)
		t.Execute(w, nil)
	}
}
