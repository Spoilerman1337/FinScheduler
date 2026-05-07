package domains

import (
	"database/sql"
	"finscheduler/pkg/qh"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	Id          uuid.UUID       `db:"id"`
	Name        string          `db:"name"`
	Price       decimal.Decimal `db:"price"`
	Description string          `db:"description"`
	IsActive    bool            `db:"is_active"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
	Cashback    int32           `db:"cashback"`
	Category    ItemCategory    `db:"category"`
}

type ItemDto struct {
	Id          *uuid.UUID    `json:"id"`
	Name        *string       `json:"name"`
	Price       *float64      `json:"price"`
	Description *string       `json:"description"`
	IsActive    *bool         `json:"isActive"`
	CreatedAt   *time.Time    `json:"createdAt"`
	UpdatedAt   *time.Time    `json:"updatedAt"`
	Cashback    *int32        `json:"cashback"`
	Category    *ItemCategory `json:"category"`
	Tags        []*Lookup     `json:"tags"`
}

type ItemFilter struct {
	Ids          []*uuid.UUID
	Name         *string
	PriceFrom    *decimal.Decimal
	PriceTo      *decimal.Decimal
	Description  *string
	IsActive     *bool
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	UpdatedFrom  *time.Time
	UpdatedTo    *time.Time
	CashbackFrom *int32
	CashbackTo   *int32
	Categories   []*ItemCategory
	TagIds       []*uuid.UUID
	Page         *int32
	PageSize     *int32
}

type ItemCreate struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	IsActive    bool            `json:"isActive"`
	Cashback    int32           `json:"cashback"`
	Category    string          `json:"category"`
	TagIds      []string        `json:"tagIds"`
}

type ItemUpdate struct {
	Name        string          `json:"name"`
	Price       decimal.Decimal `json:"price"`
	Description string          `json:"description"`
	IsActive    bool            `json:"isActive"`
	Cashback    int32           `json:"cashback"`
	Category    string          `json:"category"`
	TagIds      []string        `json:"tagIds"`
}

func NewItemFilter(r *http.Request) (ItemFilter, error) {
	queryParams := r.URL.Query()

	ids, err := qh.ParseUUIDs(queryParams, "ids")
	if err != nil {
		return ItemFilter{}, err
	}
	name := qh.ParseString(queryParams, "name")
	priceFrom, err := qh.ParseDecimal(queryParams, "priceFrom")
	if err != nil {
		return ItemFilter{}, err
	}
	priceTo, err := qh.ParseDecimal(queryParams, "priceTo")
	if err != nil {
		return ItemFilter{}, err
	}
	description := qh.ParseString(queryParams, "description")
	isActive, err := qh.ParseBool(queryParams, "isActive")
	if err != nil {
		return ItemFilter{}, err
	}
	page, err := qh.ParseInt32(queryParams, "page")
	if err != nil {
		return ItemFilter{}, err
	}
	pageSize, err := qh.ParseInt32(queryParams, "pageSize")
	if err != nil {
		return ItemFilter{}, err
	}
	createdFrom, err := qh.ParseTime(queryParams, "createdFrom")
	if err != nil {
		return ItemFilter{}, err
	}
	createdTo, err := qh.ParseTime(queryParams, "createdTo")
	if err != nil {
		return ItemFilter{}, err
	}
	updatedFrom, err := qh.ParseTime(queryParams, "updatedFrom")
	if err != nil {
		return ItemFilter{}, err
	}
	updatedTo, err := qh.ParseTime(queryParams, "updatedTo")
	if err != nil {
		return ItemFilter{}, err
	}
	cashbackFrom, err := qh.ParseInt32(queryParams, "cashbackFrom")
	if err != nil {
		return ItemFilter{}, err
	}
	cashbackTo, err := qh.ParseInt32(queryParams, "cashbackTo")
	if err != nil {
		return ItemFilter{}, err
	}
	categories, err := qh.ParseEnums[ItemCategory](queryParams, "categories")
	if err != nil {
		return ItemFilter{}, err
	}
	tagIds, err := qh.ParseUUIDs(queryParams, "tagIds")
	if err != nil {
		return ItemFilter{}, err
	}

	return ItemFilter{
		Ids:          ids,
		Name:         name,
		PriceFrom:    priceFrom,
		PriceTo:      priceTo,
		Description:  description,
		IsActive:     isActive,
		CreatedFrom:  createdFrom,
		CreatedTo:    createdTo,
		UpdatedFrom:  updatedFrom,
		UpdatedTo:    updatedTo,
		CashbackFrom: cashbackFrom,
		CashbackTo:   cashbackTo,
		Categories:   categories,
		TagIds:       tagIds,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

func NewItemDto(item Item, tags []Tag) *ItemDto {
	var updatedAt *time.Time
	if item.UpdatedAt.Valid {
		updatedAt = &item.UpdatedAt.Time
	} else {
		updatedAt = nil
	}

	price, _ := item.Price.Float64()

	var tagLookups []*Lookup
	for _, tag := range tags {
		value := tag.Id.String()
		tagLookup := Lookup{Label: &tag.Name, Value: &value}
		tagLookups = append(tagLookups, &tagLookup)
	}

	return &ItemDto{
		Id:          &item.Id,
		Name:        &item.Name,
		Description: &item.Description,
		IsActive:    &item.IsActive,
		CreatedAt:   &item.CreatedAt,
		Price:       &price,
		UpdatedAt:   updatedAt,
		Cashback:    &item.Cashback,
		Category:    &item.Category,
		Tags:        tagLookups,
	}
}

func (item *ItemCreate) Validate() error {
	if len(item.Name) < 3 {
		return fmt.Errorf("name too short")
	}
	if item.Price.IsNegative() {
		return fmt.Errorf("price is negative")
	}
	if item.Cashback < 0 {
		return fmt.Errorf("cashback is negative")
	}
	if !ItemCategory(item.Category).IsValid() {
		return fmt.Errorf("category is invalid")
	}
	if err := validateTagIds(item.TagIds); err != nil {
		return err
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
	if item.Cashback < 0 {
		return fmt.Errorf("cashback is negative")
	}
	if !ItemCategory(item.Category).IsValid() {
		return fmt.Errorf("category is invalid")
	}
	if err := validateTagIds(item.TagIds); err != nil {
		return err
	}

	return nil
}

func (item *ItemFilter) Validate() error {
	if item.Page == nil || *item.Page < 0 {
		return fmt.Errorf("page must be zero or greater")
	}
	if item.PageSize == nil || *item.PageSize <= 0 {
		return fmt.Errorf("pageSize must be positive")
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
	if item.CashbackFrom != nil && item.CashbackTo != nil && *item.CashbackTo < *item.CashbackFrom {
		return fmt.Errorf("cashbackFrom cannot be lesser than cashbackTO")
	}

	return nil
}

type ItemCategory string

const (
	FoodDrinks     ItemCategory = "FoodDrinks"
	Subscriptions  ItemCategory = "Subscriptions"
	Health         ItemCategory = "Health"
	Beauty         ItemCategory = "Beauty"
	Gifts          ItemCategory = "Gifts"
	Transport      ItemCategory = "Transport"
	Entertainments ItemCategory = "Entertainments"
	Meds           ItemCategory = "Meds"
	Travel         ItemCategory = "Travel"
	Sports         ItemCategory = "Sports"
	Telecom        ItemCategory = "Telecom"
	Education      ItemCategory = "Education"
)

func (itemCategory ItemCategory) IsValid() bool {
	switch itemCategory {
	case FoodDrinks, Subscriptions, Health, Beauty, Gifts, Transport, Entertainments, Meds, Travel, Sports, Telecom, Education:
		return true
	default:
		return false
	}
}

func validateTagIds(tagIds []string) error {
	seen := make(map[uuid.UUID]struct{}, len(tagIds))

	for _, tagId := range tagIds {
		parsedTagID, err := uuid.Parse(tagId)
		if err != nil {
			return fmt.Errorf("tagId is invalid: %s", tagId)
		}

		if _, exists := seen[parsedTagID]; exists {
			return fmt.Errorf("tagId is duplicated: %s", tagId)
		}

		seen[parsedTagID] = struct{}{}
	}

	return nil
}
