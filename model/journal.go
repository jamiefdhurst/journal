package model

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

const table = "journal"

var db *sql.DB

// Journals Collection of Journals
type Journals struct {
	Journals []Journal
}

// Journal model
type Journal struct {
	ID      int
	Slug    string
	Title   string
	Date    string
	Content string
}

// Create Create a new journal entry, saving it as necessary
func (js *Journals) Create(id int, slug string, title string, date string, content string) Journal {
	j := Journal{id, slug, title, date, content}
	if j.ID == 0 && len(slug) > 0 {
		j = js.save(j)
	}

	js.Journals = append(js.Journals, j)

	return j
}

// FetchAll Get all journals
func (js *Journals) FetchAll() {
	rows, _ := db.Query("SELECT * FROM `" + table + "` ORDER BY `date` DESC")

	defer rows.Close()
	for rows.Next() {
		js.load(rows)
	}
}

// FindBySlug Find a journal by slug.
func (js *Journals) FindBySlug(s string) Journal {
	// Attempt to find the entry
	rows, _ := db.Query("SELECT * FROM `"+table+"` WHERE `slug` = ?", strings.Replace(s, "/", "", 1))

	defer rows.Close()
	for rows.Next() {
		js.load(rows)
	}

	if len(js.Journals) == 1 {
		return js.Journals[0]
	}

	return Journal{}
}

func (js *Journals) load(rows *sql.Rows) {
	var id int
	var slug string
	var title string
	var date string
	var content string

	rows.Scan(&id, &slug, &title, &date, &content)
	js.Create(id, slug, title, date, content)
}

func (js *Journals) save(j Journal) Journal {
	var stmt *sql.Stmt
	var res sql.Result
	if j.ID == 0 {
		stmt, _ = db.Prepare("INSERT INTO `" + table + "` (`slug`, `title`, `date`, `content`) VALUES(?,?,?,?)")
		res, _ = stmt.Exec(j.Slug, j.Title, j.Date, j.Content)
	} else {
		stmt, _ = db.Prepare("UPDATE `" + table + "` SET `slug` = ?, `title` = ?, `date` = ?, `content` = ? WHERE `id` = ?")
		res, _ = stmt.Exec(j.Slug, j.Title, j.Date, j.Content, j.ID)
	}

	// Store insert ID
	if j.ID == 0 {
		id, _ := res.LastInsertId()
		j.ID = int(id)
	}

	return j
}

// GetDate Get the friendly date for the Journal
func (j Journal) GetDate() string {
	re := regexp.MustCompile("\\d{4}\\-\\d{2}\\-\\d{2}")
	date := re.FindString(j.Date)
	dateArr := strings.Split(date, "-")
	for i := 0; i < len(dateArr)/2; i++ {
		k := len(dateArr) - i - 1
		dateArr[i], dateArr[k] = dateArr[k], dateArr[i]
	}

	return strings.Join(dateArr, "/")
}

// JournalCreateTable Create the actual table
func JournalCreateTable() error {
	_, err := db.Exec("CREATE TABLE `journal` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`slug` VARCHAR(255) NOT NULL, " +
		"`title` VARCHAR(255) NOT NULL, " +
		"`date` DATE NOT NULL, " +
		"`content` TEXT NOT NULL" +
		")")

	return err
}

// Slugify Utility to convert a string into a slug
func Slugify(s string) string {
	re := regexp.MustCompile("[\\W+]")

	return strings.ToLower(re.ReplaceAllString(s, "-"))
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./data/journal.db")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
