package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite 3 driver
)

const giphyTable = "giphy"

// Giphys Collection of Giphy API key entries
type Giphys struct {
	Giphys []Giphy
}

// Giphy model
type Giphy struct {
	ID     int    `json:"id"`
	APIKey string `json:"apiKey"`
}

// GetKey Get the current API key
func (gs *Giphys) GetKey() string {
	// Attempt to find the entry
	rows, _ := db.Query("SELECT * FROM `" + giphyTable + "`")

	defer rows.Close()
	for rows.Next() {
		gs.load(rows)
	}

	if len(gs.Giphys) == 1 {
		return gs.Giphys[0].APIKey
	}

	return ""
}

func (gs *Giphys) load(rows *sql.Rows) {
	var id int
	var apiKey string

	rows.Scan(&id, &apiKey)
	g := Giphy{id, apiKey}
	gs.Giphys = append(gs.Giphys, g)
}

func (gs *Giphys) save(g Giphy) Giphy {
	var stmt *sql.Stmt
	stmt, _ = db.Prepare("REPLACE INTO `" + giphyTable + "` (`id`, `apiKey`) VALUES(?,?)")
	stmt.Exec(g.ID, g.APIKey)

	return g
}

// Update Update/save the API key
func (gs *Giphys) Update(apiKey string) Giphy {
	g := Giphy{1, apiKey}
	g = gs.save(g)

	gs.Giphys = append(gs.Giphys, g)

	return g
}

// GiphyCreateTable Create the actual table
func GiphyCreateTable() error {
	_, err := db.Exec("CREATE TABLE `" + giphyTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`apiKey` VARCHAR(64) NOT NULL" +
		")")

	return err
}
