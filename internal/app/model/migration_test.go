package model

import (
	"testing"

	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/pkg/markdown"
	"github.com/jamiefdhurst/journal/test/mocks/database"
)

func TestMigrations_CreateTable(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	migrations := Migrations{Container: container}
	migrations.CreateTable()
	if db.Queries != 1 {
		t.Errorf("Expected 1 query to have been run")
	}
}

func TestMigrations_HasMigrationRun(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	migrations := Migrations{Container: container}

	// Test error case
	db.ErrorMode = true
	if migrations.HasMigrationRun("test_migration") {
		t.Error("Should return false when database has an error")
	}

	// Test migration not found
	db.ErrorMode = false
	db.Rows = &database.MockRowsEmpty{}
	if migrations.HasMigrationRun("test_migration") {
		t.Error("Should return false when migration doesn't exist")
	}

	// Test migration found and applied
	db.Rows = &database.MockMigration_SingleRow{}
	if !migrations.HasMigrationRun("test_migration") {
		t.Error("Should return true when migration has been applied")
	}
}

func TestMigrations_RecordMigration(t *testing.T) {
	db := &database.MockSqlite{}
	container := &app.Container{Db: db}
	migrations := Migrations{Container: container}

	// Test error on initial query
	db.ErrorMode = true
	err := migrations.RecordMigration("test_migration")
	if err == nil {
		t.Error("Should return error when database has an error on query")
	}

	// Test insert new migration (no existing record)
	db.ErrorMode = false
	db.Rows = &database.MockRowsEmpty{}
	db.Result = &database.MockResult{}
	err = migrations.RecordMigration("test_migration")
	if err != nil {
		t.Errorf("Should not return error when inserting: %v", err)
	}

	// Test update existing migration record
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockMigration_SingleRow{})
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.RecordMigration("test_migration")
	if err != nil {
		t.Errorf("Should not return error when updating: %v", err)
	}

	// Test error on INSERT Exec
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // SELECT — no existing record
	db.ErrorAtQuery = 2                        // Fail the INSERT Exec
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.RecordMigration("test_migration")
	if err == nil {
		t.Error("Should return error when INSERT Exec fails")
	}
}

func TestMigrations_MigrateRandomSlugs(t *testing.T) {
	// Test already-applied case
	db := &database.MockSqlite{}
	db.Rows = &database.MockMigration_SingleRow{}
	container := &app.Container{Db: db}
	migrations := Migrations{Container: container}
	err := migrations.MigrateRandomSlugs()
	if err != nil {
		t.Errorf("Already-applied migration should return nil: %v", err)
	}
	if db.Queries != 1 {
		t.Errorf("Expected only 1 query when migration already applied, got %d", db.Queries)
	}

	// Test migration runs — no journal with "random" slug
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.AppendResult(&database.MockRowsEmpty{}) // FindBySlug("random") — not found
	db.AppendResult(&database.MockRowsEmpty{}) // RecordMigration SELECT
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateRandomSlugs()
	if err != nil {
		t.Errorf("Migration should run without errors when no 'random' slug exists: %v", err)
	}

	// Test migration runs — journal with "random" slug exists
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{})       // HasMigrationRun
	db.AppendResult(&database.MockJournal_SingleRow{}) // FindBySlug("random") — found
	db.AppendResult(&database.MockRowsEmpty{})       // RecordMigration SELECT
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateRandomSlugs()
	if err != nil {
		t.Errorf("Migration should run without errors when 'random' slug is found: %v", err)
	}

	// Test RecordMigration failure path
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.AppendResult(&database.MockRowsEmpty{}) // FindBySlug("random") — not found
	db.AppendResult(&database.MockRowsEmpty{}) // RecordMigration SELECT
	db.ErrorAtQuery = 4                        // Fail RecordMigration INSERT
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateRandomSlugs()
	if err == nil {
		t.Error("Expected error when RecordMigration fails")
	}
}

func TestMigrations_MigrateAddTimestamps(t *testing.T) {
	// Test already-applied case
	db := &database.MockSqlite{}
	db.Rows = &database.MockMigration_SingleRow{}
	container := &app.Container{Db: db}
	migrations := Migrations{Container: container}
	err := migrations.MigrateAddTimestamps()
	if err != nil {
		t.Errorf("Already-applied migration should return nil: %v", err)
	}
	if db.Queries != 1 {
		t.Errorf("Expected only 1 query when migration already applied, got %d", db.Queries)
	}

	// Test error on first ALTER TABLE (created_at)
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.ErrorAtQuery = 2                        // Fail the first Exec (ALTER TABLE created_at)
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateAddTimestamps()
	if err == nil {
		t.Error("Expected error when first ALTER TABLE fails")
	}

	// Test error on second ALTER TABLE (updated_at)
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.ErrorAtQuery = 3                        // Fail the second Exec (ALTER TABLE updated_at)
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateAddTimestamps()
	if err == nil {
		t.Error("Expected error when second ALTER TABLE fails")
	}

	// Test migration runs for first time
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.AppendResult(&database.MockRowsEmpty{}) // RecordMigration SELECT
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateAddTimestamps()
	if err != nil {
		t.Errorf("Migration should run without errors: %v", err)
	}

	// Test RecordMigration failure path
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.AppendResult(&database.MockRowsEmpty{}) // RecordMigration SELECT
	db.ErrorAtQuery = 5                        // Fail RecordMigration INSERT (after 2 ALTERs)
	db.Result = &database.MockResult{}
	container = &app.Container{Db: db}
	migrations = Migrations{Container: container}
	err = migrations.MigrateAddTimestamps()
	if err == nil {
		t.Error("Expected error when RecordMigration fails")
	}
}

func TestMigrations_MigrateHTMLToMarkdown(t *testing.T) {
	// Test already-applied case (migration skipped)
	db := &database.MockSqlite{}
	db.Rows = &database.MockMigration_SingleRow{}
	container := &app.Container{
		Db:                db,
		MarkdownProcessor: &markdown.Markdown{},
	}
	migrations := Migrations{Container: container}
	err := migrations.MigrateHTMLToMarkdown()
	if err != nil {
		t.Errorf("Already-applied migration should return nil: %v", err)
	}
	if db.Queries != 1 {
		t.Errorf("Expected only 1 query (the check) when migration already applied, got %d", db.Queries)
	}

	// Test migration runs for the first time
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{})      // HasMigrationRun check
	db.AppendResult(&database.MockJournal_MultipleRows{}) // FetchAll
	db.AppendResult(&database.MockRowsEmpty{})      // RecordMigration SELECT
	db.Result = &database.MockResult{}
	container = &app.Container{
		Db:                db,
		MarkdownProcessor: &markdown.Markdown{},
	}
	migrations = Migrations{Container: container}
	err = migrations.MigrateHTMLToMarkdown()
	if err != nil {
		t.Errorf("Migration should run without errors: %v", err)
	}

	// Test RecordMigration failure path (no journals to migrate)
	db = &database.MockSqlite{}
	db.EnableMultiMode()
	db.AppendResult(&database.MockRowsEmpty{}) // HasMigrationRun
	db.AppendResult(&database.MockRowsEmpty{}) // FetchAll — no entries
	db.AppendResult(&database.MockRowsEmpty{}) // RecordMigration SELECT
	db.ErrorAtQuery = 4                        // Fail RecordMigration INSERT
	db.Result = &database.MockResult{}
	container = &app.Container{
		Db:                db,
		MarkdownProcessor: &markdown.Markdown{},
	}
	migrations = Migrations{Container: container}
	err = migrations.MigrateHTMLToMarkdown()
	if err == nil {
		t.Error("Expected error when RecordMigration fails")
	}
}
