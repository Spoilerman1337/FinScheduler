package items

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

type ItemsRepo struct {
	db *sqlx.DB
}

func NewItemsRepo(db *sqlx.DB) *ItemsRepo {
	return &ItemsRepo{db: db}
}

func (repo *ItemsRepo) Get(filter *ItemFilter) ([]Item, int64, error) {
	//TODO: log error later
	var items []Item
	var count int64 = 0

	var query = " FROM public.items WHERE 1=1"
	args := make([]interface{}, 0)

	if filter.Ids != nil && len(filter.Ids) > 0 {
		inQuery, inArgs, err := sqlx.In(" AND Id IN (?)", filter.Ids)

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
		query += " AND isActive = ?"
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

func (repo *ItemsRepo) GetById(id uuid.UUID) (*Item, error) {
	//TODO: log error later
	var item Item

	if id == uuid.Nil {
		return nil, fmt.Errorf("id should not be nil")
	}

	query := "SELECT * FROM public.items WHERE id = ?"
	query = repo.db.Rebind(query)

	err := repo.db.Get(&item, query, id)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *ItemsRepo) Create(create *ItemCreate) (uuid.UUID, error) {
	newID, err := uuid.NewV7()

	if err != nil {
		return uuid.Nil, err
	}

	query := "INSERT INTO public.items (id, name, price, description, is_active) VALUES (?, ?, ?, ?, ?)"
	_, err = repo.db.Exec(query, newID, create.Name, create.Price, create.Description, create.IsActive)
	if err != nil {
		return uuid.Nil, err
	}

	return newID, err
}

func (repo *ItemsRepo) Update(itemID uuid.UUID, update *ItemUpdate) (bool, error) {
	transaction, err := repo.db.Beginx()
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}
		}
	}()

	if err != nil {
		return false, err
	}

	var updatedItem *Item
	err = repo.db.Get(&updatedItem, "SELECT * FROM public.items WHERE id = ?", itemID)
	if err != nil {
		return false, err
	}
	if updatedItem == nil {
		return false, nil
	}

	updatedItem.Name = update.Name
	updatedItem.Price = update.Price
	updatedItem.Description = update.Description
	updatedItem.IsActive = update.IsActive

	query := "UPDATE public.items SET name = ?, price = ?, description = ?, is_active = ?) WHERE id = ?"
	result, err := repo.db.Exec(query, updatedItem.Name, updatedItem.Price, updatedItem.Description, updatedItem.IsActive, updatedItem.Id)
	if err != nil {
		return false, err
	}

	err = transaction.Commit()
	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err
}

func (repo *ItemsRepo) Delete(itemID uuid.UUID) (bool, error) {
	transaction, err := repo.db.BeginTxx(nil, &sql.TxOptions{Isolation: 4})
	defer func() {
		if err != nil {
			rbErr := transaction.Rollback()
			if rbErr != nil {
				log.Printf("rollback failed: %v", rbErr)
			}
		}
	}()

	query := "DELETE FROM public.items WHERE id = ?"
	result, err := repo.db.Exec(query, itemID)

	err = transaction.Commit()
	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err
}
