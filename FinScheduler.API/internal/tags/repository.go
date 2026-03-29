package tags

import (
	"context"
	"finscheduler/internal/metrics"
	"finscheduler/internal/shared"
	"finscheduler/internal/traces"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"log/slog"
	"time"
)

type TagsRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

const databaseDriver string = "postgresql"
const tableName = "tags"

func NewTagsRepository(db *sqlx.DB, logger *slog.Logger) *TagsRepository {
	return &TagsRepository{db: db, logger: logger}
}

func (repository *TagsRepository) Get(ctx context.Context, filter *TagFilter) ([]Tag, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tags []Tag
	var count int64 = 0

	var query = " FROM public.tags WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND id IN (?)", filter.Ids)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"Ids\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
			traces.EnrichFailedRepositorySpanRead(span, err, count)
			return nil, 0, err
		}

		query += inQuery
		args = append(args, inArgs...)
	}

	if filter.Name != nil && len(*filter.Name) > 0 {
		query += " AND name ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.IsActive != nil {
		query += " AND is_active = ?"
		args = append(args, *filter.IsActive)
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
	err := repository.db.Select(&tags, selectQuery, selectArgs...)
	metrics.RecordDatabaseDuration(ctx, selectStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationSelect)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repository.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "query", countQuery, "args", countArgs)
	countStart := time.Now()
	err = repository.db.Get(&count, countQuery, countArgs...)
	metrics.RecordDatabaseDuration(ctx, countStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationCount)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationCount)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(tags)))
	return tags, count, err
}

func (repository *TagsRepository) GetById(ctx context.Context, id uuid.UUID) (*Tag, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tag Tag

	if id == uuid.Nil {
		repository.logger.ErrorContext(ctx, "id should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("id should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	query := "SELECT * FROM public.tags WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "id", id)
	start := time.Now()
	err := repository.db.Get(&tag, query, id)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)

	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationSelect)
	}

	traces.EnrichSuccessRepositorySpanRead(span, 1)
	return &tag, nil
}

func (repository *TagsRepository) GetLookup(ctx context.Context, filter *TagFilter) ([]shared.Lookup, int64, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tags []shared.Lookup
	var count int64 = 0

	var query = " FROM public.tags WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Name != nil && len(*filter.Name) > 0 {
		query += " AND name ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	query += " AND is_active = true"

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
	err := repository.db.Select(&tags, selectQuery, selectArgs...)
	metrics.RecordDatabaseDuration(ctx, selectStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationSelect)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repository.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "query", countQuery, "args", countArgs)
	countStart := time.Now()
	err = repository.db.Get(&count, countQuery, countArgs...)
	metrics.RecordDatabaseDuration(ctx, countStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationCount)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationCount)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(tags)))
	return tags, count, err
}

func (repository *TagsRepository) Create(ctx context.Context, create *TagCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationInsert)
	defer span.End()

	newID, err := uuid.NewV7()

	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	query := "INSERT INTO public.tags (id, name, is_active) VALUES (?, ?, ?)"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	res, err := repository.db.Exec(query, newID, create.Name, create.IsActive)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tableName, err != nil, metrics.DatabaseOperationInsert)
	var affected int64 = 0
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "newID",
			newID, "name", create.Name, "isActive", create.IsActive)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	} else {
		affected, _ = res.RowsAffected()
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationInsert)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return newID, err
}

func (repository *TagsRepository) Update(ctx context.Context, tagID uuid.UUID, update *TagUpdate) (bool, error) {
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(ctx, "tags-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	transaction, err := repository.db.Beginx()
	defer func() {
		if p := recover(); p != nil {
			if rbErr := transaction.Rollback(); rbErr != nil {
				repository.logger.ErrorContext(ctx, "rollback failed", "error", rbErr)
			}
			panic(p)
		} else if err != nil {
			if rbErr := transaction.Rollback(); rbErr != nil {
				repository.logger.ErrorContext(ctx, "rollback failed", "error", rbErr)
			}
		}
	}()

	if err != nil {
		repository.logger.ErrorContext(ctx, "transaction error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := "UPDATE public.tags SET name = ?, is_active = ? WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "updating an tag:", "id", tagID, "name", update.Name,
		"isActive", update.IsActive)
	updateStart := time.Now()
	result, err := transaction.Exec(query, update.Name, update.IsActive, tagID)
	metrics.RecordDatabaseDuration(ctx, updateStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "id", tagID, "name",
			update.Name, "isActive", update.IsActive)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit tag", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationUpdate)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return success, err
}
