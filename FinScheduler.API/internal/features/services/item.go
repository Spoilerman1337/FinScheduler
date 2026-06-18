package services

import (
	"context"
	"database/sql"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/persistence"
	"finscheduler/internal/traces"
	"finscheduler/pkg/dh"
	"finscheduler/pkg/sh"
	"fmt"
	"log/slog"
	"math"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
)

type ItemsService struct {
	uow    *persistence.UnitOfWork
	logger *slog.Logger
}

const itemsServiceName = "items"
const priceForecastMonthsAhead = 12
const priceForecastProjectedValueWindowPercent = 25

var averageDaysInMonth = decimal.RequireFromString("30.4375")
var priceForecastProjectedValueWindow = decimal.NewFromInt(priceForecastProjectedValueWindowPercent)

func NewItemsService(uow *persistence.UnitOfWork, logger *slog.Logger) *ItemsService {
	return &ItemsService{
		uow:    uow,
		logger: logger,
	}
}

func (service *ItemsService) GetListingInfo(ctx context.Context, filter *domains.ItemFilter) ([]domains.ItemListingDto, int64, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "GetListingInfo")
	defer span.End()

	if filter == nil {
		service.logger.ErrorContext(ctx, "filter is nil")
		err := fmt.Errorf("filter is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "GetListingInfo", err)
		return nil, 0, err
	}

	var items []domains.ItemListingDto
	var count int64

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawItems, rawItemsCount, err := repositories.Items.GetListingInfo(ctx, filter)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get items failed", "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "GetListingInfo", err)
			return err
		}

		count = rawItemsCount

		if len(rawItems) == 0 {
			items = make([]domains.ItemListingDto, 0)
			return nil
		}

		items = make([]domains.ItemListingDto, 0, len(rawItems))
		for _, item := range rawItems {
			items = append(items, *domains.NewItemListingDto(item))
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return items, count, err
}

func (service *ItemsService) GetDetailedInfo(ctx context.Context, itemID uuid.UUID) (*domains.ItemDetailedDto, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "GetDetailedInfo")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
		return nil, err
	}

	var item *domains.ItemDetailedDto

	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawItem, err := repositories.Items.GetDetailedInfo(ctx, itemID)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get item by id failed", "itemID", itemID, "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
			return err
		}

		rawPriceHistories, err := repositories.PriceHistories.GetByItemID(ctx, itemID)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get price histories by item id failed", "itemID", itemID, "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
			return err
		}

		priceForecastPoints := make([]domains.PriceForecastPointDto, 0)
		rawPriceForecast, err := repositories.PriceForecasts.GetLatestByItemID(ctx, itemID)
		if err != nil {
			if err != sql.ErrNoRows {
				service.logger.ErrorContext(ctx, "Get price forecast by item id failed", "itemID", itemID, "error", err)
				traces.EnrichFailedServiceSpan(span, err)
				metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
				return err
			}
		} else if rawPriceForecast != nil {
			priceForecastPoints = domains.BuildPriceForecastPoints(*rawPriceForecast, priceForecastMonthsAhead)
		}

		rawTagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{itemID})
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tag to items failed", "itemID", itemID, "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
			return err
		}

		tagIDs := make([]uuid.UUID, 0, len(rawTagToItems))
		for _, tagToItem := range rawTagToItems {
			tagIDs = append(tagIDs, tagToItem.TagId)
		}

		rawTags, err := repositories.Tags.GetByIds(ctx, tagIDs)
		if err != nil {
			service.logger.ErrorContext(ctx, "Get tags by ids failed", "itemID", itemID, "error", err)
			traces.EnrichFailedServiceSpan(span, err)
			metrics.RecordServiceFailure(ctx, itemsServiceName, "GetDetailedInfo", err)
			return err
		}

		item = domains.NewItemDetailedDto(*rawItem, rawTags, rawPriceHistories, priceForecastPoints)
		return nil
	})
	if err != nil {
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return item, nil
}

func (service *ItemsService) Create(ctx context.Context, create *domains.ItemCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Create")
	defer span.End()

	if create == nil {
		service.logger.ErrorContext(ctx, "create is nil")
		err := fmt.Errorf("create is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return uuid.Nil, err
	}

	if err := create.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "create validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return uuid.Nil, err
	}

	var newId uuid.UUID
	createTagIds := parseUUIDs(create.TagIds)

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error

		newId, err = repositories.Items.Create(ctx, create)
		if err != nil {
			return err
		}
		if newId == uuid.Nil {
			return fmt.Errorf("failed to create item: repository returned nil uuid")
		}

		if len(createTagIds) == 0 {
			return nil
		}

		success, err := repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: newId, TagIds: createTagIds})
		if err != nil {
			if details, ok := dh.GetPostgresErrorDetails(err); ok && details.Code == dh.PostgresForeignKeyViolationCode {
				return domains.ErrInvalidReference
			}
			return err
		}
		if !success {
			return fmt.Errorf("failed to create item: tag to item insert affected no rows")
		}

		return nil
	})

	if err != nil || newId == uuid.Nil {
		if err == nil {
			err = fmt.Errorf("failed to create item: repository returned nil uuid")
		}
		service.logger.ErrorContext(ctx, "error creating an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Create", err)
		return newId, err
	}

	traces.EnrichSuccessServiceSpan(span)

	return newId, err
}

func (service *ItemsService) Update(ctx context.Context, itemID uuid.UUID, update *domains.ItemUpdate) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Update")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}
	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}

	if err := update.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "update validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return false, err
	}

	var success bool
	updateTagIds := parseUUIDs(update.TagIds)

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		currentItem, err := repositories.Items.GetDetailedInfo(ctx, itemID)
		if err != nil {
			if err == sql.ErrNoRows {
				success = false
				return nil
			}

			return err
		}

		success, err = repositories.Items.Update(ctx, itemID, update)
		if err != nil {
			return err
		}
		if !success {
			return nil
		}

		if !currentItem.Price.Equal(update.Price) {
			_, err = repositories.PriceHistories.UpsertToday(ctx, itemID, &domains.PriceHistoryUpsert{Value: update.Price})
			if err != nil {
				return err
			}

			priceHistories, err := repositories.PriceHistories.GetByItemID(ctx, itemID)
			if err != nil {
				return err
			}

			priceForecastUpsert := buildPriceForecastUpsert(priceHistories)
			if priceForecastUpsert != nil {
				_, err = repositories.PriceForecasts.UpsertByItemID(ctx, itemID, priceForecastUpsert)
				if err != nil {
					return err
				}
			}
		}

		tagToItems, err := repositories.TagToItems.GetByItemIds(ctx, []uuid.UUID{itemID})
		if err != nil {
			return err
		}

		var currentTagIds []uuid.UUID
		if tagToItems != nil && len(tagToItems) > 0 {
			for _, tagToItem := range tagToItems {
				currentTagIds = append(currentTagIds, tagToItem.TagId)
			}
		}

		toDelete, toInsert := dh.Reconcile(updateTagIds, currentTagIds)

		if len(toDelete) > 0 {
			tagSuccess, err := repositories.TagToItems.BulkDelete(ctx, &domains.TagToItemDelete{ItemId: itemID, TagIds: toDelete})
			if err != nil {
				return err
			}
			if !tagSuccess {
				return fmt.Errorf("failed to update item: tag to item delete affected no rows")
			}
		}

		if len(toInsert) > 0 {
			tagSuccess, err := repositories.TagToItems.BulkInsert(ctx, &domains.TagToItemCreate{ItemId: itemID, TagIds: toInsert})
			if err != nil {
				if details, ok := dh.GetPostgresErrorDetails(err); ok && details.Code == dh.PostgresForeignKeyViolationCode {
					return domains.ErrInvalidReference
				}
				return err
			}
			if !tagSuccess {
				return fmt.Errorf("failed to update item: tag to item insert affected no rows")
			}
		}

		return nil
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "error updating an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Update", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
}

func (service *ItemsService) Delete(ctx context.Context, itemID uuid.UUID) (bool, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "Delete")
	defer span.End()

	if itemID == uuid.Nil {
		service.logger.ErrorContext(ctx, "itemID is nil")
		err := fmt.Errorf("itemID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Delete", err)
		return false, err
	}

	var success bool

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var err error
		success, err = repositories.Items.Delete(ctx, itemID)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "error deleting an item", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "Delete", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
}

func (service *ItemsService) UpdateCashbackByTag(ctx context.Context, update *domains.ItemCashbackByTagUpdate) (int64, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "UpdateCashbackByTag")
	defer span.End()

	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByTag", err)
		return 0, err
	}

	if err := update.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "update validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByTag", err)
		return 0, err
	}

	tagID, err := uuid.Parse(update.TagId)
	if err != nil {
		service.logger.ErrorContext(ctx, "failed to parse tag id", "tagId", update.TagId, "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByTag", err)
		return 0, err
	}

	var affected int64

	err = service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var repositoryErr error
		affected, repositoryErr = repositories.Items.UpdateCashbackByTag(ctx, tagID, update.Cashback)
		return repositoryErr
	})
	if err != nil {
		service.logger.ErrorContext(ctx, "error updating cashback by tag", "tagId", update.TagId, "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByTag", err)
		return 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return affected, nil
}

func (service *ItemsService) UpdateCashbackByIds(ctx context.Context, update *domains.ItemCashbackByIdsUpdate) (int64, error) {
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(ctx, "items-service")
	traces.RecordServiceSpan(span, "UpdateCashbackByIds")
	defer span.End()

	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByIds", err)
		return 0, err
	}

	if err := update.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "update validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByIds", err)
		return 0, err
	}

	var affected int64
	itemIDs := parseUUIDs(update.ItemIds)

	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var repositoryErr error
		affected, repositoryErr = repositories.Items.UpdateCashbackByIds(ctx, itemIDs, update.Cashback)
		return repositoryErr
	})
	if err != nil {
		service.logger.ErrorContext(ctx, "error updating cashback by ids", "itemIds", update.ItemIds, "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, itemsServiceName, "UpdateCashbackByIds", err)
		return 0, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return affected, nil
}

func parseUUIDs(ids []string) []uuid.UUID {
	if ids == nil {
		return nil
	}

	result := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		result[i] = uuid.MustParse(id)
	}

	return result
}

func buildPriceForecastUpsert(priceHistories []domains.PriceHistory) *domains.PriceForecastUpsert {
	if len(priceHistories) < 2 {
		return nil
	}

	winsorizedPriceHistories := winsorizePriceHistoriesForForecast(priceHistories)

	latestPriceHistory := priceHistories[0]
	winsorizedLatestPriceHistory := winsorizedPriceHistories[0]
	oldestPriceHistory := winsorizedPriceHistories[len(winsorizedPriceHistories)-1]
	daysBetween := latestPriceHistory.RecordedAt.Sub(oldestPriceHistory.RecordedAt).Hours() / 24
	if daysBetween <= 0 {
		return nil
	}

	monthsBetween := decimal.NewFromFloat(daysBetween).Div(averageDaysInMonth)
	if monthsBetween.IsZero() {
		return nil
	}

	averageMonthlyDrift := decimal.Zero
	if !oldestPriceHistory.Value.IsZero() {
		latestValue, _ := winsorizedLatestPriceHistory.Value.Float64()
		oldestValue, _ := oldestPriceHistory.Value.Float64()
		monthsBetweenFloat, _ := monthsBetween.Float64()
		monthlyGrowthFactor := math.Pow(latestValue/oldestValue, 1/monthsBetweenFloat)
		averageMonthlyDrift = decimal.NewFromFloat((monthlyGrowthFactor - 1) * 100).Round(6)
	} else if !latestPriceHistory.Value.IsZero() {
		return nil
	}

	return &domains.PriceForecastUpsert{
		LastKnownPrice:      latestPriceHistory.Value,
		AverageMonthlyDrift: averageMonthlyDrift,
	}
}

func winsorizePriceHistoriesForForecast(priceHistories []domains.PriceHistory) []domains.PriceHistory {
	result := make([]domains.PriceHistory, len(priceHistories))
	copy(result, priceHistories)

	for index := len(result) - 1; index >= 0; index-- {
		if len(result[index:]) < 3 {
			continue
		}

		projectedPriceValue := buildProjectedPriceValue(result[index:])
		if projectedPriceValue == nil {
			continue
		}

		result[index].Value = sh.WinsorizeValue(
			result[index].Value,
			*projectedPriceValue,
			priceForecastProjectedValueWindow,
		)
	}

	return result
}

func buildProjectedPriceValue(priceHistories []domains.PriceHistory) *decimal.Decimal {
	if len(priceHistories) < 3 {
		return nil
	}

	latestPriceHistory := priceHistories[0]
	previousLatestPriceHistory := priceHistories[1]
	oldestPriceHistory := priceHistories[len(priceHistories)-1]
	previousDaysBetween := previousLatestPriceHistory.RecordedAt.Sub(oldestPriceHistory.RecordedAt).Hours() / 24
	if previousDaysBetween <= 0 {
		return nil
	}

	previousMonthsBetween := decimal.NewFromFloat(previousDaysBetween).Div(averageDaysInMonth)
	if previousMonthsBetween.IsZero() || previousLatestPriceHistory.Value.IsZero() || oldestPriceHistory.Value.IsZero() {
		return nil
	}

	previousLatestValue, _ := previousLatestPriceHistory.Value.Float64()
	oldestValue, _ := oldestPriceHistory.Value.Float64()
	previousMonthsBetweenFloat, _ := previousMonthsBetween.Float64()
	previousMonthlyGrowthFactor := math.Pow(previousLatestValue/oldestValue, 1/previousMonthsBetweenFloat)

	latestIntervalDays := latestPriceHistory.RecordedAt.Sub(previousLatestPriceHistory.RecordedAt).Hours() / 24
	if latestIntervalDays <= 0 {
		return nil
	}

	latestIntervalMonths := decimal.NewFromFloat(latestIntervalDays).Div(averageDaysInMonth)
	latestIntervalMonthsFloat, _ := latestIntervalMonths.Float64()
	projectedLatestValue := previousLatestValue * math.Pow(previousMonthlyGrowthFactor, latestIntervalMonthsFloat)
	projectedPriceValue := decimal.NewFromFloat(projectedLatestValue)
	return &projectedPriceValue
}
