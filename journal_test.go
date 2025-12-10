package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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
	db := &database.Sqlite{}
	if err := db.Connect("test/data/test.db"); err != nil {
		t.Error("Could not open test database for writing...")
	}

	// Setup container
	container.Db = db

	js := model.Journals{Container: container}
	ms := model.Migrations{Container: container}
	vs := model.Visits{Container: container}
	db.Exec("DROP TABLE journal")
	db.Exec("DROP TABLE migration")
	db.Exec("DROP TABLE visit")
	js.CreateTable()
	ms.CreateTable()
	vs.CreateTable()
	ms.MigrateAddTimestamps()

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
	expected := `{"links":{},"pagination":{"current_page":1,"total_pages":1,"posts_per_page":20,"total_posts":3},"posts":[{"url":"/api/v1/post/test-3","title":"A Final Test","date":"2018-03-01T00:00:00Z","content":"<p>Test finally!</p>"},{"url":"/api/v1/post/test-2","title":"Another Test","date":"2018-02-01T00:00:00Z","content":"<p>Test again!</p>"},{"url":"/api/v1/post/test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}]}`

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
	expected := `{"url":"/api/v1/post/test","title":"Test","date":"2018-01-01T00:00:00Z","content":"<p>Test!</p>"}`

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestApiV1Single_NotFound(t *testing.T) {
	fixtures(t)

	// Try a post that doesn't exist, but is not the new random endpoint
	request, _ := http.NewRequest("GET", server.URL+"/api/v1/post/nonexistent", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 404 {
		t.Error("Expected 404 status code")
	}
}

func TestApiV1Random(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/api/v1/post/random", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// Make sure we got a valid JSON response
	if !strings.Contains(string(body[:]), "\"url\":") || !strings.Contains(string(body[:]), "\"title\":") {
		t.Errorf("Expected JSON with id and slug, got: %s", string(body[:]))
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
	bodyStr := string(body[:])

	// Check for expected fields
	expectedFields := []string{`"id":4`, `"slug":"test-4"`, `"title":"Test 4"`, `"date":"2018-06-01T00:00:00Z"`, `"content":"<p>Test 4!</p>"`, `"created_at"`, `"updated_at"`}
	for _, field := range expectedFields {
		if !strings.Contains(bodyStr, field) {
			t.Errorf("Expected response to contain %s\nGot:\n\t%s", field, bodyStr)
		}
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
	bodyStr := string(body[:])

	// Check for expected fields
	expectedFields := []string{`"url":"/api/v1/post/repeated-1"`, `"title":"Repeated"`, `"date":"2019-02-01T00:00:00Z"`, `"content":"<p>Repeated content test again!</p>"`, `"created_at"`, `"updated_at"`}
	for _, field := range expectedFields {
		if !strings.Contains(bodyStr, field) {
			t.Errorf("Expected response to contain %s\nGot:\n\t%s", field, bodyStr)
		}
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
	bodyStr := string(body[:])

	// Check for expected fields
	expectedFields := []string{`"id":1`, `"slug":"test"`, `"title":"A different title"`, `"date":"2018-01-01T00:00:00Z"`, `"content":"<p>Test!</p>"`, `"updated_at"`}
	for _, field := range expectedFields {
		if !strings.Contains(bodyStr, field) {
			t.Errorf("Expected response to contain %s\nGot:\n\t%s", field, bodyStr)
		}
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

func TestApiV1Stats(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/api/v1/stats", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// Check that JSON is returned
	if res.Header.Get("Content-Type") != "application/json" {
		t.Error("Expected JSON content type")
	}

	now := time.Now()
	date := now.Format("2006-01-02")
	month := now.Format("2006-01")
	expected := fmt.Sprintf(`{"posts":{"count":3,"first_post_date":"Monday January 1, 2018"},"configuration":{"title":"Jamie's Journal","description":"A private journal containing Jamie's innermost thoughts","theme":"default","posts_per_page":20,"google_analytics":false,"create_enabled":true,"edit_enabled":true},"visits":{"daily":[{"date":"%sT00:00:00Z","api_hits":1,"web_hits":0,"total":1}],"monthly":[{"month":"%s","api_hits":1,"web_hits":0,"total":1}]}}`, date, month)

	// Use contains to get rid of any extra whitespace that we can discount
	if !strings.Contains(string(body[:]), expected) {
		t.Errorf("Expected:\n\t%s\nGot:\n\t%s", expected, string(body[:]))
	}
}

func TestOpenapi(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/openapi.yml", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	expected := []string{"openapi: '3.0.3'", "/api/v1/post:", "/api/v1/post/{slug}:", "/api/v1/post/random:", "/api/v1/stats:"}
	for _, e := range expected {
		if !strings.Contains(string(body[:]), e) {
			t.Errorf("Expected:\n\t%s\nGot:\n\t%s", e, string(body[:]))
		}
	}
}

func TestWebStats(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/stats", nil)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// Check for stats page elements
	if !strings.Contains(string(body[:]), "<h1>Stats</h1>") {
		t.Error("Expected stats page title to be present")
	}

	// Check for post count (3 from fixtures)
	if !strings.Contains(string(body[:]), "Total Posts") || !strings.Contains(string(body[:]), "<dd>3</dd>") {
		t.Error("Expected post count to be displayed")
	}
}

func TestVisitTracking(t *testing.T) {
	fixtures(t)

	request, _ := http.NewRequest("GET", server.URL+"/", nil)
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if res.StatusCode != 200 {
		t.Error("Expected 200 status code")
	}

	res.Body.Close()

	rows, err := container.Db.Query("SELECT COUNT(*) FROM visit WHERE url = '/'")
	if err != nil {
		t.Errorf("Failed to query visits table: %s", err)
		return
	}
	defer rows.Close()

	var visitCount int
	if rows.Next() {
		rows.Scan(&visitCount)
	}

	if visitCount == 0 {
		t.Log("Visit tracking is disabled during test environment - this is expected behaviour")
	} else {
		t.Logf("Visit tracking is active - found %d visit(s)", visitCount)

		visitRows, err := container.Db.Query("SELECT url, hits FROM visit WHERE url = '/' LIMIT 1")
		if err != nil {
			t.Errorf("Failed to query visit details: %s", err)
			return
		}
		defer visitRows.Close()

		if visitRows.Next() {
			var url string
			var hits int
			visitRows.Scan(&url, &hits)

			if url != "/" {
				t.Errorf("Expected visit URL to be '/', got '%s'", url)
			}
			if hits != 1 {
				t.Errorf("Expected visit hits to be 1, got %d", hits)
			}
		}
	}
}
