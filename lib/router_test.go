package lib

import (
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/jamiefdhurst/journal/model"
)

type FakeController struct {
	HasRun bool
}

func (f *FakeController) Init(db model.Database, params []string) {}

func (f *FakeController) Run(response http.ResponseWriter, request *http.Request) {
	f.HasRun = true
}

type FakeDatabase struct{}

func (f *FakeDatabase) Close() {}

func (f *FakeDatabase) Connect() error {
	return nil
}

func (f *FakeDatabase) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (f *FakeDatabase) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

type FakeResponse struct{}

func (f FakeResponse) Header() http.Header {
	return map[string][]string{}
}

func (f FakeResponse) Write([]byte) (int, error) {
	return 0, nil
}

func (f FakeResponse) WriteHeader(statusCode int) {}

type FakeServer struct{}

func (f FakeServer) ListenAndServe() error {
	return nil
}

func TestGet(t *testing.T) {
	controller := &FakeController{}
	database := &FakeDatabase{}
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
	controller := &FakeController{}
	database := &FakeDatabase{}
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
	controller := &FakeController{}
	database := &FakeDatabase{}
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
	errorController := &FakeController{}
	indexController := &FakeController{}
	standardController := &FakeController{}
	paramController := &FakeController{}
	database := &FakeDatabase{}
	response := FakeResponse{}
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
	if errorController.HasRun == true {
		t.Errorf("Expected static file to have been served but error controller was run")
	}

	// Index
	indexURL := &url.URL{Path: "/"}
	indexRequest := &http.Request{URL: indexURL, Method: "GET"}
	router.ServeHTTP(response, indexRequest)
	if indexController.HasRun != true || errorController.HasRun == true {
		t.Errorf("Expected index controller to have been served but error controller was run")
	}

	// Standard route
	standardURL := &url.URL{Path: "/standard"}
	standardRequest := &http.Request{URL: standardURL, Method: "GET"}
	router.ServeHTTP(response, standardRequest)
	if standardController.HasRun != true || errorController.HasRun == true {
		t.Errorf("Expected standard controller to have been served but error controller was run")
	}

	// Parameterised route
	paramURL := &url.URL{Path: "/param/test1"}
	paramRequest := &http.Request{URL: paramURL, Method: "GET"}
	router.ServeHTTP(response, paramRequest)
	if paramController.HasRun != true || errorController.HasRun == true {
		t.Errorf("Expected param controller to have been served but error controller was run")
	}

	// Not found route
	notFoundURL := &url.URL{Path: "/random"}
	notFoundRequest := &http.Request{URL: notFoundURL, Method: "GET"}
	router.ServeHTTP(response, notFoundRequest)
	if errorController.HasRun != true {
		t.Errorf("Expected error controller to have been served")
	}
}

func TestStartAndServe(t *testing.T) {
	controller := &FakeController{}
	database := &FakeDatabase{}
	router := Router{Db: database, Routes: []Route{}, ErrorController: controller}
	router.StartAndServe(FakeServer{})

	if len(router.Routes) < 1 {
		t.Errorf("Expected some routes to have been defined but none were found")
	}
}
