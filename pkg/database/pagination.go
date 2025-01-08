package database

import "math"

const PAGINATION_MAX_PAGES = 9

// PaginationDisplay describes the information needed to display a pagination component
type PaginationDisplay struct {
	ShowLeft    bool
	ShowRight   bool
	CurrentPage int
	FirstPage   int
	LastPage    int
	TotalPages  int
}

// PaginationInformation is used to return information from a pagination query
type PaginationInformation struct {
	Page           int
	TotalPages     int
	ResultsPerPage int
	TotalResults   int
}

// PaginationQuery accepts current page and results per page to generate a query
type PaginationQuery struct {
	Page           int
	ResultsPerPage int
}

// PaginationResult is used to return pagination details from a count query
type PaginationResult struct {
	TotalResults int `json:"total"`
}

func DisplayPagination(info PaginationInformation) PaginationDisplay {
	display := PaginationDisplay{false, false, info.Page, 1, info.TotalPages, info.TotalPages}
	if info.TotalPages <= PAGINATION_MAX_PAGES {
		return display
	}
	half := int(math.Floor(PAGINATION_MAX_PAGES / 2))
	if info.Page-half > 1 {
		display.ShowLeft = true
		display.FirstPage = info.Page - half
		if info.TotalPages-half <= info.Page {
			display.FirstPage = info.TotalPages - PAGINATION_MAX_PAGES + 1
		}
	}
	if info.Page+half < info.TotalPages {
		display.ShowRight = true
		display.LastPage = info.Page + half
		if info.Page-half <= 1 {
			display.LastPage = PAGINATION_MAX_PAGES
		}
	}

	return display
}
