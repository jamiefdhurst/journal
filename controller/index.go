package controller

import (
	"journal/lib"
	"journal/model"
	"net/http"
	"text/template"
)

// Index Handle displaying all blog entries
type Index struct {
	lib.Controller
}

// IndexData Data for index
type IndexData struct {
	Journals []model.Journal
}

// Run Index
func (c *Index) Run(w http.ResponseWriter, r *http.Request) {
	rows, err := c.Db.Query("SELECT * FROM `journal` ORDER BY `date` DESC")
	lib.CheckErr(err)

	js := IndexData{}

	for rows.Next() {
		j := model.Journal{}
		err := rows.Scan(&j.ID, &j.Slug, &j.Title, &j.Date, &j.Content)
		lib.CheckErr(err)

		js.Journals = append(js.Journals, j)
	}

	t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/index.tmpl")
	t.ExecuteTemplate(w, "header", nil)
	t.ExecuteTemplate(w, "content", js)
	t.ExecuteTemplate(w, "footer", nil)
	t.Execute(w, nil)
}
