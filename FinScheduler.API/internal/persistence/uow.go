package persistence

import (
	"context"
	"finscheduler/internal/features/repositories"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewUnitOfWork(db *sqlx.DB, logger *slog.Logger) *UnitOfWork {
	return &UnitOfWork{db: db, logger: logger}
}

type Repositories struct {
	Items          *repositories.ItemsRepository
	PriceHistories *repositories.PriceHistoriesRepository
	Tags           *repositories.TagsRepository
	TagToItems     *repositories.TagToItemsRepository
}

func (uow *UnitOfWork) WithoutTx(fn func(Repositories) error) error {
	return fn(uow.buildRepositories(uow.db))
}

func (uow *UnitOfWork) WithTx(ctx context.Context, fn func(Repositories) error) error {
	tx, err := uow.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if err := fn(uow.buildRepositories(tx)); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	committed = true
	return nil
}

func (uow *UnitOfWork) buildRepositories(db repositories.DBTX) Repositories {
	factory := NewRepositoryFactory(db, uow.logger)

	return Repositories{
		Items:          factory.Items(),
		PriceHistories: factory.PriceHistories(),
		Tags:           factory.Tags(),
		TagToItems:     factory.TagToItems(),
	}
}
