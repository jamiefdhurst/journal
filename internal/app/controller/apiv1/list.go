package apiv1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
	"github.com/jamiefdhurst/journal/pkg/database"
)

type listResponse struct {
	Links      database.PaginationLinks       `json:"links"`
	Pagination database.PaginationInformation `json:"pagination"`
	Posts      []journalToJSON                `json:"posts"`
}

// List Display all blog entries as JSON
type List struct {
	controller.Super
}

func ListData(request *http.Request, js model.Journals) ([]model.Journal, database.PaginationInformation) {
	paginationQuery := database.PaginationQuery{Page: 1, ResultsPerPage: js.Container.Configuration.ArticlesPerPage}
	query := request.URL.Query()
	if query["page"] != nil {
		page, err := strconv.Atoi(query["page"][0])
		if err == nil {
			paginationQuery.Page = page
		}
	}

	return js.FetchPaginated(paginationQuery)
}

// Run List action
func (c *List) Run(response http.ResponseWriter, request *http.Request) {
	container := c.Super.Container().(*app.Container)
	js := model.Journals{Container: container}

	journals, paginationInfo := ListData(request, js)
	jsonResponse := listResponse{database.LinksPagination("/api/v1/post", paginationInfo), paginationInfo, MapJournalsToJSON(journals)}

	response.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	encoder.Encode(jsonResponse)
}
