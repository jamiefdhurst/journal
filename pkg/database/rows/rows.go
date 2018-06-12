package rows

// Rows Define a common interface for a result of rows
type Rows interface {
	Close() error
	Columns() ([]string, error)
	Next() bool
	Scan(dest ...interface{}) error
}
