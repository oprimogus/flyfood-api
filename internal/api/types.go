package api

import "math"

type Pagination[T any] struct {
	Data            []T  `json:"data"`
	CurrentPage     int  `json:"currentPage"`
	PageSize        int  `json:"pageSize"`
	TotalPages      int  `json:"totalPages"`
	TotalItems      int  `json:"totalItems"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

func Paginate[T any](items []T, currentPage, pageSize, totalItems int) Pagination[T] {
	if totalItems == 0 {
		return Pagination[T]{
			Data:        make([]T, 0),
			CurrentPage: 1,
		}
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	if currentPage > totalPages {
		currentPage = totalPages
	} else if currentPage < 1 {
		currentPage = 1
	}

	return Pagination[T]{
		Data:            items,
		CurrentPage:     currentPage,
		PageSize:        pageSize,
		TotalItems:      totalItems,
		TotalPages:      totalPages,
		HasNextPage:     currentPage < totalPages,
		HasPreviousPage: currentPage > 1,
	}
}
