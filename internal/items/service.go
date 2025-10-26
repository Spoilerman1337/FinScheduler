package items

import (
	"fmt"
	"github.com/google/uuid"
)

type ItemsService struct {
	repo *ItemsRepository
}

func NewItemsService(repo *ItemsRepository) *ItemsService {
	return &ItemsService{repo: repo}
}

func (service *ItemsService) Get(filter *ItemFilter) ([]ItemDto, int64, error) {
	if filter == nil {
		return nil, 0, fmt.Errorf("filter is nil")
	}

	rawItems, count, err := service.repo.Get(filter)

	items := make([]ItemDto, 0)
	if rawItems != nil && len(rawItems) > 0 {
		for _, item := range rawItems {
			items = append(items, *NewItemDto(&item))
		}
	}

	return items, count, err
}

func (service *ItemsService) GetById(id uuid.UUID) (*ItemDto, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf(`id is nil`)
	}

	rawItem, err := service.repo.GetById(id)

	return NewItemDto(rawItem), err
}

func (service *ItemsService) Create(create *ItemCreate) (uuid.UUID, error) {
	if create == nil {
		return uuid.Nil, fmt.Errorf(`item is nil`)
	}

	return service.repo.Create(create)
}

func (service *ItemsService) Update(itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	if itemID == uuid.Nil {
		return false, fmt.Errorf(`itemID is nil`)
	}
	if update == nil {
		return false, fmt.Errorf(`update is nil`)
	}

	return service.repo.Update(itemID, update)
}

func (service *ItemsService) Delete(itemID uuid.UUID) (bool, error) {
	if itemID == uuid.Nil {
		return false, fmt.Errorf(`itemID is nil`)
	}

	return service.repo.Delete(itemID)
}
