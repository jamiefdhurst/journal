package apiv1

import "github.com/jamiefdhurst/journal/internal/app/model"

type journalFromJSON struct {
    Title   string
    Date    string
    Content string
}

type journalToJSON struct {
    URL       string  `json:"url"`
    Title     string  `json:"title"`
    Date      string  `json:"date"`
    Content   string  `json:"content"`
    CreatedAt *string `json:"created_at,omitempty"`
    UpdatedAt *string `json:"updated_at,omitempty"`
}

func MapJournalToJSON(journal model.Journal) journalToJSON {
    result := journalToJSON{
        URL:     "/api/v1/post/" + journal.Slug,
        Title:   journal.Title,
        Date:    journal.GetEditableDate(),
        Content: journal.Content,
    }

    // Format timestamps in ISO 8601 format if they exist
    if journal.CreatedAt != nil {
        createdAtStr := journal.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
        result.CreatedAt = &createdAtStr
    }
    if journal.UpdatedAt != nil {
        updatedAtStr := journal.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
        result.UpdatedAt = &updatedAtStr
    }

    return result
}

func MapJournalsToJSON(journals []model.Journal) []journalToJSON {
    result := make([]journalToJSON, len(journals))
    for i, j := range journals {
        result[i] = MapJournalToJSON(j)
    }
    return result
}
