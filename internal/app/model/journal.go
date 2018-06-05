package model

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"

	"github.com/jamiefdhurst/journal/pkg/database"
	"github.com/jamiefdhurst/journal/pkg/database/rows"
)

const journalTable = "journal"

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
	if len(dateArr) != 3 {
		return ""
	}
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

// Journals Common database resource link for Journal actions
type Journals struct {
	Db database.Database
	Gs GiphysExtractor
}

// CreateTable Create the actual table
func (js *Journals) CreateTable() error {
	_, err := js.Db.Exec("CREATE TABLE `" + journalTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`slug` VARCHAR(255) NOT NULL, " +
		"`title` VARCHAR(255) NOT NULL, " +
		"`date` DATE NOT NULL, " +
		"`content` TEXT NOT NULL" +
		")")

	return err
}

// FetchAll Get all journals
func (js *Journals) FetchAll() []Journal {
	rows, err := js.Db.Query("SELECT * FROM `" + journalTable + "` ORDER BY `date` DESC")
	if err != nil {
		return []Journal{}
	}

	return js.loadFromRows(rows)
}

// FindBySlug Find a journal by slug
func (js *Journals) FindBySlug(slug string) Journal {
	// Attempt to find the entry
	rows, err := js.Db.Query("SELECT * FROM `"+journalTable+"` WHERE `slug` = ? LIMIT 1", slug)
	if err != nil {
		return Journal{}
	}
	journals := js.loadFromRows(rows)

	if len(journals) == 1 {
		return journals[0]
	}

	return Journal{}
}

// Save Save a journal entry, either inserting it or updating it in the database
func (js *Journals) Save(j Journal) Journal {
	var res sql.Result

	// Convert content for saving
	j.Content = js.Gs.ExtractContentsAndSearchAPI(j.Content)

	if j.ID == 0 {
		res, _ = js.Db.Exec("INSERT INTO `"+journalTable+"` (`slug`, `title`, `date`, `content`) VALUES(?,?,?,?)", j.Slug, j.Title, j.Date, j.Content)
	} else {
		res, _ = js.Db.Exec("UPDATE `"+journalTable+"` SET `slug` = ?, `title` = ?, `date` = ?, `content` = ? WHERE `id` = ?", j.Slug, j.Title, j.Date, j.Content, strconv.Itoa(j.ID))
	}

	// Store insert ID
	if j.ID == 0 {
		id, _ := res.LastInsertId()
		j.ID = int(id)
	}

	return j
}

func (js Journals) loadFromRows(rows rows.Rows) []Journal {
	defer rows.Close()
	journals := []Journal{}
	for rows.Next() {
		j := Journal{}
		rows.Scan(&j.ID, &j.Slug, &j.Title, &j.Date, &j.Content)
		journals = append(journals, j)
	}

	return journals
}

// Slugify Utility to convert a string into a slug
func Slugify(s string) string {
	re := regexp.MustCompile("[\\W+]")

	return strings.ToLower(re.ReplaceAllString(s, "-"))
}
