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
	Journal model.Journal
	Next    model.Journal
	Prev    model.Journal
}

// Run View action
func (c *View) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Container: c.Super.Container.(*app.Container), Gs: model.GiphyAdapter(c.Super.Container.(*app.Container))}
	c.Journal = js.FindBySlug(c.Params[1])

	if c.Journal.ID == 0 {
		RunBadRequest(response, request, c.Super.Container)
	} else {
		c.Next = js.FindNext(c.Journal.ID)
		c.Prev = js.FindPrev(c.Journal.ID)
		gs := model.Giphys{}
		c.Journal.Content = gs.ConvertIDsToIframes(c.Journal.Content)
		template, _ := template.ParseFiles(
			"./web/templates/_layout/default.tmpl",
			"./web/templates/view.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	}
}
