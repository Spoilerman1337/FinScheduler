package services

import (
	"context"
	"finscheduler/internal/features/domains"
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

const maxPageSize int32 = 1<<31 - 1

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
		return nil, 0, err
	}

	var items []domains.ItemDto
	var count int64

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawItems, rawItemsCount, err := repositories.Items.Get(ctx, filter)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		count = rawItemsCount

		var itemIds []uuid.UUID
		if rawItems != nil {
			for _, item := range rawItems {
				itemIds = append(itemIds, item.Id)
			}
		}

		//TODO: Возможно, стоит сделать отдельные методы в репозиториях для получения тегов/тти по ID итема, и без пагинации, чтобы не делать вот это вот.
		pageSize := maxPageSize
		page := int32(0)

		rawTagToItems, err := repositories.TagToItems.GetByItemIds(ctx, itemIds)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		tagIds := make([]*uuid.UUID, len(rawTagToItems))
		if len(rawTagToItems) > 0 {
			for i, tag := range rawTagToItems {
				tagIds[i] = &tag.TagId
			}
		}

		rawTags, _, err := repositories.Tags.Get(ctx, &domains.TagFilter{PageSize: &pageSize, Page: &page, Ids: tagIds})
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
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

func (service *ItemsService) GetById(ctx context.Context, id uuid.UUID) (*domains.ItemDto, error) {
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

	var item *domains.ItemDto

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawItem, err := repositories.Items.GetById(ctx, id)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		//TODO: Возможно, стоит сделать отдельные методы в репозиториях для получения тегов/тти по ID итема, и без пагинации, чтобы не делать вот это вот.
		pageSize := maxPageSize
		page := int32(0)

		rawTagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{id})
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		tagIds := make([]*uuid.UUID, len(rawTagToItems))
		if len(rawTagToItems) > 0 {
			for i, tag := range rawTagToItems {
				tagIds[i] = &tag.TagId
			}
		}

		rawTags, _, err := repositories.Tags.Get(ctx, &domains.TagFilter{PageSize: &pageSize, Page: &page, Ids: tagIds})
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		item = domains.NewItemDto(*rawItem, rawTags)

		return nil
	})
	if err != nil {
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return item, nil
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
		return uuid.Nil, err
	}

	var newId uuid.UUID

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error

		newId, err = repositories.Items.Create(ctx, create)

		tagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{newId})

		var currentTagIds []uuid.UUID
		if tagToItems != nil && len(tagToItems) > 0 {
			for _, tagToItem := range tagToItems {
				currentTagIds = append(currentTagIds, tagToItem.TagId)
			}
		}

		var updateTagIds []uuid.UUID
		if create.TagIds != nil && len(create.TagIds) > 0 {
			for _, tagId := range create.TagIds {
				uuidTagId, err := uuid.Parse(tagId)
				if err == nil {
					updateTagIds = append(updateTagIds, uuidTagId)
				}
			}
		}

		toDelete, toInsert := dh.Reconcile(updateTagIds, currentTagIds)

		if len(toDelete) > 0 {
			_, err = repositories.TagToItems.BulkDelete(ctx, &domains.TagToItemDelete{ItemId: &newId, TagIds: rh.ReferenceSlice(toDelete)})
		}

		if len(toInsert) > 0 {
			_, err = repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: &newId, TagIds: rh.ReferenceSlice(toInsert)})
		}

		if err != nil || newId == uuid.Nil {
			if err == nil {
				err = fmt.Errorf("failed to create item: repository returned nil uuid")
			}
			return err
		}

		return nil
	})

	if err != nil || newId == uuid.Nil {
		if err == nil {
			err = fmt.Errorf("failed to create item: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error creating an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
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
		return false, err
	}
	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return false, err
	}

	var success bool

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error

		success, err = repositories.Items.Update(ctx, itemID, update)

		tagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{itemID})

		var currentTagIds []uuid.UUID
		if tagToItems != nil && len(tagToItems) > 0 {
			for _, tagToItem := range tagToItems {
				currentTagIds = append(currentTagIds, tagToItem.TagId)
			}
		}

		var updateTagIds []uuid.UUID
		if update.TagIds != nil && len(update.TagIds) > 0 {
			for _, tagId := range update.TagIds {
				uuidTagId, err := uuid.Parse(tagId)
				if err == nil {
					updateTagIds = append(updateTagIds, uuidTagId)
				}
			}
		}

		toDelete, toInsert := dh.Reconcile(updateTagIds, currentTagIds)

		if len(toDelete) > 0 {
			success, err = repositories.TagToItems.BulkDelete(ctx, &domains.TagToItemDelete{ItemId: &itemID, TagIds: rh.ReferenceSlice(toDelete)})
		}

		if len(toInsert) > 0 {
			success, err = repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: &itemID, TagIds: rh.ReferenceSlice(toInsert)})
		}

		if err != nil || !success {
			if err == nil {
				err = fmt.Errorf("failed to update item: repository returned nil uuid")
			}
			return err
		}

		return nil
	})

	if err != nil || !success {
		if err == nil {
			err = fmt.Errorf("failed to update item: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error updating an item", "error", err)
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

	var success bool

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error
		success, err = repositories.Items.Delete(ctx, itemID)

		if err != nil || !success {
			if err == nil {
				err = fmt.Errorf("failed to delete item: repository returned nil uuid")
			}
			return err
		}

		return nil
	})

	if err != nil || !success {
		if err == nil {
			err = fmt.Errorf("failed to delete item: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error deleting an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, err
}
