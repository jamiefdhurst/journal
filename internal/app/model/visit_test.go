package model

import (
    "testing"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestVisits_CreateTable(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Db: db}
    visits := Visits{Container: container}

    err := visits.CreateTable()

    if err != nil {
        t.Errorf("Expected no error creating table, got: %s", err)
    }
}

func TestVisits_FindByDateAndURL(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Db: db}
    visits := Visits{Container: container}

    db.Rows = &database.MockVisit_SingleRow{}
    visit := visits.FindByDateAndURL("2023-01-01", "/test")

    if visit.ID != 1 {
        t.Errorf("Expected visit ID to be 1, got %d", visit.ID)
    }
    if visit.URL != "/test" {
        t.Errorf("Expected visit URL to be /test, got %s", visit.URL)
    }
    if visit.Hits != 5 {
        t.Errorf("Expected visit hits to be 5, got %d", visit.Hits)
    }

    // Test with no visit found
    db.Rows = &database.MockRowsEmpty{}
    emptyVisit := visits.FindByDateAndURL("2023-01-01", "/nonexistent")

    if emptyVisit.ID != 0 {
        t.Errorf("Expected empty visit ID to be 0, got %d", emptyVisit.ID)
    }
}

func TestVisits_RecordVisit(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Db: db}
    visits := Visits{Container: container}

    db.Rows = &database.MockRowsEmpty{} // No existing visit
    db.Result = &database.MockResult{}

    err := visits.RecordVisit("/new-page")

    if err != nil {
        t.Errorf("Expected no error recording new visit, got: %s", err)
    }

    db.Rows = &database.MockVisit_SingleRow{} // Existing visit
    db.Result = &database.MockResult{}

    err = visits.RecordVisit("/test")

    if err != nil {
        t.Errorf("Expected no error updating existing visit, got: %s", err)
    }
}

func TestVisits_GetDailyStats(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Db: db}
    visits := Visits{Container: container}

    // Test with mock data
    db.Rows = &database.MockVisitStats_DailyRows{}

    dailyStats := visits.GetDailyStats(14)

    if len(dailyStats) != 2 {
        t.Errorf("Expected 2 daily stats, got %d", len(dailyStats))
    }

    if len(dailyStats) > 0 {
        if dailyStats[0].Date != "2023-12-25" {
            t.Errorf("Expected first date to be 2023-12-25, got %s", dailyStats[0].Date)
        }
        if dailyStats[0].Total != 57 {
            t.Errorf("Expected first total to be 57, got %d", dailyStats[0].Total)
        }
    }
}

func TestVisits_GetMonthlyStats(t *testing.T) {
    db := &database.MockSqlite{}
    container := &app.Container{Db: db}
    visits := Visits{Container: container}

    // Test with mock data
    db.Rows = &database.MockVisitStats_MonthlyRows{}

    monthlyStats := visits.GetMonthlyStats()

    if len(monthlyStats) != 2 {
        t.Errorf("Expected 2 monthly stats, got %d", len(monthlyStats))
    }

    if len(monthlyStats) > 0 {
        if monthlyStats[0].Month != "2023-12" {
            t.Errorf("Expected first month to be 2023-12, got %s", monthlyStats[0].Month)
        }
        if monthlyStats[0].Total != 1700 {
            t.Errorf("Expected first total to be 1700, got %d", monthlyStats[0].Total)
        }
    }
}

func TestDailyVisit_GetFriendlyDate(t *testing.T) {
    visit := DailyVisit{Date: "2023-12-25"}
    
    friendly := visit.GetFriendlyDate()
    expected := "Monday December 25, 2023"
    
    if friendly != expected {
        t.Errorf("Expected friendly date to be %s, got %s", expected, friendly)
    }
    
    // Test with invalid date
    invalidVisit := DailyVisit{Date: "invalid-date"}
    if invalidVisit.GetFriendlyDate() != "invalid-date" {
        t.Error("Expected invalid date to return original string")
    }
}

func TestMonthlyVisit_GetFriendlyMonth(t *testing.T) {
    visit := MonthlyVisit{Month: "2023-12"}
    
    friendly := visit.GetFriendlyMonth()
    expected := "December 2023"
    
    if friendly != expected {
        t.Errorf("Expected friendly month to be %s, got %s", expected, friendly)
    }
    
    // Test with invalid month
    invalidVisit := MonthlyVisit{Month: "invalid-month"}
    if invalidVisit.GetFriendlyMonth() != "invalid-month" {
        t.Error("Expected invalid month to return original string")
    }
}
