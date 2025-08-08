package common

import "math"

type PageResult[T any] struct {
	Data      []T   `json:"data"`
	Page      int   `json:"page"`
	LastPage  int   `json:"lastPage"`
	TotalSize int64 `json:"totalSize"`
}

func NewPageResult[T any](data []T, page int, pageSize int64, totalSize int64) *PageResult[T] {
	lastPage := 1
	if totalSize > 0 && pageSize > 0 {
		ceil := math.Ceil(float64(totalSize) / float64(pageSize))
		lastPage = int(ceil)
	}

	return &PageResult[T]{
		Data:      data,
		Page:      page,
		LastPage:  lastPage,
		TotalSize: totalSize,
	}
}
