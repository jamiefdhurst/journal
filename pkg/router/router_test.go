package router

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jamiefdhurst/journal/test/mocks/controller"
	mockRouter "github.com/jamiefdhurst/journal/test/mocks/router"
)

type BlankContainer struct{}

func TestGet(t *testing.T) {
	ctrl := &controller.MockController{}
	router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}

	// Test normal route
	router.Get("/testing", ctrl)
	if router.Routes[0].controller != ctrl || router.Routes[0].method != "GET" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test param route
	router.Get("/[%s]/[%d]/[%a]", ctrl)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestPost(t *testing.T) {
	ctrl := &controller.MockController{}
	router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}

	// Test normal route
	router.Post("/testing", ctrl)
	if router.Routes[0].controller != ctrl || router.Routes[0].method != "POST" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test param route
	router.Post("/[%s]/[%d]/[%a]", ctrl)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestPut(t *testing.T) {
	ctrl := &controller.MockController{}
	router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}

	// Test normal route
	router.Put("/testing", ctrl)
	if router.Routes[0].controller != ctrl || router.Routes[0].method != "PUT" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test param route
	router.Put("/[%s]/[%d]/[%a]", ctrl)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestServeHTTP(t *testing.T) {
	errorController := &controller.MockController{}
	indexController := &controller.MockController{}
	standardController := &controller.MockController{}
	paramController := &controller.MockController{}
	response := controller.NewMockResponse()
	router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: errorController}
	router.Get("/standard", standardController)
	router.Get("/param/[%s]", paramController)
	router.Get("/", indexController)

	// Set CWD
	os.Chdir(os.Getenv("GOPATH") + "/src/github.com/jamiefdhurst/journal")

	// Serve static file
	staticURL := &url.URL{Path: "/css/default.min.css"}
	staticRequest := &http.Request{URL: staticURL, Method: "GET"}
	router.ServeHTTP(response, staticRequest)
	if errorController.HasRun {
		t.Errorf("Expected static file to have been served but error controller was run")
		errorController.HasRun = false
	}

	// Index
	indexURL := &url.URL{Path: "/"}
	indexRequest := &http.Request{URL: indexURL, Method: "GET"}
	router.ServeHTTP(response, indexRequest)
	if !indexController.HasRun || errorController.HasRun {
		t.Errorf("Expected index controller to have been served but error controller was run")
		errorController.HasRun = false
	}

	// Standard route
	standardURL := &url.URL{Path: "/standard"}
	standardRequest := &http.Request{URL: standardURL, Method: "GET"}
	router.ServeHTTP(response, standardRequest)
	if !standardController.HasRun || errorController.HasRun {
		t.Errorf("Expected standard controller to have been served but error controller was run")
		errorController.HasRun = false
	}

	// Param route
	paramURL := &url.URL{Path: "/param/test1"}
	paramRequest := &http.Request{URL: paramURL, Method: "GET"}
	router.ServeHTTP(response, paramRequest)
	if !paramController.HasRun || errorController.HasRun {
		t.Errorf("Expected param controller to have been served but error controller was run")
		errorController.HasRun = false
	}

	// Not found route
	notFoundURL := &url.URL{Path: "/random"}
	notFoundRequest := &http.Request{URL: notFoundURL, Method: "GET"}
	router.ServeHTTP(response, notFoundRequest)
	if !errorController.HasRun {
		t.Errorf("Expected error controller to have been served")
	}
}

func TestStartAndServe(t *testing.T) {
	ctrl := &controller.MockController{}
	router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}
	server := &mockRouter.MockServer{}
	router.StartAndServe(server)

	if !server.Listening {
		t.Errorf("Expected some routes to have been defined but none were found")
	}
}
