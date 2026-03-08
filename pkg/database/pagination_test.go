package database

import (
    "testing"
)

func TestLinksPagination(t *testing.T) {
    tables := []struct {
        url      string
        info     PaginationInformation
        previous string
        next     string
    }{
        // Single page — no links
        {"/posts", PaginationInformation{1, 1, 20, 20}, "", ""},
        // First page of many — only next
        {"/posts", PaginationInformation{1, 3, 20, 60}, "", "/posts?page=2"},
        // Last page — only previous
        {"/posts", PaginationInformation{3, 3, 20, 60}, "/posts?page=2", ""},
        // Middle page — both links
        {"/posts", PaginationInformation{2, 3, 20, 60}, "/posts?page=1", "/posts?page=3"},
    }

    for _, table := range tables {
        links := LinksPagination(table.url, table.info)
        if links.Previous != table.previous {
            t.Errorf("LinksPagination(%v): expected Previous %q, got %q", table.info, table.previous, links.Previous)
        }
        if links.Next != table.next {
            t.Errorf("LinksPagination(%v): expected Next %q, got %q", table.info, table.next, links.Next)
        }
    }
}

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
            PaginationInformation{1, 7, 20, 135},
            PaginationDisplay{false, false, 1, 1, 7, 7},
        },
        {
            PaginationInformation{1, 9, 20, 175},
            PaginationDisplay{false, true, 1, 1, 7, 9},
        },
        {
            PaginationInformation{1, 15, 20, 299},
            PaginationDisplay{false, true, 1, 1, 7, 15},
        },
        {
            PaginationInformation{15, 15, 20, 299},
            PaginationDisplay{true, false, 15, 9, 15, 15},
        },
        {
            PaginationInformation{7, 15, 20, 299},
            PaginationDisplay{true, true, 7, 4, 10, 15},
        },
    }

    for _, table := range tables {
        actual := DisplayPagination(table.input)
        if actual != table.output {
            t.Errorf("Expected DisplayPagination() to produce result of '%v', got '%v'", table.output, actual)
        }
    }
}
