package items

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"time"
)

type Item struct {
	Id          uuid.UUID       `db:"id"`
	Name        string          `db:"name"`
	Price       decimal.Decimal `db:"price"`
	Description string          `db:"description"`
	IsActive    bool            `db:"is_active"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}

type ItemDto struct {
	Id          *uuid.UUID
	Name        *string
	Price       *float64
	Description *string
	IsActive    *bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type ItemFilter struct {
	Ids         []*uuid.UUID
	Name        *string
	PriceFrom   *decimal.Decimal
	PriceTo     *decimal.Decimal
	Description *string
	IsActive    *bool
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	UpdatedFrom *time.Time
	UpdatedTo   *time.Time
	Page        *int32
	PageSize    *int32
}

type ItemCreate struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
}

type ItemUpdate struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
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
	priceFrom, _ := decimal.NewFromString(queryParams.Get("priceFrom"))
	priceTo, _ := decimal.NewFromString(queryParams.Get("priceTo"))
	description := queryParams.Get("description")
	isActive, _ := strconv.ParseBool(queryParams.Get("isActive"))
	page, _ := strconv.ParseInt(queryParams.Get("page"), 10, 32)
	pageSize, _ := strconv.ParseInt(queryParams.Get("pageSize"), 10, 32)
	createdFrom, _ := time.Parse(time.RFC3339, queryParams.Get("createdFrom"))
	createdTo, _ := time.Parse(time.RFC3339, queryParams.Get("createdTo"))
	updatedFrom, _ := time.Parse(time.RFC3339, queryParams.Get("updatedFrom"))
	updatedTo, _ := time.Parse(time.RFC3339, queryParams.Get("updatedTo"))

	page32 := int32(page)
	pageSize32 := int32(pageSize)

	return ItemFilter{
		Ids:         ids,
		Name:        &name,
		PriceFrom:   &priceFrom,
		PriceTo:     &priceTo,
		Description: &description,
		IsActive:    &isActive,
		CreatedFrom: &createdFrom,
		CreatedTo:   &createdTo,
		UpdatedFrom: &updatedFrom,
		UpdatedTo:   &updatedTo,
		Page:        &page32,
		PageSize:    &pageSize32,
	}
}

func NewItemDto(item *Item) *ItemDto {
	if item == nil {
		return nil
	}

	var updatedAt *time.Time
	if &item.UpdatedAt != nil && item.UpdatedAt.Valid {
		updatedAt = &item.UpdatedAt.Time
	} else {
		updatedAt = nil
	}

	price, _ := item.Price.Float64()

	return &ItemDto{
		Id:          &item.Id,
		Name:        &item.Name,
		Description: &item.Description,
		IsActive:    &item.IsActive,
		CreatedAt:   &item.CreatedAt,
		Price:       &price,
		UpdatedAt:   updatedAt,
	}
}
