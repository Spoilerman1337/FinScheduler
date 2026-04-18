package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type TagsService struct {
	uow    *persistence.UnitOfWork
	logger *slog.Logger
}

func NewTagsService(uow *persistence.UnitOfWork, logger *slog.Logger) *TagsService {
	return &TagsService{
		uow:    uow,
		logger: logger,
	}
}

func (service *TagsService) Get(ctx context.Context, filter *domains.TagFilter) ([]domains.TagDto, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-service")
	traces.RecordServiceSpan(span, "Get")
	defer span.End()

	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		err := fmt.Errorf("filter is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	var tags []domains.TagDto
	var count int64

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawTags, rawTagsCount, err := repositories.Tags.Get(ctx, filter)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		count = rawTagsCount

		tags = make([]domains.TagDto, 0)
		if rawTags != nil && len(rawTags) > 0 {
			for _, tag := range rawTags {
				tags = append(tags, *domains.NewTagDto(tag))
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return tags, count, err
}

func (service *TagsService) GetById(ctx context.Context, id uuid.UUID) (*domains.TagDto, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-service")
	traces.RecordServiceSpan(span, "GetById")
	defer span.End()

	if id == uuid.Nil {
		service.logger.ErrorContext(ctx, "id is nil")
		err := fmt.Errorf("id is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return nil, err
	}

	var tag *domains.TagDto

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawTag, err := repositories.Tags.GetById(ctx, id)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		tag = domains.NewTagDto(*rawTag)

		return nil
	})
	if err != nil {
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return tag, err
}

func (service *TagsService) GetLookup(ctx context.Context, filter *domains.TagFilter) ([]domains.Lookup, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-service")
	traces.RecordServiceSpan(span, "Get")
	defer span.End()

	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		err := fmt.Errorf("filter is nil")
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	var tags []domains.Lookup
	var count int64

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawTags, rawTagsCount, err := repositories.Tags.GetLookup(ctx, filter)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			return err
		}

		tags = rawTags
		count = rawTagsCount

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return tags, count, err
}

func (service *TagsService) Create(ctx context.Context, create *domains.TagCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-service")
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
		newId, err = repositories.Tags.Create(ctx, create)

		if err != nil || newId == uuid.Nil {
			if err == nil {
				err = fmt.Errorf("failed to create tag: repository returned nil uuid")
			}
			return err
		}

		return nil
	})

	if err != nil || newId == uuid.Nil {
		if err == nil {
			err = fmt.Errorf("failed to create tag: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error creating a tag", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
	}

	traces.EnrichSuccessServiceSpan(span)

	return newId, err
}

func (service *TagsService) Update(ctx context.Context, tagID uuid.UUID, update *domains.TagUpdate) (bool, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-service")
	traces.RecordServiceSpan(span, "Update")
	defer span.End()

	if tagID == uuid.Nil {
		service.logger.ErrorContext(ctx, "tagID is nil")
		err := fmt.Errorf("tagID is nil")
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
		success, err = repositories.Tags.Update(ctx, tagID, update)

		if err != nil || !success {
			if err == nil {
				err = fmt.Errorf("failed to update tag: repository returned nil uuid")
			}
			return err
		}

		return nil
	})

	if err != nil || !success {
		if err == nil {
			err = fmt.Errorf("failed to update tag: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error updating a tag", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, err
}
