package result

// Result summarises an executed database command
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
