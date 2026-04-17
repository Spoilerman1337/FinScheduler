package items

import (
	"context"
	"database/sql"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"finscheduler/pkg/rh"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type ItemsRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

const databaseDriver string = "postgresql"
const itemsTableName = "items"

func NewItemsRepository(db *sqlx.DB, logger *slog.Logger) *ItemsRepository {
	return &ItemsRepository{db: db, logger: logger}
}

func (repository *ItemsRepository) Get(ctx context.Context, filter *ItemFilter) ([]Item, int64, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var items []*Item
	var count int64 = 0

	var itemsQuery = " FROM public.items i "

	//TODO: bruh. В Go наверняка есть что-то вроде сишарписткого string.Join(',', []). Собрать параметры именно так, и не писать кринж типа 1 = 1
	itemsQuery += " WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND id IN (?)", filter.Ids)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"Ids\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
			traces.EnrichFailedRepositorySpanRead(span, err, count)
			return nil, 0, err
		}

		itemsQuery += inQuery
		args = append(args, inArgs...)
	}

	if filter.Name != nil && len(*filter.Name) > 0 {
		itemsQuery += " AND i.name ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.PriceFrom != nil {
		itemsQuery += " AND i.price >= ?"
		args = append(args, *filter.PriceFrom)
	}

	if filter.PriceTo != nil {
		itemsQuery += " AND i.price <= ?"
		args = append(args, *filter.PriceTo)
	}

	if filter.Description != nil && len(*filter.Description) > 0 {
		itemsQuery += " AND i.description ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Description))
	}

	if filter.IsActive != nil {
		itemsQuery += " AND i.is_active = ?"
		args = append(args, *filter.IsActive)
	}

	if filter.CreatedFrom != nil {
		itemsQuery += " AND i.created_at >= ?"
		args = append(args, *filter.CreatedFrom)
	}

	if filter.CreatedTo != nil {
		itemsQuery += " AND i.created_at <= ?"
		args = append(args, *filter.CreatedTo)
	}

	if filter.UpdatedFrom != nil {
		itemsQuery += " AND i.updated_at >= ?"
		args = append(args, *filter.UpdatedFrom)
	}

	if filter.UpdatedTo != nil {
		itemsQuery += " AND i.updated_at <= ?"
		args = append(args, *filter.UpdatedTo)
	}

	if filter.CashbackFrom != nil {
		itemsQuery += " AND i.cashback >= ?"
		args = append(args, *filter.CashbackFrom)
	}

	if filter.CashbackTo != nil {
		itemsQuery += " AND i.cashback <= ?"
		args = append(args, *filter.CashbackTo)
	}

	if filter.Categories != nil && len(filter.Categories) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND i.category IN (?)", filter.Categories)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"Categories\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
			traces.EnrichFailedRepositorySpanRead(span, err, count)
			return nil, 0, err
		}

		itemsQuery += inQuery
		args = append(args, inArgs...)
	}

	if filter.TagIds != nil && len(filter.TagIds) > 0 {
		inQuery, inArgs, err := sqlx.In(` AND EXISTS (
			SELECT 1 FROM public.tag_to_item tti 
			WHERE tti.item_id = i.id AND tti.tag_id IN (?)
		)`, filter.TagIds)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"TagIds\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
			traces.EnrichFailedRepositorySpanRead(span, err, count)
			return nil, 0, err
		}

		itemsQuery += inQuery
		args = append(args, inArgs...)
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

	itemsSelectQuery := fmt.Sprintf("SELECT * %s LIMIT ? OFFSET ?", itemsQuery)
	itemsSelectQuery = repository.db.Rebind(itemsSelectQuery)
	itemsSelectArgs := append(make([]interface{}, 0), args...)
	itemsSelectArgs = append(itemsSelectArgs, pageSize, offset)

	repository.logger.InfoContext(ctx, "executing operation:", "itemsQuery", itemsSelectQuery, "args", itemsSelectArgs)
	itemsSelectStart := time.Now()
	err := repository.db.SelectContext(ctx, &items, itemsSelectQuery, itemsSelectArgs...)
	metrics.RecordDatabaseDuration(ctx, itemsSelectStart, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationSelect)
	}

	itemsCountQuery := fmt.Sprintf("SELECT COUNT(*) %s", itemsQuery)
	itemsCountQuery = repository.db.Rebind(itemsCountQuery)
	itemsCountArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "itemsQuery", itemsCountQuery, "args", itemsCountArgs)
	itemsCountStart := time.Now()
	err = repository.db.Get(&count, itemsCountQuery, itemsCountArgs...)
	metrics.RecordDatabaseDuration(ctx, itemsCountStart, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationCount)
		traces.EnrichFailedRepositorySpanRead(span, err, count)
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationCount)
	}

	traces.EnrichSuccessRepositorySpanRead(span, int64(len(items)))
	return rh.DereferenceSlice(items), count, err
}

func (repository *ItemsRepository) GetById(ctx context.Context, id uuid.UUID) (*Item, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	var item Item

	if id == uuid.Nil {
		repository.logger.ErrorContext(ctx, "id should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("id should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	query := "SELECT * FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "id", id)
	start := time.Now()
	err := repository.db.Get(&item, query, id)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationSelect)

	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationSelect)
	}

	traces.EnrichSuccessRepositorySpanRead(span, 1)
	return &item, nil
}

func (repository *ItemsRepository) Create(ctx context.Context, create *ItemCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationInsert)
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

	newID, err := uuid.NewV7()
	now := time.Now().UTC()

	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	query := "INSERT INTO public.items (id, name, price, description, is_active, created_at, cashback, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	res, err := transaction.ExecContext(ctx, query, newID, create.Name, create.Price, create.Description, create.IsActive, now, create.Cashback, create.Category)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationInsert)
	var affected int64 = 0
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "newID",
			newID, "name", create.Name, "price", create.Price, "description", create.Description, "isActive",
			create.IsActive, "createdAt", now, "cashback", create.Cashback, "category", create.Category)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	} else {
		affected, _ = res.RowsAffected()
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationInsert)
	}

	err = transaction.Commit()

	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return newID, err
}

func (repository *ItemsRepository) Update(ctx context.Context, itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-repository")
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

	now := time.Now().UTC()

	if err != nil {
		repository.logger.ErrorContext(ctx, "transaction error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := "UPDATE public.items SET name = ?, price = ?, description = ?, is_active = ?, updated_at = ?, cashback = ?, category = ? WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "updating an item:", "id",
		itemID, "name", update.Name, "price", update.Price, "description", update.Description, "isActive",
		update.IsActive, "updatedAt", now, "cashback", update.Cashback, "category", update.Category)
	updateStart := time.Now()
	result, err := transaction.ExecContext(ctx, query, update.Name, update.Price, update.Description, update.IsActive,
		sql.NullTime{Time: now, Valid: true}, update.Cashback, update.Category, itemID)
	metrics.RecordDatabaseDuration(ctx, updateStart, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "id",
			itemID, "name", update.Name, "price", update.Price, "description", update.Description, "isActive",
			update.IsActive, "updatedAt", now, "cashback", update.Cashback, "category", update.Category)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationUpdate)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return success, err
}

func (repository *ItemsRepository) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
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
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := "DELETE FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "fetching delete item:", "query", query, "id", itemID)
	start := time.Now()
	result, err := transaction.ExecContext(ctx, query, itemID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, itemsTableName, err != nil, metrics.DatabaseOperationDelete)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "error", err, "id", itemID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, true, metrics.DatabaseOperationDelete)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, itemsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
	}

	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return success, err
}
