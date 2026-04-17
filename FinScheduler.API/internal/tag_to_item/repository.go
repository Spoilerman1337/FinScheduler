package tag_to_item

import (
	"context"
	"database/sql"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"finscheduler/pkg/rh"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type TagToItemsRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewTagToItemsRepository(db *sqlx.DB, logger *slog.Logger) *TagToItemsRepository {
	return &TagToItemsRepository{db: db, logger: logger}
}

const databaseDriver string = "postgresql"
const tagsToItemTableName = "tag_to_item"

func (repository *TagToItemsRepository) GetByItemIds(ctx context.Context, itemIds []uuid.UUID) ([]TagToItem, error) {
	tracer := otel.Tracer("tag-to-items")
	ctx, span := tracer.Start(ctx, "tag-to-items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var tagToItems []TagToItem

	if itemIds == nil {
		repository.logger.ErrorContext(ctx, "itemId should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("itemId should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	if len(itemIds) == 0 {
		return make([]TagToItem, 0), nil
	}

	query := `SELECT tag_id
			  FROM public.tag_to_item tti 
			  WHERE tti.item_id IN (?)`
	query, inArgs, err := sqlx.In(query, itemIds)
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "itemIds", itemIds)
	start := time.Now()
	err = repository.db.SelectContext(ctx, &tagToItems, query, inArgs)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tagsToItemTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, true, metrics.DatabaseOperationSelect)
	}

	traces.EnrichSuccessRepositorySpanRead(span, 1)

	return tagToItems, nil
}

func (repository *TagToItemsRepository) BulkInsert(ctx context.Context, create *TagToItemCreate) (bool, error) {
	tracer := otel.Tracer("tag-to-items")
	ctx, span := tracer.Start(ctx, "tag-to-items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
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

	args := make([]interface{}, 0, len(create.TagIds)*2)
	values := make([]string, 0, len(create.TagIds))

	for _, tagId := range create.TagIds {
		values = append(values, "(?, ?)")
		args = append(args, create.ItemId, tagId)
	}

	query := fmt.Sprintf("INSERT INTO public.tag_to_item (item_id, tag_id) VALUES %s",
		strings.Join(values, ","))
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	res, err := transaction.ExecContext(ctx, query, args...)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tagsToItemTableName, err != nil, metrics.DatabaseOperationInsert)
	var affected int64 = 0
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "itemId",
			create.ItemId, "tagIds", create.TagIds)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	} else {
		affected, _ = res.RowsAffected()
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, true, metrics.DatabaseOperationInsert)
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return affected > 0, err
}

func (repository *TagToItemsRepository) BulkDelete(ctx context.Context, delete *TagToItemDelete) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationDelete)
	defer span.End()

	transaction, err := repository.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
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
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := `DELETE FROM public.tag_to_item WHERE item_id = ? AND tag_id IN (?)`
	query, inArgs, err := sqlx.In(query, *delete.ItemId, rh.DereferenceSlice(delete.TagIds))
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "fetching delete tag to items:", "query", query, "itemId", delete.ItemId, "tagIds", delete.TagIds)
	start := time.Now()
	result, err := transaction.ExecContext(ctx, query, inArgs...)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tagsToItemTableName, err != nil, metrics.DatabaseOperationDelete)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "query", query, "itemId", delete.ItemId, "tagIds", delete.TagIds)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, true, metrics.DatabaseOperationDelete)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tagsToItemTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return success, err
}
