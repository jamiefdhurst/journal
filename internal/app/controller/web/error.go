package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Error Display a 404 not found page
type Error struct {
	controller.Super
}

// Run Error
func (c *Error) Run(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNotFound)

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.tmpl",
		"./web/templates/error.tmpl")
	template.ExecuteTemplate(response, "layout", nil)
}
