package shared

type PaginatedList[T any] struct {
	Data  []T   `json:"data"`
	Count int64 `json:"count"`
}

func NewPaginatedList[T any](data []T, count int64) *PaginatedList[T] {
	return &PaginatedList[T]{Data: data, Count: count}
}
