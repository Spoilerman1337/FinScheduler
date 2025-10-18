package shared

type PaginatedList[T any] struct {
	data  []T
	count int64
}

func NewPaginatedList[T any](data []T, count int64) *PaginatedList[T] {
	return &PaginatedList[T]{data: data, count: count}
}
