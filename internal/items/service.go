package items

import (
	"context"
	"finscheduler/internal/traces"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Get")
	defer span.End()

	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		err := fmt.Errorf("filter is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	rawItems, count, err := service.repository.Get(ctx, filter)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get items failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	items := make([]ItemDto, 0)
	if rawItems != nil && len(rawItems) > 0 {
		for _, item := range rawItems {
			items = append(items, *NewItemDto(&item))
		}
	}

	traces.EnrichSuccessServiceSpan(span)
	return items, count, err
}

func (service *ItemsService) GetById(ctx context.Context, id uuid.UUID) (*ItemDto, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "GetById")
	defer span.End()

	if id == uuid.Nil {
		service.logger.ErrorContext(ctx, "id is nil")
		err := fmt.Errorf("id is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return nil, err
	}

	rawItem, err := service.repository.GetById(ctx, id)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get items failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return NewItemDto(rawItem), err
}

func (service *ItemsService) Create(ctx context.Context, create *ItemCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Create")
	defer span.End()

	if create == nil {
		service.logger.ErrorContext(ctx, "create is nil")
		err := fmt.Errorf("create is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return uuid.Nil, err
	}

	newId, err := service.repository.Create(ctx, create)

	if err != nil || newId == uuid.Nil {
		service.logger.ErrorContext(ctx, "error creating an item", "error", err)
		span.SetStatus(codes.Error, "")
		if err == nil {
			err = fmt.Errorf("error creating an item")
		}
		span.RecordError(err)
	}

	traces.EnrichSuccessServiceSpan(span)
	return newId, err
}

func (service *ItemsService) Update(ctx context.Context, itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Update")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return false, err
	}
	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return false, err
	}

	success, err := service.repository.Update(ctx, itemID, update)

	if err != nil || !success {
		service.logger.ErrorContext(ctx, "error updating an item", "error", err)
		if err == nil {
			err = fmt.Errorf("error updating an item")
		}
		traces.EnrichFailedServiceSpan(span, err)
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, err
}

func (service *ItemsService) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Delete")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return false, err
	}

	success, err := service.repository.Delete(ctx, itemID)

	if err != nil || !success {
		service.logger.ErrorContext(ctx, "error deleting an item", "error", err)
		if err == nil {
			err = fmt.Errorf("error updating an item")
		}
		traces.EnrichFailedServiceSpan(span, err)
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, err
}
