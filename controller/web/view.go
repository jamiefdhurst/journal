package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

// View Handle displaying individual entry
type View struct {
	controller.Super
	Journal model.Journal
}

// Run View action
func (c *View) Run(response http.ResponseWriter, request *http.Request) {

	js := model.Journals{Db: c.Super.Db}
	c.Journal = js.FindBySlug(c.Params[1])

	if c.Journal.ID == 0 {
		errorController := Error{}
		errorController.Run(response, request)
	} else {
		gs := model.Giphys{Db: c.Super.Db}
		c.Journal.Content = gs.ConvertIDsToIframes(c.Journal.Content)
		template, _ := template.ParseFiles(
			"./views/_layout/default.tmpl",
			"./views/view.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	}
}
