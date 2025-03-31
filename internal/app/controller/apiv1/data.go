package apiv1

import "github.com/jamiefdhurst/journal/internal/app/model"

type journalFromJSON struct {
	Title   string
	Date    string
	Content string
}

type journalToJSON struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
}

func MapJournalToJSON(journal model.Journal) journalToJSON {
	return journalToJSON{
		"/api/v1/post/" + journal.Slug,
		journal.Title,
		journal.Date,
		journal.Content,
	}
}

func MapJournalsToJSON(journals []model.Journal) []journalToJSON {
	result := make([]journalToJSON, len(journals))
	for i, j := range journals {
		result[i] = MapJournalToJSON(j)
	}
	return result
}
