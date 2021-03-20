package web

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
	"github.com/jamiefdhurst/journal/pkg/database"
)

// Index Handle displaying all blog entries
type Index struct {
	controller.Super
	Journals   []model.Journal
	Pages      []int
	Pagination database.PaginationInformation
	Saved      bool
}

// Run Index action
func (c *Index) Run(response http.ResponseWriter, request *http.Request) {

	container := c.Super.Container.(*app.Container)
	js := model.Journals{Container: container, Gs: model.GiphyAdapter(container)}

	pagination := database.PaginationQuery{Page: 1, ResultsPerPage: container.Configuration.ArticlesPerPage}
	query := request.URL.Query()
	if query["page"] != nil {
		page, err := strconv.Atoi(query["page"][0])
		if err == nil {
			pagination.Page = page
		}
	}

	c.Journals, c.Pagination = js.FetchPaginated(pagination)
	c.Saved = false
	if query["saved"] != nil {
		c.Saved = true
	}

	c.Pages = make([]int, c.Pagination.TotalPages)
	for i := range c.Pages {
		c.Pages[i] = i + 1
	}

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.tmpl",
		"./web/templates/index.tmpl")
	template.ExecuteTemplate(response, "layout", c)
}
