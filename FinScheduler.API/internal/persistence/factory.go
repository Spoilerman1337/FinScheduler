package persistence

import (
	"finscheduler/internal/features/repositories"
	"log/slog"
)

type RepositoryFactory struct {
	db     repositories.DBTX
	logger *slog.Logger
}

func NewRepositoryFactory(db repositories.DBTX, logger *slog.Logger) *RepositoryFactory {
	return &RepositoryFactory{db: db, logger: logger}
}

func (factory *RepositoryFactory) Items() *repositories.ItemsRepository {
	return repositories.NewItemsRepository(factory.db, factory.logger)
}

func (factory *RepositoryFactory) PriceForecasts() *repositories.PriceForecastsRepository {
	return repositories.NewPriceForecastsRepository(factory.db, factory.logger)
}

func (factory *RepositoryFactory) PriceHistories() *repositories.PriceHistoriesRepository {
	return repositories.NewPriceHistoriesRepository(factory.db, factory.logger)
}

func (factory *RepositoryFactory) Tags() *repositories.TagsRepository {
	return repositories.NewTagsRepository(factory.db, factory.logger)
}

func (factory *RepositoryFactory) TagToItems() *repositories.TagToItemsRepository {
	return repositories.NewTagToItemsRepository(factory.db, factory.logger)
}
