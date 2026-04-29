package router

import (
    "net/http"
    "os"
    "path"
    "runtime"
    "testing"

    "github.com/jamiefdhurst/journal/test/mocks/controller"
    mockRouter "github.com/jamiefdhurst/journal/test/mocks/router"
)

type BlankContainer struct{}

func init() {
    _, filename, _, _ := runtime.Caller(0)
    dir := path.Join(path.Dir(filename), "../..")
    err := os.Chdir(dir)
    if err != nil {
        panic(err)
    }
}

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
    router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: errorController, StaticPaths: []string{"test"}}
    router.Get("/standard", standardController)
    router.Get("/param/[%s]", paramController)
    router.Get("/", indexController)

    // Serve static file
    staticRequest, _ := http.NewRequest("GET", "/style.css", nil)
    router.ServeHTTP(response, staticRequest)
    if errorController.HasRun {
        t.Errorf("Expected static file to have been served but error controller was run")
        errorController.HasRun = false
    }

    // Index
    indexRequest, _ := http.NewRequest("GET", "/", nil)
    router.ServeHTTP(response, indexRequest)
    if !indexController.HasRun || errorController.HasRun {
        t.Errorf("Expected index controller to have been served but error controller was run")
        errorController.HasRun = false
    }

    // Standard route
    standardRequest, _ := http.NewRequest("GET", "/standard", nil)
    router.ServeHTTP(response, standardRequest)
    if !standardController.HasRun || errorController.HasRun {
        t.Errorf("Expected standard controller to have been served but error controller was run")
        errorController.HasRun = false
    }

    // Param route
    paramRequest, _ := http.NewRequest("GET", "/param/test1", nil)
    router.ServeHTTP(response, paramRequest)
    if !paramController.HasRun || errorController.HasRun {
        t.Errorf("Expected param controller to have been served but error controller was run")
        errorController.HasRun = false
    }

    // Not found route
    notFoundRequest, _ := http.NewRequest("GET", "/random", nil)
    router.ServeHTTP(response, notFoundRequest)
    if !errorController.HasRun {
        t.Errorf("Expected error controller to have been served")
    }
}

func TestServeHTTP_HTTPHeaders(t *testing.T) {
    ctrl := &controller.MockController{}
    router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}
    server := &mockRouter.MockServer{}
    router.StartAndServe(server)

    response := controller.NewMockResponse()
    request, _ := http.NewRequest("GET", "/random", nil)
    router.ServeHTTP(response, request)

    csp := response.Headers.Get("Content-Security-Policy")
    xss := response.Headers.Get("X-XSS-Protection")
    sts := response.Headers.Get("Strict-Transport-Security")
    if csp == "" {
        t.Error("Expected CSP header to be present")
    }
    if xss == "" {
        t.Error("Expected XSS header to be present")
    }
    if sts != "" {
        t.Error("Expected STS header to not be present")
    }
}

func TestServeHTTP_HTTPSHeaders(t *testing.T) {
    ctrl := &controller.MockController{}
    router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}
    server := &mockRouter.MockServer{}
    router.StartAndServeTLS(server, "test/cert.pem", "test/key.pem")

    response := controller.NewMockResponse()
    request, _ := http.NewRequest("GET", "/random", nil)
    router.ServeHTTP(response, request)

    csp := response.Headers.Get("Content-Security-Policy")
    xss := response.Headers.Get("X-XSS-Protection")
    sts := response.Headers.Get("Strict-Transport-Security")
    if csp == "" {
        t.Error("Expected CSP header to be present")
    }
    if xss == "" {
        t.Error("Expected XSS header to be present")
    }
    if sts == "" {
        t.Error("Expected STS header to be present")
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

func TestStartAndServeTLS(t *testing.T) {
    ctrl := &controller.MockController{}
    router := Router{Container: &BlankContainer{}, Routes: []Route{}, ErrorController: ctrl}
    server := &mockRouter.MockServer{}
    router.StartAndServeTLS(server, "test/cert.pem", "test/key.pem")

    if !server.Listening {
        t.Errorf("Expected some routes to have been defined but none were found")
    }
}
