package controller

import (
	"net/http"
	"text/template"
)

// Error Display a 404
type Error struct {
	Controller
}

// Run Error
func (c *Error) Run(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	t, _ := template.ParseFiles(
		"./views/_layout/default.tmpl",
		"./views/error.tmpl")
	t.ExecuteTemplate(w, "layout", nil)
}
