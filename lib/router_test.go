package lib

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jamiefdhurst/journal/controller"
	"github.com/jamiefdhurst/journal/model"
)

func TestGet(t *testing.T) {
	controller := &controller.MockController{}
	database := &model.MockDatabase{}
	router := Router{Db: database, Routes: []Route{}, ErrorController: controller}

	// Test normal route
	router.Get("/testing", controller)
	if router.Routes[0].controller != controller || router.Routes[0].method != "GET" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test paramterised route
	router.Get("/[%s]/[%d]/[%a]", controller)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestPost(t *testing.T) {
	controller := &controller.MockController{}
	database := &model.MockDatabase{}
	router := Router{Db: database, Routes: []Route{}, ErrorController: controller}

	// Test normal route
	router.Post("/testing", controller)
	if router.Routes[0].controller != controller || router.Routes[0].method != "POST" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test paramterised route
	router.Post("/[%s]/[%d]/[%a]", controller)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestPut(t *testing.T) {
	controller := &controller.MockController{}
	database := &model.MockDatabase{}
	router := Router{Db: database, Routes: []Route{}, ErrorController: controller}

	// Test normal route
	router.Put("/testing", controller)
	if router.Routes[0].controller != controller || router.Routes[0].method != "PUT" || router.Routes[0].regexURI != "^\\/testing$" {
		t.Errorf("GET Route added was not as expected")
	}

	// Test paramterised route
	router.Put("/[%s]/[%d]/[%a]", controller)
	if router.Routes[1].regexURI != "^\\/([\\w\\-]+)\\/(\\d+)\\/(.+?)$" {
		t.Errorf("GET Route added was not as expected")
	}
}

func TestServeHTTP(t *testing.T) {
	errorController := &controller.MockController{}
	indexController := &controller.MockController{}
	standardController := &controller.MockController{}
	paramController := &controller.MockController{}
	database := &model.MockDatabase{}
	response := controller.NewMockResponse()
	router := Router{Db: database, Routes: []Route{}, ErrorController: errorController}
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
	}

	// Index
	indexURL := &url.URL{Path: "/"}
	indexRequest := &http.Request{URL: indexURL, Method: "GET"}
	router.ServeHTTP(response, indexRequest)
	if !indexController.HasRun || errorController.HasRun {
		t.Errorf("Expected index controller to have been served but error controller was run")
	}

	// Standard route
	standardURL := &url.URL{Path: "/standard"}
	standardRequest := &http.Request{URL: standardURL, Method: "GET"}
	router.ServeHTTP(response, standardRequest)
	if !standardController.HasRun || errorController.HasRun {
		t.Errorf("Expected standard controller to have been served but error controller was run")
	}

	// Parameterised route
	paramURL := &url.URL{Path: "/param/test1"}
	paramRequest := &http.Request{URL: paramURL, Method: "GET"}
	router.ServeHTTP(response, paramRequest)
	if !paramController.HasRun || errorController.HasRun {
		t.Errorf("Expected param controller to have been served but error controller was run")
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
	controller := &controller.MockController{}
	database := &model.MockDatabase{}
	router := Router{Db: database, Routes: []Route{}, ErrorController: controller}
	router.StartAndServe(MockServer{})

	if len(router.Routes) < 1 {
		t.Errorf("Expected some routes to have been defined but none were found")
	}
}
