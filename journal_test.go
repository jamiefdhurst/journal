package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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
	rtr = router.NewRouter(nil)
	server = httptest.NewServer(rtr)

	log.Println("Serving on " + server.URL)
}

func fixtures(t *testing.T) {
	container := &app.Container{}
	adapter := giphy.Client{Client: &json.Client{}}
	db := &database.Sqlite{}
	if err := db.Connect("test/data/test.db"); err != nil {
		t.Error("Could not open test database for writing...")
	}

	// Setup container
	container.Db = db
	container.Giphy = adapter
	rtr.Container = container

	js := model.Journals{Container: container}
	db.Exec("DROP TABLE journal")
	js.CreateTable()

	// Clear database
	db.Exec("DELETE FROM journal")

	// Set up data
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test", "Test", "<p>Test!</p>", "2018-01-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test-2", "Another Test", "<p>Test again!</p>", "2018-02-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?)", "test-3", "A Final Test", "<p>Test finally!</p>", "2018-03-01")
}

func TestApiv1List(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("GET", server.URL+"/api/v1/post", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	expected := `[{"id":3,"slug":"test-3","title":"A Final Test","date":"2018-03-01T00:00:00Z","content":"<p>Test finally!</p>"},{"id":2,"slug":"test-2","title":"Another Test","date":"2018-02-01T00:00:00Z","content":"<p>Test again!</p>"},{"id":1,"slug":"test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}]`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}

}

func TestApiV1Single(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("GET", server.URL+"/api/v1/post/test", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	expected := `{"id":1,"slug":"test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Single_NotFound(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("GET", server.URL+"/api/v1/post/random", nil)

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

	request, err := http.NewRequest("POST", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Test 4","date":"2018-06-01T00:00:00Z","content":"<p>Test 4!</p>"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	expected := `{"id":4,"slug":"test-4","title":"Test 4","date":"2018-06-01T00:00:00Z","content":"<p>Test 4!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Create_InvalidRequest(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("POST", server.URL+"/api/v1/post", nil)

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

	request, err := http.NewRequest("POST", server.URL+"/api/v1/post", strings.NewReader(`{"title":"Test 4"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 400 {
		t.Error("Expected 400 status code")
	}
}

func TestApiV1Update(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("PUT", server.URL+"/api/v1/post/test", strings.NewReader(`{"title":"A different title"}`))

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	expected := `{"id":1,"slug":"test","title":"A different title","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Update_NotFound(t *testing.T) {
	fixtures(t)

	request, err := http.NewRequest("PUT", server.URL+"/api/v1/post/random", strings.NewReader(`{"title":"A different title"}`))

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

	request, err := http.NewRequest("PUT", server.URL+"/api/v1/post/test", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 400 {
		t.Error("Expected 400 status code")
	}
}
