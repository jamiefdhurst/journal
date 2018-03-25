package lib

import (
	"journal/controller"
	"log"
	"net/http"
	"os"
	"regexp"
)

// Route Define a route
type Route struct {
	method     string
	uri        string
	matchable  bool
	controller controller.Interface
}

// Router Contain routes
type Router struct {
	err    controller.Interface
	routes []Route
	server *Server
}

// Add A new route
func (m *Router) Add(t string, u string, a bool, c controller.Interface) {
	r := Route{t, u, a, c}
	m.routes = append(m.routes, r)
}

// ServeHTTP Serve the HTTP request
func (m *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s: %s", r.Method, r.URL.Path)

	// Attempt static file first
	file := "src/journal/public" + r.URL.Path
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
