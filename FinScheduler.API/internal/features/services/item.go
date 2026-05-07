package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/persistence"
	"finscheduler/internal/traces"
	"finscheduler/pkg/dh"
	"finscheduler/pkg/rh"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type ItemsService struct {
	uow    *persistence.UnitOfWork
	logger *slog.Logger
}

const itemsServiceName = "items"

func NewItemsService(uow *persistence.UnitOfWork, logger *slog.Logger) *ItemsService {
	return &ItemsService{
		uow:    uow,
		logger: logger,
	}
}

func (service *ItemsService) Get(ctx context.Context, filter *domains.ItemFilter) ([]domains.ItemDto, int64, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Get")
	defer span.End()

	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		err := fmt.Errorf("filter is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Get", err)
		return nil, 0, err
	}

	var items []domains.ItemDto
	var count int64

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawItems, rawItemsCount, err := repositories.Items.Get(ctx, filter)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "Get", err)
			return err
		}

		count = rawItemsCount

		if len(rawItems) == 0 {
			items = make([]domains.ItemDto, 0)
			return nil
		}

		var itemIds []uuid.UUID
		for _, item := range rawItems {
			itemIds = append(itemIds, item.Id)
		}

		rawTagToItems, err := repositories.TagToItems.GetByItemIds(ctx, itemIds)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "Get", err)
			return err
		}

		tagIds := make([]uuid.UUID, len(rawTagToItems))
		if len(rawTagToItems) > 0 {
			for i, tag := range rawTagToItems {
				tagIds[i] = tag.TagId
			}
		}

		rawTags, err := repositories.Tags.GetByIds(ctx, tagIds)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "Get", err)
			return err
		}

		tagsById := make(map[uuid.UUID]domains.Tag, len(rawTags))
		for _, tag := range rawTags {
			tagsById[tag.Id] = tag
		}

		itemsWithTags := make(map[uuid.UUID][]domains.Tag)
		for _, tti := range rawTagToItems {
			if tag, exists := tagsById[tti.TagId]; exists {
				itemsWithTags[tti.ItemId] = append(itemsWithTags[tti.ItemId], tag)
			}
		}

		items = make([]domains.ItemDto, 0)
		for _, item := range rawItems {
			items = append(items, *domains.NewItemDto(item, itemsWithTags[item.Id]))
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return items, count, err
}

func (service *ItemsService) Create(ctx context.Context, create *domains.ItemCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Create")
	defer span.End()

	if create == nil {
		service.logger.ErrorContext(ctx, "create is nil")
		err := fmt.Errorf("create is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return uuid.Nil, err
	}

	if err := create.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "create validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return uuid.Nil, err
	}

	var newId uuid.UUID
	updateTagIds := parseTagIDs(create.TagIds)

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error

		newId, err = repositories.Items.Create(ctx, create)
		if err != nil {
			return err
		}
		if newId == uuid.Nil {
			return fmt.Errorf("failed to create item: repository returned nil uuid")
		}

		tagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{newId})
		if err != nil {
			return err
		}

		var currentTagIds []uuid.UUID
		if tagToItems != nil && len(tagToItems) > 0 {
			for _, tagToItem := range tagToItems {
				currentTagIds = append(currentTagIds, tagToItem.TagId)
			}
		}

		toDelete, toInsert := dh.Reconcile(updateTagIds, currentTagIds)

		if len(toDelete) > 0 {
			success, err := repositories.TagToItems.BulkDelete(ctx, &domains.TagToItemDelete{ItemId: &newId, TagIds: rh.ReferenceSlice(toDelete)})
			if err != nil {
				return err
			}
			if !success {
				return fmt.Errorf("failed to create item: tag to item delete affected no rows")
			}
		}

		if len(toInsert) > 0 {
			success, err := repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: &newId, TagIds: rh.ReferenceSlice(toInsert)})
			if err != nil {
				if details, ok := dh.GetPostgresErrorDetails(err); ok && details.Code == dh.PostgresForeignKeyViolationCode {
					return domains.ErrInvalidReference
				}
				return err
			}
			if !success {
				return fmt.Errorf("failed to create item: tag to item insert affected no rows")
			}
		}

		return nil
	})

	if err != nil || newId == uuid.Nil {
		if err == nil {
			err = fmt.Errorf("failed to create item: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error creating an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return newId, err
	}

	traces.EnrichSuccessServiceSpan(span)

	return newId, err
}

func (service *ItemsService) Update(ctx context.Context, itemID uuid.UUID, update *domains.ItemUpdate) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Update")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}
	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}

	if err := update.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "update validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}

	var success bool
	updateTagIds := parseTagIDs(update.TagIds)

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error

		success, err = repositories.Items.Update(ctx, itemID, update)
		if err != nil {
			return err
		}
		if !success {
			return nil
		}

		tagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{itemID})
		if err != nil {
			return err
		}

		var currentTagIds []uuid.UUID
		if tagToItems != nil && len(tagToItems) > 0 {
			for _, tagToItem := range tagToItems {
				currentTagIds = append(currentTagIds, tagToItem.TagId)
			}
		}

		toDelete, toInsert := dh.Reconcile(updateTagIds, currentTagIds)

		if len(toDelete) > 0 {
			tagSuccess, err := repositories.TagToItems.BulkDelete(ctx, &domains.TagToItemDelete{ItemId: &itemID, TagIds: rh.ReferenceSlice(toDelete)})
			if err != nil {
				return err
			}
			if !tagSuccess {
				return fmt.Errorf("failed to update item: tag to item delete affected no rows")
			}
		}

		if len(toInsert) > 0 {
			tagSuccess, err := repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: &itemID, TagIds: rh.ReferenceSlice(toInsert)})
			if err != nil {
				if details, ok := dh.GetPostgresErrorDetails(err); ok && details.Code == dh.PostgresForeignKeyViolationCode {
					return domains.ErrInvalidReference
				}
				return err
			}
			if !tagSuccess {
				return fmt.Errorf("failed to update item: tag to item insert affected no rows")
			}
		}

		return nil
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "error updating an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
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
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Delete", err)
		return false, err
	}

	var success bool

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error
		success, err = repositories.Items.Delete(ctx, itemID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "error deleting an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Delete", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
}

func parseTagIDs(tagIDs []string) []uuid.UUID {
	if tagIDs == nil {
		return nil
	}

	result := make([]uuid.UUID, len(tagIDs))
	for i, tagID := range tagIDs {
		result[i] = uuid.MustParse(tagID)
	}

	return result
}
