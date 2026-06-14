package domains

import "github.com/google/uuid"

type TagToItem struct {
	ItemId uuid.UUID `db:"item_id"`
	TagId  uuid.UUID `db:"tag_id"`
}

type TagToItemCreate struct {
	ItemId *uuid.UUID   `json:"itemId"`
	TagIds []*uuid.UUID `json:"tagId"`
}

type TagToItemDelete struct {
	ItemId *uuid.UUID   `json:"itemId"`
	TagIds []*uuid.UUID `json:"tagId"`
}
