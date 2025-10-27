package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/controller/apiv1"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
	"github.com/jamiefdhurst/journal/pkg/database"
)

// Index Handle displaying all blog entries
type Index struct {
	controller.Super
}

type indexTemplateData struct {
	Container  interface{}
	Excerpt    func(model.Journal) string
	Journals   []model.Journal
	Pages      []int
	Pagination database.PaginationDisplay
	Saved      bool
}

// Run Index action
func (c *Index) Run(response http.ResponseWriter, request *http.Request) {

	data := indexTemplateData{}

	container := c.Super.Container().(*app.Container)
	data.Container = container
	js := model.Journals{Container: container}

	var paginationInfo database.PaginationInformation
	data.Journals, paginationInfo = apiv1.ListData(request, js)
	data.Pagination = database.DisplayPagination(paginationInfo)
	data.Saved = false
	flash := c.Session().GetFlash()
	if flash != nil && flash[0] == "saved" {
		data.Saved = true
	}

	data.Pages = make([]int, database.PAGINATION_MAX_PAGES)
	i := 0
	for p := data.Pagination.FirstPage; p <= data.Pagination.LastPage; p++ {
		data.Pages[i] = p
		i++
	}

	data.Excerpt = func(j model.Journal) string {
		return j.GetHTMLExcerpt(container.Configuration.ExcerptWords)
	}

	c.SaveSession(response)
	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/index.html.tmpl")
	template.ExecuteTemplate(response, "layout", data)
}
