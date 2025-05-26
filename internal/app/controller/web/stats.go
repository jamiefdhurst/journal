package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Stats Handle displaying journal stats and configuration
type Stats struct {
	controller.Super
}

type statsTemplateData struct {
	Container       *app.Container
	PostCount       int
	FirstPostDate   string
	TitleSet        bool
	DescriptionSet  bool
	ThemeSet        bool
	ArticlesPerPage int
	GACodeSet       bool
	CreateEnabled   bool
	EditEnabled     bool
	DailyVisits     []model.DailyVisit
	MonthlyVisits   []model.MonthlyVisit
}

// Run Stats action
func (c *Stats) Run(response http.ResponseWriter, request *http.Request) {
	data := statsTemplateData{}

	container := c.Super.Container().(*app.Container)
	data.Container = container

	js := model.Journals{Container: container}
	allJournals := js.FetchAll()
	data.PostCount = len(allJournals)

	if data.PostCount > 0 {
		firstPost := allJournals[data.PostCount-1]
		data.FirstPostDate = firstPost.GetDate()
	} else {
		data.FirstPostDate = "No posts yet"
	}

	// Settings status
	defaultConfig := app.DefaultConfiguration()
	data.TitleSet = container.Configuration.Title != defaultConfig.Title
	data.DescriptionSet = container.Configuration.Description != defaultConfig.Description
	data.ThemeSet = container.Configuration.Theme != defaultConfig.Theme
	data.ArticlesPerPage = container.Configuration.ArticlesPerPage
	data.GACodeSet = container.Configuration.GoogleAnalyticsCode != ""
	data.CreateEnabled = container.Configuration.EnableCreate
	data.EditEnabled = container.Configuration.EnableEdit

	vs := model.Visits{Container: container}
	data.DailyVisits = vs.GetDailyStats(14)
	data.MonthlyVisits = vs.GetMonthlyStats()

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/stats.html.tmpl")
	template.ExecuteTemplate(response, "layout", data)
}
