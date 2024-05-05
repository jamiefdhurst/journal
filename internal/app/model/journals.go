package model

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jamiefdhurst/journal/internal/app"
	"github.com/jamiefdhurst/journal/pkg/database"
	"github.com/jamiefdhurst/journal/pkg/database/dynamodb"
	"github.com/jamiefdhurst/journal/pkg/database/rows"
	"github.com/jamiefdhurst/journal/pkg/database/sqlite"
)

type JournalStore interface {
	CreateTable() error
	EnsureUniqueSlug(string, int) string
	FetchAll() []Journal
	FetchPaginated(database.PaginationQuery) ([]Journal, database.PaginationInformation)
	FindBySlug(string) Journal
	FindNext(int) Journal
	FindPrev(int) Journal
	Save(Journal) Journal
}

const journalTable = "journal"

// Journals Common database resource link for Journal actions
type Journals struct {
	Container *app.Container
	Gs        GiphysExtractor
}

// CreateTable Create the actual table
func (js *Journals) CreateTable() error {
	db := js.Container.Db.(sqlite.SqliteLike)
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `" + journalTable + "` (" +
		"`id` INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"`slug` VARCHAR(255) NOT NULL, " +
		"`title` VARCHAR(255) NOT NULL, " +
		"`date` DATE NOT NULL, " +
		"`content` TEXT NOT NULL" +
		")")

	return err
}

// EnsureUniqueSlug Make sure the current slug is unique
func (js *Journals) EnsureUniqueSlug(slug string, addition int) string {
	newSlug := slug
	if addition > 0 {
		newSlug = strings.Join([]string{slug, "-", strconv.Itoa(addition)}, "")
	}
	exists := js.FindBySlug(newSlug)
	if exists.ID > 0 {
		addition++
		return js.EnsureUniqueSlug(slug, addition)
	}

	return newSlug
}

// FetchAll Get all journals
func (js *Journals) FetchAll() []Journal {
	db := js.Container.Db.(sqlite.SqliteLike)
	rows, err := db.Query("SELECT * FROM `" + journalTable + "` ORDER BY `date` DESC")
	if err != nil {
		return []Journal{}
	}

	return js.loadFromRows(rows)
}

// FetchPaginated returns a set of paginated journal entries
func (js *Journals) FetchPaginated(query database.PaginationQuery) ([]Journal, database.PaginationInformation) {
	pagination := database.PaginationInformation{
		Page:           query.Page,
		ResultsPerPage: query.ResultsPerPage,
	}

	db := js.Container.Db.(sqlite.SqliteLike)
	countResult, err := db.Query("SELECT COUNT(*) AS `total` FROM `" + journalTable + "`")
	if err != nil {
		return []Journal{}, pagination
	}
	countResult.Next()
	countResult.Scan(&pagination.TotalResults)
	countResult.Close()
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalResults) / float64(query.ResultsPerPage)))

	if query.Page > pagination.TotalPages {
		return []Journal{}, pagination
	}

	rows, _ := db.Query(fmt.Sprintf("SELECT * FROM `"+journalTable+"` ORDER BY `date` DESC LIMIT %d OFFSET %d", query.ResultsPerPage, (query.Page-1)*query.ResultsPerPage))
	return js.loadFromRows(rows), pagination
}

// FindBySlug Find a journal by slug
func (js *Journals) FindBySlug(slug string) Journal {
	db := js.Container.Db.(sqlite.SqliteLike)

	return js.loadSingle(db.Query("SELECT * FROM `"+journalTable+"` WHERE `slug` = ? LIMIT 1", slug))
}

// FindNext returns the next entry after an ID
func (js *Journals) FindNext(id int) Journal {
	db := js.Container.Db.(sqlite.SqliteLike)

	return js.loadSingle(db.Query("SELECT * FROM `"+journalTable+"` WHERE `id` > ? ORDER BY `id` LIMIT 1", strconv.Itoa(id)))
}

// FindNext returns the previous entry before an ID
func (js *Journals) FindPrev(id int) Journal {
	db := js.Container.Db.(sqlite.SqliteLike)

	return js.loadSingle(db.Query("SELECT * FROM `"+journalTable+"` WHERE `id` < ? ORDER BY `id` DESC LIMIT 1", strconv.Itoa(id)))
}

// Save Save a journal entry, either inserting it or updating it in the database
func (js *Journals) Save(j Journal) Journal {
	var res sql.Result

	// Convert content for saving
	j.Content = js.Gs.ExtractContentsAndSearchAPI(j.Content)
	if j.Slug == "" {
		j.Slug = Slugify(j.Title)
	}

	db := js.Container.Db.(sqlite.SqliteLike)
	if j.ID == 0 {
		j.Slug = js.EnsureUniqueSlug(j.Slug, 0)
		res, _ = db.Exec("INSERT INTO `"+journalTable+"` (`slug`, `title`, `date`, `content`) VALUES(?,?,?,?)", j.Slug, j.Title, j.Date, j.Content)
	} else {
		res, _ = db.Exec("UPDATE `"+journalTable+"` SET `slug` = ?, `title` = ?, `date` = ?, `content` = ? WHERE `id` = ?", j.Slug, j.Title, j.Date, j.Content, strconv.Itoa(j.ID))
	}

	// Store insert ID
	if j.ID == 0 {
		id, _ := res.LastInsertId()
		j.ID = int(id)
	}

	return j
}

func (js Journals) loadFromRows(rows rows.Rows) []Journal {
	defer rows.Close()
	journals := []Journal{}
	for rows.Next() {
		j := Journal{}
		rows.Scan(&j.ID, &j.Slug, &j.Title, &j.Date, &j.Content)
		journals = append(journals, j)
	}

	return journals
}

func (js *Journals) loadSingle(rows rows.Rows, err error) Journal {
	if err != nil {
		return Journal{}
	}
	journals := js.loadFromRows(rows)

	if len(journals) == 1 {
		return journals[0]
	}

	return Journal{}
}

// DynamoJournals DynamoDB resource link for Journal actions
type DynamoJournals struct {
	Container *app.Container
	Gs        GiphysExtractor
}

// CreateTable Non-functional, here for interface composition
func (js *DynamoJournals) CreateTable() error {
	return nil
}

// EnsureUniqueSlug Make sure the current slug is unique
func (js *DynamoJournals) EnsureUniqueSlug(slug string, addition int) string {
	newSlug := slug
	if addition > 0 {
		newSlug = strings.Join([]string{slug, "-", strconv.Itoa(addition)}, "")
	}
	exists := js.FindBySlug(newSlug)
	if exists.ID > 0 {
		addition++
		return js.EnsureUniqueSlug(slug, addition)
	}

	return newSlug
}

// FetchAll Get all journals
func (js *DynamoJournals) FetchAll() []Journal {
	db := js.Container.Db.(dynamodb.DynamodbLike)
	expr, _ := expression.NewBuilder().Build()
	journals := []Journal{}
	err := db.Scan(expr, &journals)
	if err != nil {
		return []Journal{}
	}
	return journals
}

// FetchPaginated returns a set of paginated journal entries
func (js *DynamoJournals) FetchPaginated(query database.PaginationQuery) ([]Journal, database.PaginationInformation) {
	pagination := database.PaginationInformation{
		Page:           query.Page,
		ResultsPerPage: query.ResultsPerPage,
	}

	db := js.Container.Db.(dynamodb.DynamodbLike)
	expr, _ := expression.NewBuilder().Build()
	count, _ := db.ScanCount(expr)
	pagination.TotalResults = int(count)
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalResults) / float64(query.ResultsPerPage)))

	if query.Page > pagination.TotalPages {
		return []Journal{}, pagination
	}

	journals := []Journal{}
	err := db.ScanLimit(expr, (query.Page-1)*query.ResultsPerPage, query.ResultsPerPage, &journals)
	if err != nil {
		return []Journal{}, pagination
	}

	return journals, pagination
}

// FindBySlug Find a journal by slug
func (js *DynamoJournals) FindBySlug(slug string) Journal {
	db := js.Container.Db.(dynamodb.DynamodbLike)
	expr, _ := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("slug"), expression.Value(slug)),
	).Build()
	journals := []Journal{}
	err := db.ScanLimit(expr, 0, 1, &journals)
	if err != nil || len(journals) < 1 {
		return Journal{}
	}
	return journals[0]
}

// FindNext returns the next entry after an ID
func (js *DynamoJournals) FindNext(id int) Journal {
	return Journal{}
	// TODO: Write this using scan so it works correctly
	// db := js.Container.Db.(dynamodb.DynamodbLike)
	// expr, _ := expression.NewBuilder().WithKeyCondition(
	// 	expression.Key("id").GreaterThan(expression.Value(id))).Build()
	// journals := []Journal{}
	// err := db.Query(expr, true, &journals)
	// if err != nil || len(journals) < 1 {
	// 	return Journal{}
	// }
	// return journals[0]
}

// FindNext returns the previous entry before an ID
func (js *DynamoJournals) FindPrev(id int) Journal {
	return Journal{}
	// TODO: Write this using scan so it works correctly
	// db := js.Container.Db.(dynamodb.DynamodbLike)
	// expr, _ := expression.NewBuilder().WithKeyCondition(
	// 	expression.Key("id").LessThan(expression.Value(id))).Build()
	// journals := []Journal{}
	// err := db.Query(expr, false, &journals)
	// if err != nil || len(journals) < 1 {
	// 	return Journal{}
	// }
	// return journals[0]
}

// Save a journal entry, either inserting it or updating it in the database
func (js *DynamoJournals) Save(j Journal) Journal {
	db := js.Container.Db.(dynamodb.DynamodbLike)
	if j.ID == 0 {
		lastJournal := js.FindPrev(math.MaxInt)
		if lastJournal.ID == 0 {
			j.ID = 1
		} else {
			j.ID = lastJournal.ID - 1
		}
	}

	db.PutItem(map[string]types.AttributeValue{
		"id":      &types.AttributeValueMemberN{Value: strconv.Itoa(j.ID)},
		"slug":    &types.AttributeValueMemberS{Value: j.Slug},
		"title":   &types.AttributeValueMemberS{Value: j.Title},
		"date":    &types.AttributeValueMemberS{Value: j.Date},
		"content": &types.AttributeValueMemberS{Value: j.Content},
	})

	return j
}

func NewJournalStore(container *app.Container, giphys GiphysExtractor) JournalStore {
	if container.Configuration.Database == database.Dynamodb {
		return &DynamoJournals{Container: container, Gs: giphys}
	}

	return &Journals{Container: container, Gs: giphys}
}
