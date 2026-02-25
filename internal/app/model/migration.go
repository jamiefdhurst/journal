package model

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/jamiefdhurst/journal/internal/app"
    "github.com/jamiefdhurst/journal/pkg/database/rows"
)

const migrationTable = "migration"

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
func (ms *Migrations) CreateTable() error {
    _, err := ms.Container.Db.Exec("CREATE TABLE IF NOT EXISTS `" + migrationTable + "` (" +
        "`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
        "`name` VARCHAR(255) NOT NULL, " +
        "`applied` BOOLEAN NOT NULL DEFAULT 0" +
        ")")

    return err
}

// HasMigrationRun checks if a specific migration has been applied
func (ms *Migrations) HasMigrationRun(name string) bool {
    rows, err := ms.Container.Db.Query("SELECT * FROM `"+migrationTable+"` WHERE `name` = ? LIMIT 1", name)
    if err != nil {
        return false
    }

    migrations := ms.loadFromRows(rows)

    return len(migrations) > 0 && migrations[0].Applied
}

// RecordMigration marks a migration as applied
func (ms *Migrations) RecordMigration(name string) error {
    // Check if migration exists first
    rows, err := ms.Container.Db.Query("SELECT * FROM `"+migrationTable+"` WHERE `name` = ? LIMIT 1", name)
    if err != nil {
        return err
    }

    migrations := ms.loadFromRows(rows)

    var res sql.Result
    if len(migrations) == 0 {
        // Create new migration record
        res, err = ms.Container.Db.Exec("INSERT INTO `"+migrationTable+"` (`name`, `applied`) VALUES(?, ?)", name, true)
    } else {
        // Update existing migration record
        res, err = ms.Container.Db.Exec("UPDATE `"+migrationTable+"` SET `applied` = ? WHERE `id` = ?", true, migrations[0].ID)
    }

    if err != nil {
        return err
    }

    _, err = res.RowsAffected()
    return err
}

func (ms *Migrations) loadFromRows(rows rows.Rows) []Migration {
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
func (ms *Migrations) MigrateHTMLToMarkdown() error {
    const migrationName = "html_to_markdown"

    // Skip if already migrated
    if ms.HasMigrationRun(migrationName) {
        log.Println("HTML to Markdown migration already applied. Skipping...")
        return nil
    }

    log.Println("Running HTML to Markdown migration...")

    // Get all journal entries
    js := Journals{Container: ms.Container}
    journalEntries := js.FetchAll()

    log.Printf("Found %d journal entries to migrate\n", len(journalEntries))

    count := 0
    for _, journal := range journalEntries {
        // Convert HTML content to Markdown
        markdownContent := ms.Container.MarkdownProcessor.FromHTML(journal.Content)
        journal.Content = markdownContent

        // Save the entry with the new markdown content
        js.Save(journal)
        count++

        log.Printf("Migrated entry: %s (%d)\n", journal.Title, journal.ID)
    }

    log.Printf("Migration complete. Converted %d journal entries from HTML to Markdown.\n", count)

    // Record migration as completed
    err := ms.RecordMigration(migrationName)
    if err != nil {
        return fmt.Errorf("migration completed but failed to record status: %w", err)
    }

    return nil
}

// MigrateRandomSlugs fixes any journal entries that have the "random" slug
func (ms *Migrations) MigrateRandomSlugs() error {
    const migrationName = "random_slug_fix"

    // Skip if already migrated
    if ms.HasMigrationRun(migrationName) {
        log.Println("Random slug fix migration already applied. Skipping...")
        return nil
    }

    log.Println("Running random slug fix migration...")

    // Get the journal with the 'random' slug if it exists
    js := Journals{Container: ms.Container}
    randomJournal := js.FindBySlug("random")

    if randomJournal.ID == 0 {
        log.Println("No journal entry found with 'random' slug. Migration not needed.")
    } else {
        // Rename the slug to 'random-post'
        randomJournal.Slug = "random-post"
        js.Save(randomJournal)
        log.Printf("Migrated journal entry: %s (ID: %d) from 'random' to 'random-post'\n", randomJournal.Title, randomJournal.ID)
    }

    // Record migration as completed
    err := ms.RecordMigration(migrationName)
    if err != nil {
        return fmt.Errorf("migration completed but failed to record status: %w", err)
    }

    return nil
}

// MigrateAddTimestamps adds created_at and updated_at columns to the journal table
func (ms *Migrations) MigrateAddTimestamps() error {
    const migrationName = "add_timestamps"

    // Skip if already migrated
    if ms.HasMigrationRun(migrationName) {
        log.Println("Add timestamps migration already applied. Skipping...")
        return nil
    }

    log.Println("Running add timestamps migration...")

    // Add created_at column
    _, err := ms.Container.Db.Exec("ALTER TABLE `" + journalTable + "` ADD COLUMN `created_at` DATETIME DEFAULT NULL")
    if err != nil {
        return fmt.Errorf("failed to add created_at column: %w", err)
    }

    // Add updated_at column
    _, err = ms.Container.Db.Exec("ALTER TABLE `" + journalTable + "` ADD COLUMN `updated_at` DATETIME DEFAULT NULL")
    if err != nil {
        return fmt.Errorf("failed to add updated_at column: %w", err)
    }

    log.Println("Successfully added created_at and updated_at columns to journal table.")

    // Record migration as completed
    err = ms.RecordMigration(migrationName)
    if err != nil {
        return fmt.Errorf("migration completed but failed to record status: %w", err)
    }

    return nil
}
