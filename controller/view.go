package controller

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/model"
)

// View Handle displaying individual entry
type View struct {
	Controller
}

type viewData struct {
	Journal model.Journal
}

// Run View
func (c *View) Run(w http.ResponseWriter, r *http.Request) {

	js := model.Journals{}
	j := js.FindBySlug(c.Params[1])

	if j.ID == 0 {
		e := Error{}
		e.Run(w, r)
	} else {
		j.Content = model.ConvertIDsForDisplay(j.Content)
		data := viewData{j}
		t, _ := template.ParseFiles(
			"./src/journal/views/_layout/default.tmpl",
			"./src/journal/views/view.tmpl")
		t.ExecuteTemplate(w, "layout", data)
	}
}
