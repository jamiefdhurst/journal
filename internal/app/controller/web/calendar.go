package web

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Calendar Handle displaying a calendar with blog entries for given days
type Calendar struct {
	controller.Super
}

type day struct {
	Date    time.Time
	IsEmpty bool
}

type calendarTemplateData struct {
	Container    interface{}
	Days         map[int][]model.Journal
	Weeks        [][]day
	CurrentDate  time.Time
	PrevYear     int
	PrevYearUrl  string
	NextYear     int
	NextYearUrl  string
	PrevMonth    string
	PrevMonthUrl string
	NextMonth    string
	NextMonthUrl string
}

// Run Calendar action
func (c *Calendar) Run(response http.ResponseWriter, request *http.Request) {

	data := calendarTemplateData{}

	container := c.Super.Container().(*app.Container)
	data.Container = container
	js := model.Journals{Container: container}

	// Load date from parameters if available (either 2006/jan or 2006)
	date := time.Now()
	var err error
	if len(c.Params()) == 3 {
		date, err = time.Parse("2006 Jan 02", c.Params()[1]+" "+cases.Title(language.English, cases.NoLower).String(c.Params()[2])+" 25")
	} else if len(c.Params()) == 2 {
		date, err = time.Parse("2006-01-02", c.Params()[1]+"-01-01")
	}
	if err != nil {
		log.Print(err)
		RunBadRequest(response, request, container)
		return
	}

	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	startWeekday := int(firstOfMonth.Weekday())

	// Find number of days in month
	nextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := nextMonth.AddDate(0, 0, -1)
	daysInMonth := lastOfMonth.Day()

	data.Days = map[int][]model.Journal{}
	data.Weeks = [][]day{}
	week := []day{}

	// Fill in blanks before first day
	for range startWeekday {
		week = append(week, day{IsEmpty: true})
	}

	// Fill in actual days
	for d := 1; d <= daysInMonth; d++ {
		thisDate := time.Date(date.Year(), date.Month(), d, 0, 0, 0, 0, date.Location())
		data.Days[d] = js.FetchByDate(thisDate.Format("2006-01-02"))
		week = append(week, day{
			Date:    thisDate,
			IsEmpty: false,
		})

		// If Saturday, start a new week
		if thisDate.Weekday() == time.Saturday {
			data.Weeks = append(data.Weeks, week)
			week = []day{}
		}
	}

	// Fill in blanks after last day
	if len(week) > 0 {
		for len(week) < 7 {
			week = append(week, day{IsEmpty: true})
		}
		data.Weeks = append(data.Weeks, week)
	}

	// Load prev/next year and month
	firstEntry := js.FindNext(0)
	firstEntryDate, _ := time.Parse("2006-01-02", firstEntry.GetEditableDate())
	if date.Year() < time.Now().Year() {
		data.NextYear = date.Year() + 1
		data.NextYearUrl = strconv.Itoa(data.NextYear) + "/" + strings.ToLower(date.Format("Jan"))
		if date.AddDate(1, 0, 0).After(time.Now()) {
			data.NextYearUrl = strconv.Itoa(data.NextYear) + "/" + strings.ToLower(time.Now().Format("Jan"))
		}
	}
	if date.Year() > firstEntryDate.Year() {
		data.PrevYear = date.Year() - 1
		data.PrevYearUrl = strconv.Itoa(data.PrevYear) + "/" + strings.ToLower(date.Format("Jan"))
		if date.AddDate(-1, 0, 0).Before(firstEntryDate) {
			data.PrevYearUrl = strconv.Itoa(data.PrevYear) + "/" + strings.ToLower(firstEntryDate.Format("Jan"))
		}
	}
	if date.Year() < time.Now().Year() || date.Month() < time.Now().Month() {
		data.NextMonth = date.AddDate(0, 0, 31).Format("January")
		data.NextMonthUrl = strings.ToLower(date.AddDate(0, 0, 31).Format("2006/Jan"))
	}
	if date.Year() > firstEntryDate.Year() || date.Month() > firstEntryDate.Month() {
		data.PrevMonth = date.AddDate(0, 0, -31).Format("January")
		data.PrevMonthUrl = strings.ToLower(date.AddDate(0, 0, -31).Format("2006/Jan"))
	}
	data.CurrentDate = date

	template, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/calendar.html.tmpl")
	template.ExecuteTemplate(response, "layout", data)
}
