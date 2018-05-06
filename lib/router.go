package lib

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/controller/apiv1"
	"github.com/jamiefdhurst/journal/controller/web"
	"github.com/jamiefdhurst/journal/model"
)

// Server Common interface for HTTP
type Server interface {
	ListenAndServe() error
}

// Route A route contains a method (GET), URI, and a controller
type Route struct {
	method     string
	regexURI   string
	controller controller.Controller
}

// Router A router contains routes and links back to the application and implements the ServeHTTP interface
type Router struct {
	Db              model.Database
	Routes          []Route
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

	// Attempt to serve a file first
	if request.URL.Path != "/" {
		file := "public" + request.URL.Path
		_, err := os.Stat(file)
		if !os.IsNotExist(err) {
			http.ServeFile(response, request, file)
			return
		}
	}

	// Go through each route and attempt to match
	for _, route := range r.Routes {
		matched, _ := regexp.MatchString(route.regexURI, request.URL.Path)
		if matched && (request.Method == route.method || (request.Method == "" && route.method == "GET")) {
			re := regexp.MustCompile(route.regexURI)
			route.controller.Init(r.Db, re.FindStringSubmatch(request.URL.Path))
			route.controller.Run(response, request)
			return
		}
	}

	r.ErrorController.Run(response, request)
}

// StartAndServe Start the HTTP server and listen for connections
func (r *Router) StartAndServe(server Server) error {
	r.Get("/new", &web.New{})
	r.Post("/new", &web.New{})
	r.Get("/api/v1/post", &apiv1.List{})
	r.Post("/api/v1/post", &apiv1.Create{})
	r.Get("/api/v1/post/[%s]", &apiv1.Single{})
	r.Put("/api/v1/post/[%s]", &apiv1.Update{})
	r.Get("/[%s]/edit", &web.Edit{})
	r.Post("/[%s]/edit", &web.Edit{})
	r.Get("/[%s]", &web.View{})
	r.Get("/", &web.Index{})

	return server.ListenAndServe()
}
