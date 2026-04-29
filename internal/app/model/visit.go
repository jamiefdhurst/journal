package model

import (
    "regexp"
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

// DailyVisit represents daily visit statistics
type DailyVisit struct {
    Date    string `json:"date"`
    APIHits int    `json:"api_hits"`
    WebHits int    `json:"web_hits"`
    Total   int    `json:"total"`
}

// GetFriendlyDate returns a human-readable date format
func (d DailyVisit) GetFriendlyDate() string {
    re := regexp.MustCompile(`\d{4}\-\d{2}\-\d{2}`)
    date := re.FindString(d.Date)
    timeObj, err := time.Parse("2006-01-02", date)
    if err != nil {
        return d.Date
    }
    return timeObj.Format("Monday January 2, 2006")
}

// MonthlyVisit represents monthly visit statistics
type MonthlyVisit struct {
    Month   string `json:"month"`
    APIHits int    `json:"api_hits"`
    WebHits int    `json:"web_hits"`
    Total   int    `json:"total"`
}

// GetFriendlyMonth returns a human-readable month format
func (m MonthlyVisit) GetFriendlyMonth() string {
    timeObj, err := time.Parse("2006-01", m.Month)
    if err != nil {
        return m.Month
    }
    return timeObj.Format("January 2006")
}

// GetDailyStats returns visit statistics for the last N days
func (vs *Visits) GetDailyStats(days int) []DailyVisit {
    // Calculate the date N days ago
    startDate := time.Now().AddDate(0, 0, -days+1).Format("2006-01-02")

    query := `
        SELECT 
            DATE(date),
            COALESCE(SUM(CASE WHEN url LIKE '/api/%' THEN hits ELSE 0 END), 0) as api_hits,
            COALESCE(SUM(CASE WHEN url NOT LIKE '/api/%' THEN hits ELSE 0 END), 0) as web_hits,
            COALESCE(SUM(hits), 0) as total
        FROM ` + visitTable + `
        WHERE date >= ?
        GROUP BY date
        ORDER BY date DESC
    `

    rows, err := vs.Container.Db.Query(query, startDate)
    if err != nil {
        return []DailyVisit{}
    }
    defer rows.Close()

    var dailyStats []DailyVisit
    for rows.Next() {
        var stat DailyVisit
        rows.Scan(&stat.Date, &stat.APIHits, &stat.WebHits, &stat.Total)
        dailyStats = append(dailyStats, stat)
    }

    return dailyStats
}

// GetMonthlyStats returns visit statistics aggregated by month
func (vs *Visits) GetMonthlyStats() []MonthlyVisit {
    query := `
        SELECT 
            strftime('%Y-%m', date) as month,
            COALESCE(SUM(CASE WHEN url LIKE '/api/%' THEN hits ELSE 0 END), 0) as api_hits,
            COALESCE(SUM(CASE WHEN url NOT LIKE '/api/%' THEN hits ELSE 0 END), 0) as web_hits,
            COALESCE(SUM(hits), 0) as total
        FROM ` + visitTable + `
        GROUP BY strftime('%Y-%m', date)
        ORDER BY month DESC
    `

    rows, err := vs.Container.Db.Query(query)
    if err != nil {
        return []MonthlyVisit{}
    }
    defer rows.Close()

    var monthlyStats []MonthlyVisit
    for rows.Next() {
        var stat MonthlyVisit
        rows.Scan(&stat.Month, &stat.APIHits, &stat.WebHits, &stat.Total)
        monthlyStats = append(monthlyStats, stat)
    }

    return monthlyStats
}
