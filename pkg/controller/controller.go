package controller

import (
	"net/http"

	"github.com/jamiefdhurst/journal/pkg/session"
)

// Controller Main interface for controllers
type Controller interface {
	Init(app interface{}, params []string, request *http.Request)
	Run(response http.ResponseWriter, request *http.Request)
}

// Super Super-struct for all controllers.
type Super struct {
	Controller
	Container    interface{}
	Params       []string
	Session      *session.Session
	SessionStore session.Store
}

// Init Initialise the controller
func (c *Super) Init(app interface{}, params []string, request *http.Request) {
	c.Container = app
	c.Params = params
	c.SessionStore = session.NewDefaultStore("defaultdefaultdefaultdefault1234")
	c.Session, _ = c.SessionStore.Get(request)
}
