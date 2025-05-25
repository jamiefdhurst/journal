package model

import (
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
