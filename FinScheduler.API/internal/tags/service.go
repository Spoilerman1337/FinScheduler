package tags

import (
	"context"
	"finscheduler/internal/shared"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type TagsService struct {
	repository *TagsRepository
	logger     *slog.Logger
}

func NewTagsService(repository *TagsRepository, logger *slog.Logger) *TagsService {
	return &TagsService{
		repository: repository,
		logger:     logger,
	}
}

func (service *TagsService) Get(ctx context.Context, filter *TagFilter) ([]TagDto, int64, error) {
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

	rawTags, count, err := service.repository.Get(ctx, filter)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	tags := make([]TagDto, 0)
	if rawTags != nil && len(rawTags) > 0 {
		for _, tag := range rawTags {
			tags = append(tags, *NewTagDto(tag))
		}
	}

	traces.EnrichSuccessServiceSpan(span)
	return tags, count, err
}

func (service *TagsService) GetById(ctx context.Context, id uuid.UUID) (*TagDto, error) {
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

	rawTag, err := service.repository.GetById(ctx, id)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return NewTagDto(*rawTag), err
}

func (service *TagsService) GetLookup(ctx context.Context, filter *TagFilter) ([]shared.Lookup, int64, error) {
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

	tags, count, err := service.repository.GetLookup(ctx, filter)
	if err != nil {
		service.logger.ErrorContext(ctx, "Get tags failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		return nil, 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return tags, count, err
}

func (service *TagsService) Create(ctx context.Context, create *TagCreate) (uuid.UUID, error) {
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

	newId, err := service.repository.Create(ctx, create)

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

func (service *TagsService) Update(ctx context.Context, tagID uuid.UUID, update *TagUpdate) (bool, error) {
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

	success, err := service.repository.Update(ctx, tagID, update)

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
