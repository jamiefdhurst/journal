package web

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestCalendarRun(t *testing.T) {
	db := &database.MockSqlite{}
	configuration := app.DefaultConfiguration()
	container := &app.Container{Configuration: configuration, Db: db}
	response := controller.NewMockResponse()
	controller := &Calendar{}
	controller.DisableTracking()

	// Test showing current year/month (only prev nav)
	today := time.Now()
	firstOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	daysInMonth := firstOfMonth.AddDate(0, 1, -1).Day()
	db.EnableMultiMode()
	db.AppendResult(&database.MockJournal_SingleRow{})
	for d := 2; d <= daysInMonth; d++ {
		db.AppendResult(&database.MockRowsEmpty{})
	}
	db.AppendResult(&database.MockJournal_SingleRow{})
	request, _ := http.NewRequest("GET", "/calendar", strings.NewReader(""))
	controller.Init(container, []string{}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected title of journal to be shown in calendar")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-year\"") {
		t.Error("Expected previous year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-month\"") {
		t.Error("Expected previous month link to be shown")
	}
	if strings.Contains(response.Content, "class=\"next next-year\"") {
		t.Error("Expected next year link to be missing")
	}
	if strings.Contains(response.Content, "class=\"next next-month\"") {
		t.Error("Expected next month link to be missing")
	}

	// Test showing beginning (only next nav)
	response.Reset()
	db.EnableMultiMode()
	db.AppendResult(&database.MockJournal_SingleRow{})
	for d := 2; d <= 28; d++ {
		db.AppendResult(&database.MockRowsEmpty{})
	}
	db.AppendResult(&database.MockJournal_SingleRow{})
	request, _ = http.NewRequest("GET", "/calendar/2018/feb", strings.NewReader(""))
	controller.Init(container, []string{"", "2018", "feb"}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected title of journal to be shown in calendar")
	}
	if !strings.Contains(response.Content, "<h2>2018</h2>") || !strings.Contains(response.Content, "<h2>February</h2") {
		t.Error("Expected correct year and month to be shown")
	}
	if strings.Contains(response.Content, "class=\"prev prev-year\"") {
		t.Error("Expected previous year link to be missing")
	}
	if strings.Contains(response.Content, "class=\"prev prev-month\"") {
		t.Error("Expected previous month link to be missing")
	}
	if !strings.Contains(response.Content, "class=\"next next-year\"") {
		t.Error("Expected next year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"next next-month\"") {
		t.Error("Expected next month link to be shown")
	}

	// Test showing middle (both prev and next nav)
	response.Reset()
	lastYear := today.Year() - 1
	nextMonth := strings.ToLower(today.AddDate(0, 1, 0).Format("Jan"))
	firstOfMonth = time.Date(lastYear, today.AddDate(0, 1, 0).Month(), 1, 0, 0, 0, 0, today.Location())
	daysInMonth = firstOfMonth.AddDate(0, 1, -1).Day()
	db.EnableMultiMode()
	db.AppendResult(&database.MockJournal_SingleRow{})
	for d := 2; d <= daysInMonth; d++ {
		db.AppendResult(&database.MockRowsEmpty{})
	}
	db.AppendResult(&database.MockJournal_SingleRow{})
	request, _ = http.NewRequest("GET", "/calendar/"+strconv.Itoa(lastYear)+"/"+nextMonth, strings.NewReader(""))
	controller.Init(container, []string{"", strconv.Itoa(lastYear), nextMonth}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected title of journal to be shown in calendar")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-year\"") {
		t.Error("Expected previous year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-month\"") {
		t.Error("Expected previous month link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"next next-year\"") {
		t.Error("Expected next year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"next next-month\"") {
		t.Error("Expected next month link to be shown")
	}

	// Test year only
	response.Reset()
	db.EnableMultiMode()
	db.AppendResult(&database.MockJournal_SingleRow{})
	for d := 2; d <= 31; d++ {
		db.AppendResult(&database.MockRowsEmpty{})
	}
	db.AppendResult(&database.MockJournal_SingleRow{})
	request, _ = http.NewRequest("GET", "/calendar/2019", strings.NewReader(""))
	controller.Init(container, []string{"", "2019"}, request)
	controller.Run(response, request)
	if !strings.Contains(response.Content, "Title") {
		t.Error("Expected title of journal to be shown in calendar")
	}
	if !strings.Contains(response.Content, "<h2>2019</h2>") || !strings.Contains(response.Content, "<h2>January</h2") {
		t.Error("Expected correct year and month to be shown")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-year\"") {
		t.Error("Expected previous year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"prev prev-month\"") {
		t.Error("Expected previous month link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"next next-year\"") {
		t.Error("Expected next year link to be shown")
	}
	if !strings.Contains(response.Content, "class=\"next next-month\"") {
		t.Error("Expected next month link to be shown")
	}
}
