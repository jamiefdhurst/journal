package main

import (
	"net/http/httptest"
	"testing"

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
}

func fixtures(t *testing.T) {
	db := &database.Sqlite{}
	if err := db.Connect("test/data/test.db"); err != nil {
		t.Error("Could not open test database for writing...")
	}

	defer db.Close()

	// Clear database
	db.Exec("DELETE FROM journal")

	// Set up data
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?", "test", "Test", "<p>Test!</p>", "2018-01-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?", "test-2", "Another Test", "<p>Test again!</p>", "2018-02-01")
	db.Exec("INSERT INTO journal (slug, title, content, date) VALUES (?, ?, ?, ?", "test-3", "A Final Test", "<p>Test finally!</p>", "2018-03-01")

	// Use database
	rtr.Db = db
}

func TestApiv1List(t *testing.T) {
	fixtures(t)
}
