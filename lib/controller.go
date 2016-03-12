package lib

import (
	"database/sql"
	"net/http"
)

// Controller Super-struct for all controllers.
type Controller struct {
	Db     *sql.DB
	Params []string
}

// ControllerInterface Interface to satisfy being a controller.
type ControllerInterface interface {
	SetDb(db *sql.DB)
	Run(w http.ResponseWriter, r *http.Request)
	SetParams(p []string)
}

// SetDb Set the database pointer
func (c *Controller) SetDb(db *sql.DB) {
	c.Db = db
}

// SetParams Set the current parameters on the controller.
func (c *Controller) SetParams(p []string) {
	c.Params = p
}
