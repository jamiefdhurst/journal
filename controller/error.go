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
		"./src/github.com/jamiefdhurst/journal/views/_layout/default.tmpl",
		"./src/github.com/jamiefdhurst/journal/views/error.tmpl")
	t.ExecuteTemplate(w, "layout", nil)
}
