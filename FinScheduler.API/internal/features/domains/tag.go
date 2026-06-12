package domains

import (
	"finscheduler/pkg/qh"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Tag struct {
	Id       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	IsActive bool      `db:"is_active"`
}

type TagListingDto struct {
	Id       *uuid.UUID `json:"id"`
	Name     *string    `json:"name"`
	IsActive *bool      `json:"isActive"`
}

type TagDetailedDto struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"isActive"`
}

type TagFilter struct {
	Ids      []*uuid.UUID
	Name     *string
	IsActive *bool
	Page     *int32
	PageSize *int32
}

type TagLookupFilter struct {
	Name     *string
	Page     *int32
	PageSize *int32
}

type TagCreate struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type TagUpdate struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

func NewTagFilter(r *http.Request) (TagFilter, error) {
	queryParams := r.URL.Query()

	ids, err := qh.ParseUUIDs(queryParams, "ids")
	if err != nil {
		return TagFilter{}, err
	}
	name := qh.ParseString(queryParams, "name")
	isActive, err := qh.ParseBool(queryParams, "isActive")
	if err != nil {
		return TagFilter{}, err
	}
	page, err := qh.ParseInt32(queryParams, "page")
	if err != nil {
		return TagFilter{}, err
	}
	pageSize, err := qh.ParseInt32(queryParams, "pageSize")
	if err != nil {
		return TagFilter{}, err
	}

	return TagFilter{
		Ids:      ids,
		Name:     name,
		IsActive: isActive,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func NewTagLookupFilter(r *http.Request) (TagLookupFilter, error) {
	queryParams := r.URL.Query()

	name := qh.ParseString(queryParams, "name")
	page, err := qh.ParseInt32(queryParams, "page")
	if err != nil {
		return TagLookupFilter{}, err
	}
	pageSize, err := qh.ParseInt32(queryParams, "pageSize")
	if err != nil {
		return TagLookupFilter{}, err
	}

	return TagLookupFilter{
		Name:     name,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func NewTagListingDto(tag Tag) *TagListingDto {
	return &TagListingDto{
		Id:       &tag.Id,
		Name:     &tag.Name,
		IsActive: &tag.IsActive,
	}
}

func NewTagDetailedDto(tag Tag) *TagDetailedDto {
	return &TagDetailedDto{
		Name:     &tag.Name,
		IsActive: &tag.IsActive,
	}
}

func (item *TagCreate) Validate() error {
	if len(item.Name) < 3 {
		return fmt.Errorf("name too short")
	}

	return nil
}

func (item *TagUpdate) Validate() error {
	if len(item.Name) < 3 {
		return fmt.Errorf("name too short")
	}

	return nil
}

func (item *TagFilter) Validate() error {
	if item.Page == nil || *item.Page < 0 {
		return fmt.Errorf("page must be zero or greater")
	}
	if item.PageSize == nil || *item.PageSize <= 0 {
		return fmt.Errorf("pageSize must be positive")
	}

	return nil
}

func (item *TagLookupFilter) Validate() error {
	if item.Page == nil || *item.Page < 0 {
		return fmt.Errorf("page must be zero or greater")
	}
	if item.PageSize == nil || *item.PageSize <= 0 {
		return fmt.Errorf("pageSize must be positive")
	}

	return nil
}
