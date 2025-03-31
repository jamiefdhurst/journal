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
	Journals   []model.Journal
	Pages      []int
	Pagination database.PaginationDisplay
	Saved      bool
}

// Run Index action
func (c *Index) Run(response http.ResponseWriter, request *http.Request) {

	container := c.Super.Container.(*app.Container)
	js := model.Journals{Container: container}

	var paginationInfo database.PaginationInformation
	c.Journals, paginationInfo = apiv1.ListData(request, js)
	c.Pagination = database.DisplayPagination(paginationInfo)
	c.Saved = false
	flash := c.Session.GetFlash()
	if flash != nil && flash[0] == "saved" {
		c.Saved = true
	}

	c.Pages = make([]int, database.PAGINATION_MAX_PAGES)
	i := 0
	for p := c.Pagination.FirstPage; p <= c.Pagination.LastPage; p++ {
		c.Pages[i] = p
		i++
	}

	c.SessionStore.Save(response)
	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/index.html.tmpl")
	template.ExecuteTemplate(response, "layout", c)
}
