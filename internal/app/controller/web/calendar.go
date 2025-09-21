package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Calendar Handle displaying a calendar with blog entries for given days
type Calendar struct {
	controller.Super
}

type calendarTemplateData struct {
	Container interface{}
}

// Run Calendar action
func (c *Calendar) Run(response http.ResponseWriter, request *http.Request) {

	data := calendarTemplateData{}

	container := c.Super.Container().(*app.Container)
	data.Container = container

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/calendar.html.tmpl")
	template.ExecuteTemplate(response, "layout", data)
}
