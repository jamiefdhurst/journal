package lib

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/jamiefdhurst/journal/controller"
)

// Route A route contains a method (GET), URI, and a controller
type Route struct {
	method     string
	uri        string
	matchable  bool
	controller controller.Interface
}

// Router A router contains routes and links back to the application and implements the ServeHTTP interface
type Router struct {
	err    controller.Interface
	routes []Route
	app    *App
}

// Add Create and add a new route into the router
func (m *Router) Add(method string, uri string, matchable bool, controller controller.Interface) {
	r := Route{method, uri, matchable, controller}
	m.routes = append(m.routes, r)
}

// ServeHTTP Serve a given HTTP request
func (m *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Debug output into the console
	log.Printf("%s: %s", r.Method, r.URL.Path)

	// Attempt to serve a file first - still uses the full GOPATH
	file := "src/github.com/jamiefdhurst/journal/public" + r.URL.Path
	if r.URL.Path != "/" {
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			http.ServeFile(w, r, file)
			return
		}
	}

	for _, route := range m.routes {
		if r.URL.Path == route.uri && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
			route.controller.Run(w, r)
			return
		}

		// Attempt regex match
		if route.matchable {
			matched, _ := regexp.MatchString(route.uri, r.URL.Path)
			if matched && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
				re := regexp.MustCompile(route.uri)
				route.controller.SetParams(re.FindStringSubmatch(r.URL.Path))
				route.controller.Run(w, r)
				return
			}
		}
	}

	m.err.Run(w, r)
}

// NewRouter Create a new router with an error controller provided
func NewRouter(s *Server, e controller.Interface) Router {
	var r []Route

	return Router{e, r, s}
}
