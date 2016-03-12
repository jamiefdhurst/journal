package controller

import (
	"journal/lib"
	"journal/model"
	"net/http"
	"strings"
	"text/template"
)

// View Handle displaying individual entry
type View struct {
	lib.Controller
}

// ViewData Data for view
type ViewData struct {
	Journal model.Journal
}

// Run View
func (c *View) Run(w http.ResponseWriter, r *http.Request) {

	// Attempt to find the entry
	rows, err := c.Db.Query("SELECT * FROM `journal` WHERE `slug` = ?", strings.Replace(c.Params[0], "/", "", 1))
	lib.CheckErr(err)

	v := ViewData{}

	for rows.Next() {
		v.Journal = model.Journal{}
		err := rows.Scan(&v.Journal.ID, &v.Journal.Slug, &v.Journal.Title, &v.Journal.Date, &v.Journal.Content)
		lib.CheckErr(err)
	}

	if v.Journal.ID == 0 {
		http.NotFound(w, r)
	} else {
		t, _ := template.ParseFiles("./src/journal/views/_layout/header.tmpl", "./src/journal/views/_layout/footer.tmpl", "./src/journal/views/view.tmpl")
		t.ExecuteTemplate(w, "header", nil)
		t.ExecuteTemplate(w, "content", v)
		t.ExecuteTemplate(w, "footer", nil)
		t.Execute(w, nil)
	}
}
