package controller

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/model"
)

// View Handle displaying individual entry
type View struct {
	Controller
	Journal model.Journal
}

// Run View action
func (c *View) Run(response http.ResponseWriter, request *http.Request) {

	c.Journal = model.FindJournalBySlug(c.Params[1])

	if c.Journal.ID == 0 {
		errorController := Error{}
		errorController.Run(response, request)
	} else {
		c.Journal.Content = model.ConvertGiphyIDsToIframes(c.Journal.Content)
		template, _ := template.ParseFiles(
			"./views/_layout/default.tmpl",
			"./views/view.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	}
}
