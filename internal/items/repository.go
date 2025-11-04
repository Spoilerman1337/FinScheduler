package items

import (
	"context"
	"database/sql"
	"finscheduler/internal/metrics"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type ItemsRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

const databaseDriver string = "postgres"
const tableName = "items"

func NewItemsRepository(db *sqlx.DB, logger *slog.Logger) *ItemsRepository {
	return &ItemsRepository{db: db, logger: logger}
}

func (repository *ItemsRepository) Get(ctx context.Context, filter *ItemFilter) ([]Item, int64, error) {
	var items []Item
	var count int64 = 0

	var query = " FROM public.items WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND id IN (?)", filter.Ids)

		if err != nil {
			repository.logger.ErrorContext(ctx, "error binding \"Ids\" array to IN filter", "error", err)
			metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
			return nil, 0, err
		}

		query += inQuery
		args = append(args, inArgs...)
	}

	if filter.Name != nil && len(*filter.Name) > 0 {
		query += " AND name ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.PriceFrom != nil {
		query += " AND price >= ?"
		args = append(args, *filter.PriceFrom)
	}

	if filter.PriceTo != nil {
		query += " AND price <= ?"
		args = append(args, *filter.PriceTo)
	}

	if filter.Description != nil && len(*filter.Description) > 0 {
		query += " AND description ILIKE ?"
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Description))
	}

	if filter.IsActive != nil {
		query += " AND is_active = ?"
		args = append(args, *filter.IsActive)
	}

	if filter.CreatedFrom != nil {
		query += " AND created_at >= ?"
		args = append(args, *filter.CreatedFrom)
	}

	if filter.CreatedTo != nil {
		query += " AND created_at <= ?"
		args = append(args, *filter.CreatedTo)
	}

	if filter.UpdatedFrom != nil {
		query += " AND updated_at >= ?"
		args = append(args, *filter.UpdatedFrom)
	}

	if filter.UpdatedTo != nil {
		query += " AND updated_at <= ?"
		args = append(args, *filter.UpdatedTo)
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
	err := repository.db.Select(&items, selectQuery, selectArgs...)
	metrics.RecordDatabaseDuration(ctx, selectStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationSelect)
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
		return nil, 0, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationCount)
	}

	return items, count, err
}

func (repository *ItemsRepository) GetById(ctx context.Context, id uuid.UUID) (*Item, error) {
	var item Item

	if id == uuid.Nil {
		repository.logger.ErrorContext(ctx, "id should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
		return nil, fmt.Errorf("id should not be nil")
	}

	query := "SELECT * FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "id", id)
	start := time.Now()
	err := repository.db.Get(&item, query, id)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tableName, err != nil, metrics.DatabaseOperationSelect)

	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationSelect)
		return nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationSelect)
	}

	return &item, nil
}

func (repository *ItemsRepository) Create(ctx context.Context, create *ItemCreate) (uuid.UUID, error) {
	newID, err := uuid.NewV7()
	now := time.Now().UTC()

	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
		return uuid.Nil, err
	}

	query := "INSERT INTO public.items (id, name, price, description, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	_, err = repository.db.Exec(query, newID, create.Name, create.Price, create.Description, create.IsActive, now)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tableName, err != nil, metrics.DatabaseOperationInsert)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "newID",
			newID, "name", create.Name, "price", create.Price, "description", create.Description, "isActive",
			create.IsActive, "createdAt", now)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationInsert)
		return uuid.Nil, err
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationInsert)
	}

	return newID, err
}

func (repository *ItemsRepository) Update(ctx context.Context, itemID uuid.UUID, update *ItemUpdate) (bool, error) {
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
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationNone)
		return false, err
	}

	query := "UPDATE public.items SET name = ?, price = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "updating an item:", "id",
		itemID, "name", update.Name, "price", update.Price, "description", update.Description, "isActive",
		update.IsActive, "updatedAt", now)
	updateStart := time.Now()
	result, err := transaction.Exec(query, update.Name, update.Price, update.Description, update.IsActive,
		sql.NullTime{Time: now, Valid: true}, itemID)
	metrics.RecordDatabaseDuration(ctx, updateStart, databaseDriver, tableName, err != nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "id",
			itemID, "name", update.Name, "price", update.Price, "description", update.Description, "isActive",
			update.IsActive, "updatedAt", now)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationUpdate)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationUpdate)
	}

	return success, err
}

func (repository *ItemsRepository) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
	transaction, err := repository.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
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
		return false, err
	}

	query := "DELETE FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "fetching delete item:", "query", query, "id", itemID)
	start := time.Now()
	result, err := transaction.Exec(query, itemID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, tableName, err != nil, metrics.DatabaseOperationDelete)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "error", err, "id", itemID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationDelete)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationDelete)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationDelete)
		return false, err
	}

	success := rowsAffected > 0
	if success {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, true, metrics.DatabaseOperationDelete)
	} else {
		metrics.RecordDatabaseRequest(ctx, databaseDriver, tableName, false, metrics.DatabaseOperationDelete)
	}

	return success, err
}
