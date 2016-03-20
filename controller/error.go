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

	t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/error.tmpl")
	t.ExecuteTemplate(w, "header", nil)
	t.ExecuteTemplate(w, "content", nil)
	t.ExecuteTemplate(w, "footer", nil)
	t.Execute(w, nil)
}
