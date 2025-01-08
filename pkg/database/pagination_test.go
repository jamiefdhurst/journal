package database

import (
	"testing"
)

func TestDisplayPagination(t *testing.T) {
	tables := []struct {
		input  PaginationInformation
		output PaginationDisplay
	}{
		{
			PaginationInformation{1, 1, 20, 1},
			PaginationDisplay{false, false, 1, 1, 1, 1},
		},
		{
			PaginationInformation{1, 4, 20, 70},
			PaginationDisplay{false, false, 1, 1, 4, 4},
		},
		{
			PaginationInformation{1, 9, 20, 175},
			PaginationDisplay{false, false, 1, 1, 9, 9},
		},
		{
			PaginationInformation{1, 15, 20, 299},
			PaginationDisplay{false, true, 1, 1, 9, 15},
		},
		{
			PaginationInformation{15, 15, 20, 299},
			PaginationDisplay{true, false, 15, 7, 15, 15},
		},
		{
			PaginationInformation{7, 15, 20, 299},
			PaginationDisplay{true, true, 7, 3, 11, 15},
		},
	}

	for _, table := range tables {
		actual := DisplayPagination(table.input)
		if actual != table.output {
			t.Errorf("Expected DisplayPagination() to produce result of '%v', got '%v'", table.output, actual)
		}
	}
}
