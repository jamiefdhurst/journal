package model

import (
	"regexp"
	"strings"
	"time"
)

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
	re := regexp.MustCompile(`\d{4}\-\d{2}\-\d{2}`)
	date := re.FindString(j.Date)
	timeObj, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	return timeObj.Format("Monday January 2, 2006")
}

// GetEditableDate Get the date string for editing
func (j Journal) GetEditableDate() string {
	re := regexp.MustCompile(`\d{4}\-\d{2}\-\d{2}`)
	return re.FindString(j.Date)
}

// GetExcerpt returns a small extract of the entry
func (j Journal) GetExcerpt() string {
	strip := regexp.MustCompile("\b+")
	text := strings.ReplaceAll(j.Content, "<p>", "")
	text = strings.ReplaceAll(text, "</p>", " ")
	text = strip.ReplaceAllString(text, " ")
	words := strings.Split(text, " ")

	if len(words) > 50 {
		return strings.Join(words[:50], " ") + "..."
	}
	return strings.TrimSuffix(strings.Join(words, " "), " ")
}

// Slugify Utility to convert a string into a slug
func Slugify(s string) string {
	re := regexp.MustCompile(`[\W+]`)

	return strings.ToLower(re.ReplaceAllString(s, "-"))
}
