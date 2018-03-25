package model

import (
	"database/sql"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

const journalTable = "journal"

// Journals Collection of Journals
type Journals struct {
	Journals []Journal
}

// Journal model
type Journal struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
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
	rows, _ := db.Query("SELECT * FROM `" + journalTable + "` ORDER BY `date` DESC")

	defer rows.Close()
	for rows.Next() {
		js.load(rows)
	}
}

// FindBySlug Find a journal by slug.
func (js *Journals) FindBySlug(s string) Journal {
	// Attempt to find the entry
	rows, _ := db.Query("SELECT * FROM `"+journalTable+"` WHERE `slug` = ?", s)

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

	// Convert content for saving
	j.Content = ConvertSearchesToIDs(j.Content)

	if j.ID == 0 {
		stmt, _ = db.Prepare("INSERT INTO `" + journalTable + "` (`slug`, `title`, `date`, `content`) VALUES(?,?,?,?)")
		res, _ = stmt.Exec(j.Slug, j.Title, j.Date, j.Content)
	} else {
		stmt, _ = db.Prepare("UPDATE `" + journalTable + "` SET `slug` = ?, `title` = ?, `date` = ?, `content` = ? WHERE `id` = ?")
		res, _ = stmt.Exec(j.Slug, j.Title, j.Date, j.Content, j.ID)
	}

	// Store insert ID
	if j.ID == 0 {
		id, _ := res.LastInsertId()
		j.ID = int(id)
	}

	return j
}

// Update Save an existing journal entry's changes
func (js *Journals) Update(j Journal) Journal {
	return js.save(j)
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

// GetEditableDate Get the date string for editing
func (j Journal) GetEditableDate() string {
	re := regexp.MustCompile("\\d{4}\\-\\d{2}\\-\\d{2}")
	return re.FindString(j.Date)
}

// JournalCreateTable Create the actual table
func JournalCreateTable() error {
	_, err := db.Exec("CREATE TABLE `" + journalTable + "` (" +
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
