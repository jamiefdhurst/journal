package controller

import (
	"net/http"
)

// Controller Main interface for controllers
type Controller interface {
	Init(app interface{}, params []string)
	Run(response http.ResponseWriter, request *http.Request)
}

// Super Super-struct for all controllers.
type Super struct {
	Controller
	Container interface{}
	Params    []string
}

// Init Initialise the controller
func (c *Super) Init(app interface{}, params []string) {
	c.Container = app
	c.Params = params
}
