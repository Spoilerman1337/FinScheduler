package items

import (
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Item struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Price       int64     `db:"price"`
	Description string    `db:"description"`
	IsActive    bool      `db:"is_active"`
}

type ItemFilter struct {
	Ids         []*uuid.UUID
	Name        *string
	PriceFrom   *int64
	PriceTo     *int64
	Description *string
	IsActive    *bool
	Page        *int32
	PageSize    *int32
}

type ItemCreate struct {
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type ItemUpdate struct {
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func NewItemFilter(r *http.Request) ItemFilter {
	queryParams := r.URL.Query()

	var ids []*uuid.UUID
	idsParams := queryParams["ids"]
	for _, idParam := range idsParams {
		id, err := uuid.Parse(idParam)
		if err == nil {
			ids = append(ids, &id)
		}
	}

	name := queryParams.Get("name")
	priceFrom, _ := strconv.ParseInt(queryParams.Get("priceFrom"), 10, 64)
	priceTo, _ := strconv.ParseInt(queryParams.Get("priceTo"), 10, 64)
	description := queryParams.Get("description")
	isActive, _ := strconv.ParseBool(queryParams.Get("isActive"))
	page, _ := strconv.ParseInt(queryParams.Get("page"), 10, 32)
	pageSize, _ := strconv.ParseInt(queryParams.Get("pageSize"), 10, 32)

	page32 := int32(page)
	pageSize32 := int32(pageSize)

	return ItemFilter{
		Ids:         ids,
		Name:        &name,
		PriceFrom:   &priceFrom,
		PriceTo:     &priceTo,
		Description: &description,
		IsActive:    &isActive,
		Page:        &page32,
		PageSize:    &pageSize32,
	}
}
