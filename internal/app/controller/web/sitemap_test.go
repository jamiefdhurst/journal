package web

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/test/mocks/controller"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestSitemap_Run(t *testing.T) {
	db := &database.MockSqlite{}
	configuration := app.DefaultConfiguration()
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &Sitemap{}
	controller.DisableTracking()

	// Test showing all Journals in sitemap
	db.Rows = &database.MockJournal_MultipleRows{}
	request, _ := http.NewRequest("GET", "/sitemap.xml", strings.NewReader(""))
	request.Host = "example.com"
	controller.Init(container, []string{"", "0"}, request)
	controller.Run(response, request)

	if !strings.Contains(response.Content, "<loc>https://example.com/slug</loc>") || !strings.Contains(response.Content, "<loc>https://example.com/slug-2</loc>") {
		t.Error("Expected all journals to be rendered in sitemap")
	}
}
