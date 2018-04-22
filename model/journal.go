package model

import (
	"database/sql"
	"regexp"
	"strings"
)

const journalTable = "journal"

// CreateJournalTable Create the actual table
func CreateJournalTable() error {
	_, err := db.Exec("CREATE TABLE `" + journalTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`slug` VARCHAR(255) NOT NULL, " +
		"`title` VARCHAR(255) NOT NULL, " +
		"`date` DATE NOT NULL, " +
		"`content` TEXT NOT NULL" +
		")")

	return err
}

// FetchAllJournals Get all journals
func FetchAllJournals() []Journal {
	rows, _ := db.Query("SELECT * FROM `" + journalTable + "` ORDER BY `date` DESC")
	journals := loadJournalsFromRows(rows)

	return journals
}

// FindJournalBySlug Find a journal by slug
func FindJournalBySlug(slug string) Journal {
	// Attempt to find the entry
	rows, _ := db.Query("SELECT * FROM `"+journalTable+"` WHERE `slug` = ? LIMIT 1", slug)
	journals := loadJournalsFromRows(rows)

	if len(journals) == 1 {
		return journals[0]
	}

	return Journal{}
}

// Slugify Utility to convert a string into a slug
func Slugify(s string) string {
	re := regexp.MustCompile("[\\W+]")

	return strings.ToLower(re.ReplaceAllString(s, "-"))
}

func loadJournalsFromRows(rows *sql.Rows) []Journal {
	var id int
	var slug string
	var title string
	var date string
	var content string

	defer rows.Close()
	journals := []Journal{}
	for rows.Next() {
		rows.Scan(&id, &slug, &title, &date, &content)
		journals = append(journals, Journal{id, slug, title, date, content})
	}

	return journals
}

// Journal model
type Journal struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
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

// Save Save a journal entry, either inserting it or updating it in the database
func (j *Journal) Save() {
	var stmt *sql.Stmt
	var res sql.Result

	// Convert content for saving
	j.Content = ExtractContentsAndSearchGiphyAPI(j.Content)

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
}
