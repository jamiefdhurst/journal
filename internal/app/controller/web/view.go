package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// View Handle displaying individual entry
type View struct {
	controller.Super
}

type viewTemplateData struct {
	Container interface{}
	Journal   model.Journal
	Next      model.Journal
	Prev      model.Journal
}

// Run View action
func (c *View) Run(response http.ResponseWriter, request *http.Request) {

	data := viewTemplateData{}
	container := c.Super.Container().(*app.Container)
	data.Container = container
	js := model.Journals{Container: container}
	data.Journal = js.FindBySlug(c.Params()[1])

	if data.Journal.ID == 0 {
		RunBadRequest(response, request, c.Super.Container)
	} else {
		data.Next = js.FindNext(data.Journal.ID)
		data.Prev = js.FindPrev(data.Journal.ID)
		template, _ := template.ParseFiles(
			"./web/templates/_layout/default.html.tmpl",
			"./web/templates/view.html.tmpl")
		template.ExecuteTemplate(response, "layout", data)
	}
}
