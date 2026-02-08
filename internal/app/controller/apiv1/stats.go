package apiv1

import (
	"encoding/json"
	"net/http"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Stats Provide statistics about the journal system
type Stats struct {
	controller.Super
}

type statsJSON struct {
	Posts         statsPostsJSON  `json:"posts"`
	Configuration statsConfigJSON `json:"configuration"`
	Visits        statsVisitsJSON `json:"visits"`
}

type statsVisitsJSON struct {
	Daily   []model.DailyVisit   `json:"daily"`
	Monthly []model.MonthlyVisit `json:"monthly"`
}

type statsPostsJSON struct {
	Count         int    `json:"count"`
	FirstPostDate string `json:"first_post_date,omitempty"`
}

type statsConfigJSON struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Theme           string `json:"theme"`
	PostsPerPage    int    `json:"posts_per_page"`
	GoogleAnalytics bool   `json:"google_analytics"`
	CreateEnabled   bool   `json:"create_enabled"`
	EditEnabled     bool   `json:"edit_enabled"`
}

// Run Stats action
func (c *Stats) Run(response http.ResponseWriter, request *http.Request) {
	stats := statsJSON{}

	container := c.Super.Container().(*app.Container)

	js := model.Journals{Container: container}
	allJournals := js.FetchAll()
	stats.Posts.Count = len(allJournals)

	if stats.Posts.Count > 0 {
		firstPost := allJournals[stats.Posts.Count-1]
		stats.Posts.FirstPostDate = firstPost.GetEditableDate()
	}

	stats.Configuration.Title = container.Configuration.Title
	stats.Configuration.Description = container.Configuration.Description
	stats.Configuration.Theme = container.Configuration.Theme
	stats.Configuration.PostsPerPage = container.Configuration.PostsPerPage
	stats.Configuration.GoogleAnalytics = container.Configuration.GoogleAnalyticsCode != ""
	stats.Configuration.CreateEnabled = container.Configuration.EnableCreate
	stats.Configuration.EditEnabled = container.Configuration.EnableEdit

	vs := model.Visits{Container: container}
	stats.Visits.Daily = vs.GetDailyStats(14)
	stats.Visits.Monthly = vs.GetMonthlyStats()

	// Send JSON response
	response.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	encoder.SetEscapeHTML(false)
	encoder.Encode(stats)
}
