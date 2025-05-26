package controller

import (
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/session"
)

// Controller Main interface for controllers
type Controller interface {
	Init(app interface{}, params []string, request *http.Request)
	Run(response http.ResponseWriter, request *http.Request)
	Container() interface{}
	Host() string
	Params() []string
	SaveSession(w http.ResponseWriter)
	Session() *session.Session
}

// Super Super-struct for all controllers.
type Super struct {
	Controller
	container       interface{}
	disableTracking bool
	host            string
	params          []string
	session         *session.Session
	sessionStore    session.Store
}

// Init Initialise the controller
func (c *Super) Init(app interface{}, params []string, request *http.Request) {
	c.container = app
	c.host = request.Host
	c.params = params
	c.sessionStore = session.NewDefaultStore("defaultdefaultdefaultdefault1234")
	c.session, _ = c.sessionStore.Get(request)

	c.trackVisit(request)
}

func (c *Super) Container() interface{} {
	return c.container
}

func (c *Super) DisableTracking() {
	c.disableTracking = true
}

func (c *Super) Host() string {
	return c.host
}

func (c *Super) Params() []string {
	return c.params
}

// SaveSession saves the session with the current response
func (c *Super) SaveSession(w http.ResponseWriter) {
	c.sessionStore.Save(w)
}

// Session gets the private session value
func (c *Super) Session() *session.Session {
	return c.session
}

func (c *Super) trackVisit(request *http.Request) {
	if c.disableTracking {
		return
	}

	if c.container == nil || request == nil || request.URL == nil {
		return
	}

	appContainer, ok := c.container.(*app.Container)
	if !ok || appContainer.Db == nil {
		return
	}

	visits := model.Visits{Container: appContainer}
	visits.RecordVisit(request.URL.Path)
}
