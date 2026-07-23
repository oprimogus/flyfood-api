package api

import (
	"math"
	"testing"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name        string
		items       []int
		currentPage int
		pageSize    int
		totalItems  int
		expected    Pagination[int]
	}{
		{
			name:        "Normal pagination",
			items:       []int{1, 2, 3, 4, 5},
			currentPage: 2,
			pageSize:    2,
			totalItems:  5,
			expected: Pagination[int]{
				Data:            []int{1, 2, 3, 4, 5},
				CurrentPage:     2,
				PageSize:        2,
				TotalPages:      int(math.Ceil(float64(5) / 2)),
				TotalItems:      5,
				HasNextPage:     true,
				HasPreviousPage: true,
			},
		},
		{
			name:        "First page boundary",
			items:       []int{1, 2, 3},
			currentPage: 0,
			pageSize:    2,
			totalItems:  3,
			expected: Pagination[int]{
				Data:            []int{1, 2, 3},
				CurrentPage:     1,
				PageSize:        2,
				TotalPages:      2,
				TotalItems:      3,
				HasNextPage:     true,
				HasPreviousPage: false,
			},
		},
		{
			name:        "Last page boundary",
			items:       []int{1, 2, 3, 4},
			currentPage: 10,
			pageSize:    2,
			totalItems:  4,
			expected: Pagination[int]{
				Data:            []int{1, 2, 3, 4},
				CurrentPage:     2,
				PageSize:        2,
				TotalPages:      2,
				TotalItems:      4,
				HasNextPage:     false,
				HasPreviousPage: true,
			},
		},
		{
			name:        "Empty items",
			items:       []int{},
			currentPage: 1,
			pageSize:    2,
			totalItems:  0,
			expected: Pagination[int]{
				CurrentPage: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Paginate(tt.items, tt.currentPage, tt.pageSize, tt.totalItems)

			if len(result.Data) != len(tt.expected.Data) {
				t.Errorf("Data length mismatch: got %v, expected %v", len(result.Data), len(tt.expected.Data))
			}
			if result.CurrentPage != tt.expected.CurrentPage {
				t.Errorf("CurrentPage mismatch: got %d, expected %d", result.CurrentPage, tt.expected.CurrentPage)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("PageSize mismatch: got %d, expected %d", result.PageSize, tt.expected.PageSize)
			}
			if result.TotalPages != tt.expected.TotalPages {
				t.Errorf("TotalPages mismatch: got %d, expected %d", result.TotalPages, tt.expected.TotalPages)
			}
			if result.TotalItems != tt.expected.TotalItems {
				t.Errorf("TotalItems mismatch: got %d, expected %d", result.TotalItems, tt.expected.TotalItems)
			}
			if result.HasNextPage != tt.expected.HasNextPage {
				t.Errorf("HasNextPage mismatch: got %v, expected %v", result.HasNextPage, tt.expected.HasNextPage)
			}
			if result.HasPreviousPage != tt.expected.HasPreviousPage {
				t.Errorf("HasPreviousPage mismatch: got %v, expected %v", result.HasPreviousPage, tt.expected.HasPreviousPage)
			}
		})
	}
}
