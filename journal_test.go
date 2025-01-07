package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/pkg/adapter/giphy"
	"github.com/jamiefdhurst/journal/pkg/adapter/json"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/internal/app/router"
	"github.com/jamiefdhurst/journal/pkg/database"
	pkgrouter "github.com/jamiefdhurst/journal/pkg/router"
)

var (
	rtr    *pkgrouter.Router
	server *httptest.Server
)

func init() {
	container = &app.Container{Configuration: app.DefaultConfiguration()}
	rtr = router.NewRouter(container)
	server = httptest.NewServer(rtr)

	log.Println("Serving on " + server.URL)
}

func fixtures(t *testing.T) {
	adapter := giphy.Client{Client: &json.Client{}}
	db := &database.Sqlite{}
	if err := db.Connect("test/data/test.db"); err != nil {
		t.Error("Could not open test database for writing...")
	}

	// Setup container
	container.Db = db
	container.Giphy = adapter

	js := model.Journals{Container: container}
	db.Exec("DROP TABLE journal")
	js.CreateTable()

	// Set up data
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test", "Test", "<p>Test!</p>", "2018-01-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test-2", "Another Test", "<p>Test again!</p>", "2018-02-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test-3", "A Final Test", "<p>Test finally!</p>", "2018-03-01")
}

func TestConfig(t *testing.T) {
	os.Setenv("J_TITLE", "A Test Title")

	configuration := config()

	if configuration.Title != "A Test Title" {
		t.Error("Expected title to be set through environment")
	}
	if configuration.Port != "3000" {
		t.Errorf("Expected default port to be set, got %s", configuration.Port)
	}
	if configuration.Theme != "default" {
		t.Errorf("Expected default theme to be set, got %s", configuration.Theme)
	}
}

func TestLoadDatabase(t *testing.T) {
	container.Configuration.DatabasePath = "test/data/test.db"
	closeFunc := loadDatabase()
	closeFunc()
}

func TestLoadGiphy(t *testing.T) {
	existing := os.Getenv("J_GIPHY_API_KEY")
	os.Setenv("J_GIPHY_API_KEY", "foobar")
	loadGiphy()
	os.Setenv("J_GIPHY_API_KEY", existing)

	if container.Giphy == nil {
		t.Error("Expected Giphy adapter to be setup")
	}
}

func TestApiv1List(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/api/v1/post", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := `[{"id":3,"slug":"test-3","title":"A Final Test","date":"2018-03-01T00:00:00Z","content":"<p>Test finally!</p>"},{"id":2,"slug":"test-2","title":"Another Test","date":"2018-02-01T00:00:00Z","content":"<p>Test again!</p>"},{"id":1,"slug":"test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}]`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}

}

func TestApiV1Single(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/api/v1/post/test", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := `{"id":1,"slug":"test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Single_NotFound(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/api/v1/post/random", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 404 {
		t.Error("Expected 404 status code")
	}
}

func TestApiV1Create(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("PUT", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Test 4","date":"2018-06-01T00:00:00Z","content":"<p>Test 4!</p>"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 201 {
		t.Error("Expected 201 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := `{"id":4,"slug":"test-4","title":"Test 4","date":"2018-06-01T00:00:00Z","content":"<p>Test 4!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Create_InvalidRequest(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("PUT", server.URL+"/api/v1/post", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 400 {
		t.Error("Expected 400 status code")
	}
}

func TestApiV1Create_MissingData(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("PUT", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Test 4"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 400 {
		t.Error("Expected 400 status code")
	}
}

func TestApiV1Create_RepeatTitles(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("PUT", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Repeated","date":"2018-02-01T00:00:00Z","content":"<p>Repeated content test!</p>"}`))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if res.StatusCode != 201 {
		t.Error("Expected 201 status code")
	}

	request, _ = http.NewRequest("PUT", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Repeated","date":"2019-02-01T00:00:00Z","content":"<p>Repeated content test again!</p>"}`))
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if res.StatusCode != 201 {
		t.Error("Expected 201 status code")
	}

	request, _ = http.NewRequest("GET", server.URL+"/api/v1/post/repeated-1", nil)
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := `{"id":5,"slug":"repeated-1","title":"Repeated","date":"2019-02-01T00:00:00Z","content":"<p>Repeated content test again!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Update(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("POST", server.URL+"/api/v1/post/test", strings.NewReader(`{"title":"A different title"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := `{"id":1,"slug":"test","title":"A different title","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Update_NotFound(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("POST", server.URL+"/api/v1/post/random", strings.NewReader(`{"title":"A different title"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 404 {
		t.Error("Expected 404 status code")
	}
}

func TestApiV1Update_InvalidRequest(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("POST", server.URL+"/api/v1/post/test", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 400 {
		t.Error("Expected 400 status code")
	}
}
