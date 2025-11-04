package items

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
)

type ItemsService struct {
	repository *ItemsRepository
	logger     *slog.Logger
}

func NewItemsService(repository *ItemsRepository, logger *slog.Logger) *ItemsService {
	return &ItemsService{
		repository: repository,
		logger:     logger,
	}
}

func (service *ItemsService) Get(ctx context.Context, filter *ItemFilter) ([]ItemDto, int64, error) {
	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		return nil, 0, fmt.Errorf("filter is nil")
	}

	rawItems, count, err := service.repository.Get(ctx, filter)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get items failed", "error", err)
		return nil, 0, err
	}

	items := make([]ItemDto, 0)
	if rawItems != nil && len(rawItems) > 0 {
		for _, item := range rawItems {
			items = append(items, *NewItemDto(&item))
		}
	}

	return items, count, err
}

func (service *ItemsService) GetById(ctx context.Context, id uuid.UUID) (*ItemDto, error) {
	if id == uuid.Nil {
		service.logger.ErrorContext(ctx, "id is nil")
		return nil, fmt.Errorf(`id is nil`)
	}

	rawItem, err := service.repository.GetById(ctx, id)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get items failed", "error", err)
		return nil, err
	}

	return NewItemDto(rawItem), err
}

func (service *ItemsService) Create(ctx context.Context, create *ItemCreate) (uuid.UUID, error) {
	if create == nil {
		service.logger.ErrorContext(ctx, "create is nil")
		return uuid.Nil, fmt.Errorf(`create is nil`)
	}

	newId, err := service.repository.Create(ctx, create)

	if err != nil || newId == uuid.Nil {
		service.logger.ErrorContext(ctx, "error creating an item", "error", err)
	}

	return newId, err
}

func (service *ItemsService) Update(ctx context.Context, itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		return false, fmt.Errorf(`itemID is nil`)
	}
	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		return false, fmt.Errorf(`update is nil`)
	}

	success, err := service.repository.Update(ctx, itemID, update)

	if err != nil || !success {
		service.logger.ErrorContext(ctx, "error updating an item", "error", err)
	}

	return success, err
}

func (service *ItemsService) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		return false, fmt.Errorf(`itemID is nil`)
	}

	success, err := service.repository.Delete(ctx, itemID)

	if err != nil || !success {
		service.logger.ErrorContext(ctx, "error deleting an item", "error", err)
	}

	return success, err
}
