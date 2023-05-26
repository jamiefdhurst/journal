package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

// Sitemap Generate an XML sitemap
type Sitemap struct {
	controller.Super
	Journals []model.Journal
}

// Run Sitemap
func (c *Sitemap) Run(response http.ResponseWriter, request *http.Request) {

	container := c.Super.Container.(*app.Container)
	js := model.Journals{Container: container, Gs: model.GiphyAdapter(container)}

	c.Journals = js.FetchAll()

	log.Println(c.Host)

	response.Header().Add("Content-type", "text/xml")
	template, _ := template.ParseFiles("./web/templates/sitemap.xml.tmpl")
	template.ExecuteTemplate(response, "content", c)
}
