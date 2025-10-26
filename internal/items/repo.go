package items

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type ItemsRepository struct {
	db *sqlx.DB
}

func NewItemsRepository(db *sqlx.DB) *ItemsRepository {
	return &ItemsRepository{db: db}
}

func (repo *ItemsRepository) Get(filter *ItemFilter) ([]Item, int64, error) {
	// TODO: log error later
	var items []Item
	var count int64 = 0

	var query = " FROM public.items WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND id IN (?)", filter.Ids)

		if err != nil {
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
	selectQuery = repo.db.Rebind(selectQuery)
	selectArgs := append(make([]interface{}, 0), args...)
	selectArgs = append(selectArgs, pageSize, offset)

	err := repo.db.Select(&items, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) %s", query)
	countQuery = repo.db.Rebind(countQuery)
	countArgs := append(make([]interface{}, 0), args...)

	err = repo.db.Get(&count, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	return items, count, err
}

func (repo *ItemsRepository) GetById(id uuid.UUID) (*Item, error) {
	//TODO: log error later
	var item Item

	if id == uuid.Nil {
		return nil, fmt.Errorf("id should not be nil")
	}

	query := "SELECT * FROM public.items WHERE id = ?"
	query = repo.db.Rebind(query)

	err := repo.db.Get(&item, query, id)
	if err != nil {
		fmt.Print("1112")
		return nil, err
	}

	return &item, nil
}

func (repo *ItemsRepository) Create(create *ItemCreate) (uuid.UUID, error) {
	newID, err := uuid.NewV7()
	now := time.Now().UTC()

	if err != nil {
		return uuid.Nil, err
	}

	query := "INSERT INTO public.items (id, name, price, description, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	query = repo.db.Rebind(query)
	_, err = repo.db.Exec(query, newID, create.Name, create.Price, create.Description, create.IsActive, now)
	if err != nil {
		return uuid.Nil, err
	}

	return newID, err
}

func (repo *ItemsRepository) Update(itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	transaction, err := repo.db.Beginx()
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}
		}
	}()

	now := time.Now().UTC()

	if err != nil {
		return false, err
	}

	var updatedItem Item
	getQuery := "SELECT * FROM public.items WHERE id = ?"
	getQuery = repo.db.Rebind(getQuery)
	err = transaction.Get(&updatedItem, getQuery, itemID)
	if err != nil {
		return false, err
	}

	updatedItem.Name = update.Name
	updatedItem.Price = update.Price
	updatedItem.Description = update.Description
	updatedItem.IsActive = update.IsActive
	updatedItem.UpdatedAt = sql.NullTime{now, true}

	query := "UPDATE public.items SET name = ?, price = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?"
	query = repo.db.Rebind(query)
	result, err := transaction.Exec(query, updatedItem.Name, updatedItem.Price, updatedItem.Description, updatedItem.IsActive,
		updatedItem.UpdatedAt, updatedItem.Id)
	if err != nil {
		return false, err
	}

	err = transaction.Commit()
	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err
}

func (repo *ItemsRepository) Delete(itemID uuid.UUID) (bool, error) {
	ctx := context.Background()
	transaction, err := repo.db.BeginTxx(ctx, &sql.TxOptions{Isolation: 4})
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}
		}
	}()

	query := "DELETE FROM public.items WHERE id = ?"
	query = repo.db.Rebind(query)
	result, err := transaction.Exec(query, itemID)

	err = transaction.Commit()
	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err
}
