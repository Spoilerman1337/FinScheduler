package items

import (
	"context"
	"database/sql"
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
	err := repository.db.Select(&items, selectQuery, selectArgs...)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repository.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	repository.logger.InfoContext(ctx, "executing operation:", "query", countQuery, "args", countArgs)
	err = repository.db.Get(&count, countQuery, countArgs...)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		return nil, 0, err
	}

	return items, count, err
}

func (repository *ItemsRepository) GetById(ctx context.Context, id uuid.UUID) (*Item, error) {
	var item Item

	if id == uuid.Nil {
		repository.logger.ErrorContext(ctx, "id should not be nil")
		return nil, fmt.Errorf("id should not be nil")
	}

	query := "SELECT * FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "id", id)
	err := repository.db.Get(&item, query, id)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on COUNT operation", "error", err)
		return nil, err
	}

	return &item, nil
}

func (repository *ItemsRepository) Create(ctx context.Context, create *ItemCreate) (uuid.UUID, error) {
	newID, err := uuid.NewV7()
	now := time.Now().UTC()

	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		return uuid.Nil, err
	}

	query := "INSERT INTO public.items (id, name, price, description, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	_, err = repository.db.Exec(query, newID, create.Name, create.Price, create.Description, create.IsActive, now)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "newID",
			newID, "name", create.Name, "price", create.Price, "description", create.Description, "isActive",
			create.IsActive, "createdAt", now)
		return uuid.Nil, err
	}

	return newID, err
}

func (repository *ItemsRepository) Update(ctx context.Context, itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	transaction, err := repository.db.Beginx()
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				slog.ErrorContext(ctx, "rollback failed: %v", "error", rbErr)
			}
		}
	}()

	now := time.Now().UTC()

	if err != nil {
		repository.logger.ErrorContext(ctx, "transaction error", "error", err)
		return false, err
	}

	var updatedItem Item
	getQuery := "SELECT * FROM public.items WHERE id = ?"
	getQuery = repository.db.Rebind(getQuery)
	repository.logger.InfoContext(ctx, "fetching update item:", "query", getQuery, "id", itemID)
	err = transaction.Get(&updatedItem, getQuery, itemID)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on fetching item", "error", err, "id", itemID)
		return false, err
	}

	updatedItem.Name = update.Name
	updatedItem.Price = update.Price
	updatedItem.Description = update.Description
	updatedItem.IsActive = update.IsActive
	updatedItem.UpdatedAt = sql.NullTime{now, true}

	query := "UPDATE public.items SET name = ?, price = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "fetching update item:", "query", getQuery, "id", itemID)
	result, err := transaction.Exec(query, updatedItem.Name, updatedItem.Price, updatedItem.Description, updatedItem.IsActive,
		updatedItem.UpdatedAt, updatedItem.Id)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "id",
			updatedItem.Id, "name", updatedItem.Name, "price", updatedItem.Price, "description", updatedItem.Description, "isActive",
			updatedItem.IsActive, "updatedAt", now)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		return false, err
	}

	return rowsAffected > 0, err
}

func (repository *ItemsRepository) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
	transaction, err := repository.db.BeginTxx(ctx, &sql.TxOptions{Isolation: 4})
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				slog.ErrorContext(ctx, "rollback failed: %v", "error", rbErr)
			}
		}
	}()

	if err != nil {
		repository.logger.ErrorContext(ctx, "transaction error", "error", err)
		return false, err
	}

	query := "DELETE FROM public.items WHERE id = ?"
	query = repository.db.Rebind(query)
	repository.logger.InfoContext(ctx, "fetching delete item:", "query", query, "id", itemID)
	result, err := transaction.Exec(query, itemID)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "error", err, "id", itemID)
		return false, err
	}

	err = transaction.Commit()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on commit item", "error", err)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		return false, err
	}

	return rowsAffected > 0, err
}
