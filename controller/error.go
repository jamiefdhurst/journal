package controller

import (
	"net/http"
	"text/template"
)

// Error Display a 404 not found page
type Error struct {
	Controller
}

// Run Error
func (c *Error) Run(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNotFound)

	template, _ := template.ParseFiles(
		"./views/_layout/default.tmpl",
		"./views/error.tmpl")
	template.ExecuteTemplate(response, "layout", nil)
}
