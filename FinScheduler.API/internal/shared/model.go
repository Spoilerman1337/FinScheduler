package shared

import (
	"encoding/json"
	"errors"
)

type PaginatedList[T any] struct {
	Data  []T   `json:"data"`
	Count int64 `json:"count"`
}

func NewPaginatedList[T any](data []T, count int64) *PaginatedList[T] {
	return &PaginatedList[T]{Data: data, Count: count}
}

type Lookup struct {
	Value *string `json:"value" db:"value"`
	Label *string `json:"label" db:"label"`
}

type LookupJSON []Lookup

func (t *LookupJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, t)
}
