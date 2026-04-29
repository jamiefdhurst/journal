package web

import (
    "net/http"
    "text/template"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/internal/app/model"
    "github.com/jamiefdhurst/journal/pkg/controller"
)

// Sitemap Generate an XML sitemap
type Sitemap struct {
    controller.Super
}

type sitemapTemplateData struct {
    Host     string
    Journals []model.Journal
}

// Run Sitemap
func (c *Sitemap) Run(response http.ResponseWriter, request *http.Request) {

    data := sitemapTemplateData{}
    container := c.Super.Container().(*app.Container)
    data.Host = request.Host
    js := model.Journals{Container: container}

    data.Journals = js.FetchAll()

    response.Header().Add("Content-type", "text/xml")
    template, _ := template.ParseFiles("./web/templates/sitemap.xml.tmpl")
    template.ExecuteTemplate(response, "content", data)
}
