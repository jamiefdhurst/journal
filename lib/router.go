package lib

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jamiefdhurst/journal/controller"
)

// Route A route contains a method (GET), URI, and a controller
type Route struct {
	method     string
	regexURI   string
	controller controller.Interface
}

// Router A router contains routes and links back to the application and implements the ServeHTTP interface
type Router struct {
	app             *App
	routes          []Route
	errorController controller.Interface
}

func (m *Router) convertSimpleURIToRegex(uri string) string {
	uri = strings.Replace(uri, "/", "\\/", -1)

	// Match slugs
	uri = strings.Replace(uri, "[%s]", "([\\w\\-]+)", -1)

	// Match IDs
	uri = strings.Replace(uri, "[%d]", "(\\d+)", -1)

	// Match anything
	uri = strings.Replace(uri, "[%a]", "(.+?)", -1)

	return uri
}

// Get Create and add a new route into the router to handle a GET request
func (m *Router) Get(uri string, controller controller.Interface) {
	m.routes = append(m.routes, Route{"GET", m.convertSimpleURIToRegex(uri), controller})
}

// Post Create and add a new route into the router to handle a POST request
func (m *Router) Post(uri string, controller controller.Interface) {
	m.routes = append(m.routes, Route{"POST", m.convertSimpleURIToRegex(uri), controller})
}

// Put Create and add a new route into the router to handle a PUT request
func (m *Router) Put(uri string, controller controller.Interface) {
	m.routes = append(m.routes, Route{"PUT", m.convertSimpleURIToRegex(uri), controller})
}

// ServeHTTP Serve a given HTTP request
func (m *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Debug output into the console
	log.Printf("%s: %s", r.Method, r.URL.Path)

	// Attempt to serve a file first
	if r.URL.Path != "/" {
		file := "public" + r.URL.Path
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			http.ServeFile(w, r, file)
			return
		}
	}

	// Go through each route and attempt to match
	for _, route := range m.routes {
		matched, _ := regexp.MatchString(route.regexURI, r.URL.Path)
		if matched && (r.Method == route.method || (r.Method == "" && route.method == "GET")) {
			re := regexp.MustCompile(route.regexURI)
			route.controller.SetParams(re.FindStringSubmatch(r.URL.Path))
			route.controller.Run(w, r)
			return
		}
	}

	m.errorController.Run(w, r)
}
