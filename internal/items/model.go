package items

import (
	"database/sql"
	qh "finscheduler/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
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
	IsActive    bool            `json:"isActive"`
}

type ItemUpdate struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	IsActive    bool            `json:"isActive"`
}

func NewItemFilter(r *http.Request) ItemFilter {
	queryParams := r.URL.Query()

	ids := qh.ParseUUIDs(queryParams, "ids")
	name := qh.ParseString(queryParams, "name")
	priceFrom := qh.ParseDecimal(queryParams, "priceFrom")
	priceTo := qh.ParseDecimal(queryParams, "priceTo")
	description := qh.ParseString(queryParams, "description")
	isActive := qh.ParseBool(queryParams, "isActive")
	page := qh.ParseInt32(queryParams, "page")
	pageSize := qh.ParseInt32(queryParams, "pageSize")
	createdFrom := qh.ParseTime(queryParams, "createdFrom")
	createdTo := qh.ParseTime(queryParams, "createdTo")
	updatedFrom := qh.ParseTime(queryParams, "updatedFrom")
	updatedTo := qh.ParseTime(queryParams, "updatedTo")

	return ItemFilter{
		Ids:         ids,
		Name:        name,
		PriceFrom:   priceFrom,
		PriceTo:     priceTo,
		Description: description,
		IsActive:    isActive,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		UpdatedFrom: updatedFrom,
		UpdatedTo:   updatedTo,
		Page:        page,
		PageSize:    pageSize,
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

func (item *ItemCreate) Validate() error {
	if len(item.Name) < 3 {
		return fmt.Errorf("name too short")
	}
	if item.Price.IsNegative() {
		return fmt.Errorf("price is negative")
	}

	return nil
}

func (item *ItemUpdate) Validate() error {
	if len(item.Name) < 3 {
		return fmt.Errorf("name too short")
	}
	if item.Price.IsNegative() {
		return fmt.Errorf("price is negative")
	}

	return nil
}

func (item *ItemFilter) Validate() error {
	if item.Page == nil || *item.Page < 0 {
		return fmt.Errorf("page is negative")
	}
	if item.PageSize == nil || *item.PageSize < 0 {
		return fmt.Errorf("pageSize is negative")
	}
	if item.PriceFrom != nil && item.PriceTo != nil && (*item.PriceTo).LessThan(*item.PriceFrom) {
		return fmt.Errorf("priceTo cannot be lesser than priceFrom")
	}
	if item.CreatedFrom != nil && item.CreatedTo != nil && (*item.CreatedTo).Before(*item.CreatedFrom) {
		return fmt.Errorf("createTo cannot be earlier than createFrom")
	}
	if item.UpdatedFrom != nil && item.UpdatedTo != nil && (*item.UpdatedTo).Before(*item.UpdatedFrom) {
		return fmt.Errorf("updateTo cannot be earlier than updateFrom")
	}

	return nil
}
