package common

type PageResult[T any] struct {
	Data      []T   `json:"data"`
	Page      int   `json:"page"`
	LastPage  int   `json:"lastPage"`
	TotalSize int64 `json:"totalSize"`
}

func NewPageResult[T any](data []T, page int, pageSize int64, totalSize int64) *PageResult[T] {
	lastPage := 1
	if totalSize > 0 && pageSize > 0 {
		lastPage = int(totalSize / pageSize)
	}

	return &PageResult[T]{
		Data:      data,
		Page:      page,
		LastPage:  lastPage,
		TotalSize: totalSize,
	}
}
