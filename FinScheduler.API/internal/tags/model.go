package tags

import (
	"finscheduler/internal/shared"
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

type TagDto struct {
	Id       *uuid.UUID `json:"id"`
	Name     *string    `json:"name"`
	IsActive *bool      `json:"isActive"`
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

type ItemTags struct {
	ItemId *uuid.UUID        `json:"item_id" db:"item_id"`
	Tags   shared.LookupJSON `json:"tags" db:"tags"`
}

type TagToItem struct {
	TagId  *uuid.UUID `db:"tag_id"`
	ItemId *uuid.UUID `db:"item_id"`
}

func NewTagsFilter(r *http.Request) TagFilter {
	queryParams := r.URL.Query()

	ids := qh.ParseUUIDs(queryParams, "ids")
	name := qh.ParseString(queryParams, "name")
	isActive := qh.ParseBool(queryParams, "isActive")
	page := qh.ParseInt32(queryParams, "page")
	pageSize := qh.ParseInt32(queryParams, "pageSize")

	return TagFilter{
		Ids:      ids,
		Name:     name,
		IsActive: isActive,
		Page:     page,
		PageSize: pageSize,
	}
}

func NewTagDto(tag Tag) *TagDto {
	return &TagDto{
		Id:       &tag.Id,
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
		return fmt.Errorf("page is negative")
	}
	if item.PageSize == nil || *item.PageSize < 0 {
		return fmt.Errorf("pageSize is negative")
	}

	return nil
}
