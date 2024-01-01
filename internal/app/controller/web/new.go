package web

import (
	"net/http"
	"text/template"
	"time"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// New Handle creating a new entry
type New struct {
	controller.Super
	Error   bool
	Journal model.Journal
}

// Run New action
func (c *New) Run(response http.ResponseWriter, request *http.Request) {
	container := c.Super.Container.(*app.Container)
	if !container.Configuration.EnableCreate {
		RunBadRequest(response, request, c.Super.Container)
		return
	}

	if request.Method == "GET" {
		c.Error = false
		flash := c.Session.GetFlash()
		if flash != nil && flash[0] == "error" {
			c.Error = true
		}

		c.Journal.Date = time.Now().Format("2006-01-02")

		c.SessionStore.Save(response)
		template, _ := template.ParseFiles(
			"./web/templates/_layout/default.html.tmpl",
			"./web/templates/new.html.tmpl",
			"./web/templates/_partial/form.html.tmpl")
		template.ExecuteTemplate(response, "layout", c)
	} else {
		if request.FormValue("title") == "" || request.FormValue("date") == "" || request.FormValue("content") == "" {
			c.Session.AddFlash("error")
			c.SessionStore.Save(response)
			http.Redirect(response, request, "/new", http.StatusFound)
			return
		}

		js := model.NewJournalStore(container, model.GiphyAdapter(container))
		journal := model.Journal{ID: 0, Slug: model.Slugify(request.FormValue("title")), Title: request.FormValue("title"), Date: request.FormValue("date"), Content: request.FormValue("content")}
		js.Save(journal)

		c.Session.AddFlash("saved")
		c.SessionStore.Save(response)
		http.Redirect(response, request, "/", http.StatusFound)
	}
}
