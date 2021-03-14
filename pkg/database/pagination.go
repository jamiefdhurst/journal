package database

// PaginationInformation is used to return information from a pagination query
type PaginationInformation struct {
	Page           int
	TotalPages     int
	ResultsPerPage int
	TotalResults   int
}

// PaginationQuery accepts current page and resuilts per page to generate a query
type PaginationQuery struct {
	Page           int
	ResultsPerPage int
}

// PaginationResult is used to return pagination details from a count query
type PaginationResult struct {
	TotalResults int `json:"total"`
}
