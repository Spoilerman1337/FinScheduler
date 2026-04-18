package repositories

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type TagsRepository struct {
	db     DBTX
	logger *slog.Logger
}

func NewTagsRepository(db DBTX, logger *slog.Logger) *TagsRepository {
	return &TagsRepository{db: db, logger: logger}
}

func (repository *TagsRepository) Get(ctx context.Context, filter *domains.TagFilter) ([]domains.Tag, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tags []domains.Tag
	var count int64 = 0

	query := "FROM public.tags"
	filters := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In("id IN (?)", filter.Ids)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"Ids\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationNone)
			traces.EnrichFailedRepositorySpanRead(span, err, count)
			return nil, 0, err
		}

		filters = append(filters, inQuery)
		args = append(args, inArgs...)
	}

	if filter.Name != nil && len(*filter.Name) > 0 {
		filters = append(filters, "name ILIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.IsActive != nil {
		filters = append(filters, "is_active = ?")
		args = append(args, *filter.IsActive)
	}

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	var pageSize int32 = 20
	if filter.PageSize != nil {
		pageSize = *filter.PageSize
	}
	var page int32 = 0
	if filter.Page != nil {
		page = *filter.Page
	}
	offset := page * pageSize

	selectQuery := fmt.Sprintf("SELECT * %s LIMIT ? OFFSET ?", query)
	selectQuery = repository.db.Rebind(selectQuery)
	selectArgs := append(make([]interface{}, 0), args...)
	selectArgs = append(selectArgs, pageSize, offset)

	repository.logger.InfoContext(ctx, "executing operation:", "query", selectQuery, "args", selectArgs)
	selectStart := time.Now()
	err := sqlx.SelectContext(ctx, repository.db, &tags, selectQuery, selectArgs...)
	metrics.RecordDatabaseDuration(ctx, selectStart, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationSelect)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repository.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "query", countQuery, "args", countArgs)
	countStart := time.Now()
	err = sqlx.GetContext(ctx, repository.db, &count, countQuery, countArgs...)
	metrics.RecordDatabaseDuration(ctx, countStart, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationCount)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationCount)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(tags)))
	return tags, count, err
}

func (repository *TagsRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]domains.Tag, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tags []domains.Tag

	if ids == nil {
		repository.logger.ErrorContext(ctx, "ids should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("ids should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	if len(ids) == 0 {
		return make([]domains.Tag, 0), nil
	}

	query := "SELECT * FROM public.tags WHERE id IN (?)"
	query, inArgs, err := sqlx.In(query, ids)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error binding \"Ids\" array to IN filter", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "ids", ids)
	start := time.Now()
	err = sqlx.SelectContext(ctx, repository.db, &tags, query, inArgs...)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationSelect)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(tags)))
	return tags, nil
}

func (repository *TagsRepository) GetLookup(ctx context.Context, filter *domains.TagLookupFilter) ([]domains.Lookup, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tags []domains.Lookup
	var count int64 = 0

	query := "FROM public.tags"
	filters := make([]string, 0)
	args := make([]interface{}, 0)

	if filter.Name != nil && len(*filter.Name) > 0 {
		filters = append(filters, "name ILIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	filters = append(filters, "is_active = true")

	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	var pageSize int32 = 20
	if filter.PageSize != nil {
		pageSize = *filter.PageSize
	}
	var page int32 = 0
	if filter.Page != nil {
		page = *filter.Page
	}
	offset := page * pageSize

	selectQuery := fmt.Sprintf("SELECT id as value, name as label %s LIMIT ? OFFSET ?", query)
	selectQuery = repository.db.Rebind(selectQuery)
	selectArgs := append(make([]interface{}, 0), args...)
	selectArgs = append(selectArgs, pageSize, offset)

	repository.logger.InfoContext(ctx, "executing operation:", "query", selectQuery, "args", selectArgs)
	selectStart := time.Now()
	err := sqlx.SelectContext(ctx, repository.db, &tags, selectQuery, selectArgs...)
	metrics.RecordDatabaseDuration(ctx, selectStart, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationSelect)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repository.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "query", countQuery, "args", countArgs)
	countStart := time.Now()
	err = sqlx.GetContext(ctx, repository.db, &count, countQuery, countArgs...)
	metrics.RecordDatabaseDuration(ctx, countStart, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationCount)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationCount)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(tags)))
	return tags, count, err
}

func (repository *TagsRepository) Create(ctx context.Context, create *domains.TagCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationInsert)
	defer span.End()

	newID, err := uuid.NewV7()

	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	query := "INSERT INTO public.tags (id, name, is_active) VALUES (?, ?, ?)"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	res, err := repository.db.ExecContext(ctx, query, newID, create.Name, create.IsActive)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationInsert)
	var affected int64 = 0
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "newID",
			newID, "name", create.Name, "isActive", create.IsActive)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	} else {
		affected, _ = res.RowsAffected()
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationInsert)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return newID, err
}

func (repository *TagsRepository) Update(ctx context.Context, tagID uuid.UUID, update *domains.TagUpdate) (bool, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	query := "UPDATE public.tags SET name = ?, is_active = ? WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "updating an tag:", "id", tagID, "name", update.Name,
		"isActive", update.IsActive)
	updateStart := time.Now()
	result, err := repository.db.ExecContext(ctx, query, update.Name, update.IsActive, tagID)
	metrics.RecordDatabaseDuration(ctx, updateStart, databaseDriver, tagsTableName, err != nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "id", tagID, "name",
			update.Name, "isActive", update.IsActive)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, true, metrics.DatabaseOperationUpdate)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return success, err
}
