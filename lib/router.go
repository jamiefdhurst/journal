package lib

import (
	"database/sql"
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
	controller ControllerInterface
}

// Router Contian routes
type Router struct {
	db     *sql.DB
	err    ControllerInterface
	routes []Route
}

// Add A new route
func (m *Router) Add(t string, u string, a bool, c ControllerInterface) {
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
			route.controller.SetDb(m.db)
			route.controller.Run(w, r)
			return
		}

		// Attempt regex match
		if route.matchable {
			matched, _ := regexp.MatchString(route.uri, r.URL.Path)
			if matched && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
				re := regexp.MustCompile(route.uri)
				route.controller.SetDb(m.db)
				route.controller.SetParams(re.FindAllString(r.URL.Path, -1))
				route.controller.Run(w, r)
				return
			}
		}
	}

	m.err.Run(w, r)
}

// SetDb Set the db
func (m *Router) SetDb(db *sql.DB) {
	m.db = db
}

// SetErr Set the err controller
func (m *Router) SetErr(err ControllerInterface) {
	m.err = err
}
