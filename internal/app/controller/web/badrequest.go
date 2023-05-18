package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/pkg/controller"
)

// BadRequest Display a 404 not found page
type BadRequest struct {
	controller.Super
}

// Run BadRequest
func (c *BadRequest) Run(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusNotFound)

	c.SessionStore.Save(response)
	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/error.html.tmpl")
	template.ExecuteTemplate(response, "layout", c)
}

// RunBadRequest calls the bad request from an existing controller
func RunBadRequest(response http.ResponseWriter, request *http.Request, container interface{}) {
	errorController := BadRequest{}
	errorController.Init(container, []string{}, request)
	errorController.Run(response, request)
}
