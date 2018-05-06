package controller

import (
	"net/http"

	"github.com/jamiefdhurst/journal/model"
)

// Controller Main interface for controllers
type Controller interface {
	Init(db model.Database, params []string)
	Run(response http.ResponseWriter, request *http.Request)
}

// Super Super-struct for all controllers.
type Super struct {
	Controller
	Db     model.Database
	Params []string
}

// Init Initialise the controller
func (c *Super) Init(db model.Database, params []string) {
	c.Db = db
	c.Params = params
}
