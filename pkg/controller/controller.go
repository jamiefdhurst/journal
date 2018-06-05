package controller

import (
	"net/http"

	"github.com/jamiefdhurst/journal/pkg/database"
)

// Controller Main interface for controllers
type Controller interface {
	Init(db database.Database, params []string)
	Run(response http.ResponseWriter, request *http.Request)
}

// Super Super-struct for all controllers.
type Super struct {
	Controller
	Db     database.Database
	Params []string
}

// Init Initialise the controller
func (c *Super) Init(db database.Database, params []string) {
	c.Db = db
	c.Params = params
}
