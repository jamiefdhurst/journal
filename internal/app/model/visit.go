package model

import (
	"strconv"
	"time"

	"github.com/jamiefdhurst/journal/internal/app"
)

const visitTable = "visit"

// Visit stores a record of daily visits for a given endpoint/web address
type Visit struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	URL  string `json:"url"`
	Hits int    `json:"hits"`
}

// Visits manages tracking API hits
type Visits struct {
	Container *app.Container
}

// CreateTable initializes the visits table
func (vs *Visits) CreateTable() error {
	_, err := vs.Container.Db.Exec("CREATE TABLE IF NOT EXISTS `" + visitTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`date` DATE NOT NULL, " +
		"`url` VARCHAR(255) NOT NULL, " +
		"`hits` INTEGER UNSIGNED NOT NULL DEFAULT 0" +
		")")

	return err
}

// FindByDateAndURL finds a visit record for a specific date and URL
func (vs *Visits) FindByDateAndURL(date, url string) Visit {
	visit := Visit{}
	rows, err := vs.Container.Db.Query("SELECT * FROM `"+visitTable+"` WHERE `date` = ? AND `url` = ? LIMIT 1", date, url)
	if err != nil {
		return visit
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&visit.ID, &visit.Date, &visit.URL, &visit.Hits)
		return visit
	}

	return Visit{}
}

// RecordVisit records or updates a visit for the given URL and current date
func (vs *Visits) RecordVisit(url string) error {
	today := time.Now().Format("2006-01-02")

	existingVisit := vs.FindByDateAndURL(today, url)
	var err error
	if existingVisit.ID > 0 {
		_, err = vs.Container.Db.Exec("UPDATE `"+visitTable+"` SET `hits` = `hits` + 1 WHERE `id` = ?", strconv.Itoa(existingVisit.ID))
	} else {
		_, err = vs.Container.Db.Exec("INSERT INTO `"+visitTable+"` (`date`, `url`, `hits`) VALUES (?, ?, 1)", today, url)
	}

	return err
}
