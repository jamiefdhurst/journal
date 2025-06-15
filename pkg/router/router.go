package router

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Server Common interface for HTTP
type Server interface {
	ListenAndServe() error
	ListenAndServeTLS(string, string) error
}

// Route A route contains a method (GET), URI, and a controller
type Route struct {
	method     string
	regexURI   string
	controller controller.Controller
}

// Router A router contains routes and links back to the application and implements the ServeHTTP interface
type Router struct {
	isHTTPS         bool `default:"false"`
	Container       interface{}
	Routes          []Route
	StaticPaths     []string
	ErrorController controller.Controller
}

func (r Router) convertSimpleURIToRegex(uri string) string {
	uri = strings.Replace(uri, "/", "\\/", -1)

	// Match slugs
	uri = strings.Replace(uri, "[%s]", "([\\w\\-]+)", -1)

	// Match IDs
	uri = strings.Replace(uri, "[%d]", "(\\d+)", -1)

	// Match anything
	uri = strings.Replace(uri, "[%a]", "(.+?)", -1)

	return "^" + uri + "$"
}

// Get Create and add a new route into the router to handle a GET request
func (r *Router) Get(uri string, controller controller.Controller) {
	r.Routes = append(r.Routes, Route{"GET", r.convertSimpleURIToRegex(uri), controller})
}

// Post Create and add a new route into the router to handle a POST request
func (r *Router) Post(uri string, controller controller.Controller) {
	r.Routes = append(r.Routes, Route{"POST", r.convertSimpleURIToRegex(uri), controller})
}

// Put Create and add a new route into the router to handle a PUT request
func (r *Router) Put(uri string, controller controller.Controller) {
	r.Routes = append(r.Routes, Route{"PUT", r.convertSimpleURIToRegex(uri), controller})
}

// ServeHTTP Serve a given HTTP request
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// Debug output into the console
	log.Printf("%s: %s", request.Method, request.URL.Path)

	// Security headers
	if r.isHTTPS {
		response.Header().Add("Strict-Transport-Security", "max-age=15552000; includeSubDomains; preload")
	}
	response.Header().Add("Content-Security-Policy", "default-src: 'self'; font-src: 'fonts.googleapis.com'; frame-src: 'none'")
	response.Header().Add("X-XSS-Protection", "mode=block")

	// Attempt to serve a file first from available static paths
	for _, staticPath := range r.StaticPaths {
		if request.URL.Path != "/" {
			file := staticPath + request.URL.Path
			_, err := os.Stat(file)
			if !os.IsNotExist(err) {
				response.Header().Add("Cache-Control", "public, max-age=15552000")
				http.ServeFile(response, request, file)
				return
			}
		}
	}

	// Go through each route and attempt to match
	var matchedController controller.Controller = r.ErrorController
	var matchedParams []string = []string{}
	for _, route := range r.Routes {
		matched, _ := regexp.MatchString(route.regexURI, request.URL.Path)
		if matched && (request.Method == route.method || (request.Method == "" && route.method == "GET")) {
			re := regexp.MustCompile(route.regexURI)
			matchedParams = re.FindStringSubmatch(request.URL.Path)
			matchedController = route.controller
			break
		}
	}

	matchedController.Init(r.Container, matchedParams, request)
	matchedController.Run(response, request)
}

// StartAndServe Start the HTTP server and listen for connections
func (r *Router) StartAndServe(server Server) error {
	r.isHTTPS = false
	return server.ListenAndServe()
}

// StartAndServeTls Start the HTTP server and listen for connections with Tls
func (r *Router) StartAndServeTLS(server Server, cert string, key string) error {
	r.isHTTPS = true
	return server.ListenAndServeTLS(cert, key)
}
