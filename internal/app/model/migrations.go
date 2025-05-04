package model

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/pkg/database/rows"
)

const migrationsTable = "migrations"

// Migration stores a record of migrations that have been applied
type Migration struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Applied bool   `json:"applied"`
}

// Migrations manages database migrations
type Migrations struct {
	Container *app.Container
}

// CreateTable initializes the migrations table
func (m *Migrations) CreateTable() error {
	_, err := m.Container.Db.Exec("CREATE TABLE IF NOT EXISTS `" + migrationsTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`name` VARCHAR(255) NOT NULL, " +
		"`applied` BOOLEAN NOT NULL DEFAULT 0" +
		")")

	return err
}

// HasMigrationRun checks if a specific migration has been applied
func (m *Migrations) HasMigrationRun(name string) bool {
	rows, err := m.Container.Db.Query("SELECT * FROM `"+migrationsTable+"` WHERE `name` = ? LIMIT 1", name)
	if err != nil {
		return false
	}

	migrations := m.loadFromRows(rows)

	return len(migrations) > 0 && migrations[0].Applied
}

// RecordMigration marks a migration as applied
func (m *Migrations) RecordMigration(name string) error {
	// Check if migration exists first
	rows, err := m.Container.Db.Query("SELECT * FROM `"+migrationsTable+"` WHERE `name` = ? LIMIT 1", name)
	if err != nil {
		return err
	}

	migrations := m.loadFromRows(rows)

	var res sql.Result
	if len(migrations) == 0 {
		// Create new migration record
		res, err = m.Container.Db.Exec("INSERT INTO `"+migrationsTable+"` (`name`, `applied`) VALUES(?, ?)", name, true)
	} else {
		// Update existing migration record
		res, err = m.Container.Db.Exec("UPDATE `"+migrationsTable+"` SET `applied` = ? WHERE `id` = ?", true, migrations[0].ID)
	}

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	return err
}

func (m *Migrations) loadFromRows(rows rows.Rows) []Migration {
	defer rows.Close()
	migrations := []Migration{}
	for rows.Next() {
		migration := Migration{}
		rows.Scan(&migration.ID, &migration.Name, &migration.Applied)
		migrations = append(migrations, migration)
	}

	return migrations
}

// MigrateHTMLToMarkdown converts all journal entries from HTML to Markdown
func (m *Migrations) MigrateHTMLToMarkdown() error {
	const migrationName = "html_to_markdown"

	// Skip if already migrated
	if m.HasMigrationRun(migrationName) {
		log.Println("HTML to Markdown migration already applied. Skipping...")
		return nil
	}

	log.Println("Running HTML to Markdown migration...")

	// Get all journal entries
	js := Journals{Container: m.Container}
	journalEntries := js.FetchAll()

	log.Printf("Found %d journal entries to migrate\n", len(journalEntries))

	count := 0
	for _, journal := range journalEntries {
		// Convert HTML content to Markdown
		markdownContent := m.Container.MarkdownProcessor.FromHTML(journal.Content)
		journal.Content = markdownContent

		// Save the entry with the new markdown content
		js.Save(journal)
		count++

		log.Printf("Migrated entry: %s (%d)\n", journal.Title, journal.ID)
	}

	log.Printf("Migration complete. Converted %d journal entries from HTML to Markdown.\n", count)

	// Record migration as completed
	err := m.RecordMigration(migrationName)
	if err != nil {
		return fmt.Errorf("migration completed but failed to record status: %w", err)
	}

	return nil
}