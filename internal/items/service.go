package items

import (
	"fmt"
	"github.com/google/uuid"
)

type ItemsService struct {
	repo *ItemsRepo
}

func NewItemsService(repo *ItemsRepo) *ItemsService {
	return &ItemsService{repo: repo}
}

func (service *ItemsService) Get(filter *ItemFilter) ([]Item, int64, error) {
	if filter == nil {
		return nil, 0, fmt.Errorf("filter is nil")
	}

	return service.repo.Get(filter)
}

func (service *ItemsService) GetById(id uuid.UUID) (*Item, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf(`id is nil`)
	}

	return service.repo.GetById(id)
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
