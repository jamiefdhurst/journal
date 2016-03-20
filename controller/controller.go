package controller

import "net/http"

// Controller Super-struct for all controllers.
type Controller struct {
	Params []string
}

// Interface Interface to satisfy being a controller.
type Interface interface {
	Run(w http.ResponseWriter, r *http.Request)
	SetParams(p []string)
}

// SetParams Set the current parameters on the controller.
func (c *Controller) SetParams(p []string) {
	c.Params = p
}
