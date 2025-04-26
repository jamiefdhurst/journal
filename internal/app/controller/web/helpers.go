package web

import (
	"net/http"
	"text/template"

	"github.com/jamiefdhurst/journal/internal/app/model"
	"github.com/jamiefdhurst/journal/pkg/controller"
)

type formTemplateData struct {
	Error   bool
	Journal model.Journal
}

func RenderFromSession(c controller.Controller, data formTemplateData, templateFile string, response http.ResponseWriter) {
	data.Error = false
	flash := c.Session().GetFlash()
	if flash != nil && flash[0] == "error" {
		data.Error = true

		formData := c.Session().Get("form_data")
		if formData != nil {
			formMap := formData.(map[string]string)
			data.Journal.Title = formMap["title"]
			data.Journal.Date = formMap["date"]
			data.Journal.Content = formMap["content"]

			c.Session().Delete("form_data")
		}
	}

	c.SaveSession(response)
	responseTemplate, _ := template.ParseFiles(
		"./web/templates/_layout/default.html.tmpl",
		"./web/templates/"+templateFile+".html.tmpl",
		"./web/templates/_partial/form.html.tmpl")
	responseTemplate.ExecuteTemplate(response, "layout", data)
}
