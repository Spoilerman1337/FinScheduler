package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"log/slog"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemsServiceGetListingInfo_ShouldReturnErrorOnNilFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var filter *domains.ItemFilter
	service := NewItemsService(uow, logger)

	// Act
	items, count, err := service.GetListingInfo(ctx, filter)

	// Assert
	require.EqualError(t, err, "filter is nil")
	assert.Nil(t, items)
	assert.Zero(t, count)
}

func TestItemsServiceCreate_ShouldReturnErrorOnNilCreate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var create *domains.ItemCreate
	service := NewItemsService(uow, logger)

	// Act
	newID, err := service.Create(ctx, create)

	// Assert
	require.EqualError(t, err, "create is nil")
	assert.Equal(t, uuid.Nil, newID)
}

func TestItemsServiceCreate_ShouldReturnValidationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	duplicateTagID := uuid.New().String()
	create := &domains.ItemCreate{
		Name:     "Coffee",
		Category: "FoodDrinks",
		TagIds:   []string{duplicateTagID, duplicateTagID},
	}
	service := NewItemsService(uow, logger)

	// Act
	newID, err := service.Create(ctx, create)

	// Assert
	require.EqualError(t, err, "tagId is duplicated: "+duplicateTagID)
	assert.Equal(t, uuid.Nil, newID)
}

func TestItemsServiceUpdate_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	nilID := uuid.Nil
	validID := uuid.New()
	update := &domains.ItemUpdate{}
	var nilUpdate *domains.ItemUpdate
	service := NewItemsService(uow, logger)

	// Act
	successOnNilID, errOnNilID := service.Update(ctx, nilID, update)
	successOnNilUpdate, errOnNilUpdate := service.Update(ctx, validID, nilUpdate)

	// Assert
	require.EqualError(t, errOnNilID, "itemID is nil")
	require.EqualError(t, errOnNilUpdate, "update is nil")
	assert.False(t, successOnNilID)
	assert.False(t, successOnNilUpdate)
}

func TestItemsServiceUpdate_ShouldReturnValidationError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	itemID := uuid.New()
	duplicateTagID := uuid.New().String()
	update := &domains.ItemUpdate{
		Name:     "Coffee",
		Category: "FoodDrinks",
		TagIds:   []string{duplicateTagID, duplicateTagID},
	}
	service := NewItemsService(uow, logger)

	// Act
	success, err := service.Update(ctx, itemID, update)

	// Assert
	require.EqualError(t, err, "tagId is duplicated: "+duplicateTagID)
	assert.False(t, success)
}

func TestItemsServiceDelete_ShouldReturnErrorOnNilItemID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	itemID := uuid.Nil
	service := NewItemsService(uow, logger)

	// Act
	success, err := service.Delete(ctx, itemID)

	// Assert
	require.EqualError(t, err, "itemID is nil")
	assert.False(t, success)
}

func TestItemsServiceUpdateCashbackByTag_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	service := NewItemsService(uow, logger)
	var nilUpdate *domains.ItemCashbackByTagUpdate
	invalidUpdate := &domains.ItemCashbackByTagUpdate{
		Cashback: 1,
		TagId:    "bad-uuid",
	}

	// Act
	affectedOnNilUpdate, errOnNilUpdate := service.UpdateCashbackByTag(ctx, nilUpdate)
	affectedOnInvalidUpdate, errOnInvalidUpdate := service.UpdateCashbackByTag(ctx, invalidUpdate)

	// Assert
	require.EqualError(t, errOnNilUpdate, "update is nil")
	require.EqualError(t, errOnInvalidUpdate, "tagId is invalid: bad-uuid")
	assert.Zero(t, affectedOnNilUpdate)
	assert.Zero(t, affectedOnInvalidUpdate)
}

func TestItemsServiceUpdateCashbackByIds_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	service := NewItemsService(uow, logger)
	var nilUpdate *domains.ItemCashbackByIdsUpdate
	invalidUpdate := &domains.ItemCashbackByIdsUpdate{
		Cashback: 1,
		ItemIds:  []string{"bad-uuid"},
	}

	// Act
	affectedOnNilUpdate, errOnNilUpdate := service.UpdateCashbackByIds(ctx, nilUpdate)
	affectedOnInvalidUpdate, errOnInvalidUpdate := service.UpdateCashbackByIds(ctx, invalidUpdate)

	// Assert
	require.EqualError(t, errOnNilUpdate, "update is nil")
	require.EqualError(t, errOnInvalidUpdate, "itemId is invalid: bad-uuid")
	assert.Zero(t, affectedOnNilUpdate)
	assert.Zero(t, affectedOnInvalidUpdate)
}

func TestBuildPriceForecastUpsert_ShouldReturnNilWhenThereAreLessThanTwoPoints(t *testing.T) {
	// Arrange
	priceHistories := []domains.PriceHistory{
		{
			RecordedAt: time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC),
			Value:      decimal.RequireFromString("100.00"),
		},
	}

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	assert.Nil(t, upsert)
}

func TestBuildPriceForecastUpsert_ShouldBuildMonthlyDriftFromTheWholeHistoryWindow(t *testing.T) {
	// Arrange
	latestDate := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	middleDate := latestDate.AddDate(0, 0, -15)
	oldestDate := latestDate.AddDate(0, -3, 0)
	priceHistories := []domains.PriceHistory{
		{
			RecordedAt: latestDate,
			Value:      decimal.RequireFromString("112.00"),
		},
		{
			RecordedAt: middleDate,
			Value:      decimal.RequireFromString("111.00"),
		},
		{
			RecordedAt: oldestDate,
			Value:      decimal.RequireFromString("100.00"),
		},
	}
	monthsBetween := decimal.NewFromFloat(latestDate.Sub(oldestDate).Hours() / 24).Div(averageDaysInMonth)
	monthsBetweenFloat, _ := monthsBetween.Float64()
	expectedDrift := decimal.NewFromFloat((math.Pow(112.0/100.0, 1/monthsBetweenFloat) - 1) * 100).Round(6)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.True(t, decimal.RequireFromString("112.00").Equal(upsert.LastKnownPrice))
	assert.True(t, expectedDrift.Equal(upsert.AverageMonthlyDrift))
}

func TestBuildPriceForecastUpsert_ShouldNotOverweightShortRecentSegments(t *testing.T) {
	// Arrange
	latestDate := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	recentDate := latestDate.AddDate(0, 0, -1)
	oldestDate := latestDate.AddDate(0, -1, 0)
	priceHistories := []domains.PriceHistory{
		{
			RecordedAt: latestDate,
			Value:      decimal.RequireFromString("115.00"),
		},
		{
			RecordedAt: recentDate,
			Value:      decimal.RequireFromString("110.00"),
		},
		{
			RecordedAt: oldestDate,
			Value:      decimal.RequireFromString("100.00"),
		},
	}
	monthsBetween := decimal.NewFromFloat(latestDate.Sub(oldestDate).Hours() / 24).Div(averageDaysInMonth)
	monthsBetweenFloat, _ := monthsBetween.Float64()
	expectedDrift := decimal.NewFromFloat((math.Pow(115.0/100.0, 1/monthsBetweenFloat) - 1) * 100).Round(6)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.True(t, expectedDrift.Equal(upsert.AverageMonthlyDrift))
	assert.True(t, upsert.AverageMonthlyDrift.LessThan(decimal.RequireFromString("20.00")))
}

func TestBuildPriceForecastUpsert_ShouldUseWinsorizedHistoryWhenCalculatingDrift(t *testing.T) {
	// Arrange
	latestDate := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	priceHistories := []domains.PriceHistory{
		{RecordedAt: latestDate, Value: decimal.RequireFromString("1000.00")},
		{RecordedAt: latestDate.AddDate(0, -1, 0), Value: decimal.RequireFromString("120.00")},
		{RecordedAt: latestDate.AddDate(0, -2, 0), Value: decimal.RequireFromString("118.00")},
		{RecordedAt: latestDate.AddDate(0, -3, 0), Value: decimal.RequireFromString("115.00")},
		{RecordedAt: latestDate.AddDate(0, -4, 0), Value: decimal.RequireFromString("112.00")},
		{RecordedAt: latestDate.AddDate(0, -5, 0), Value: decimal.RequireFromString("110.00")},
		{RecordedAt: latestDate.AddDate(0, -6, 0), Value: decimal.RequireFromString("108.00")},
		{RecordedAt: latestDate.AddDate(0, -7, 0), Value: decimal.RequireFromString("105.00")},
		{RecordedAt: latestDate.AddDate(0, -8, 0), Value: decimal.RequireFromString("102.00")},
		{RecordedAt: latestDate.AddDate(0, -9, 0), Value: decimal.RequireFromString("100.00")},
	}
	expectedClampedLatestValue := 137.5
	monthsBetween := decimal.NewFromFloat(latestDate.Sub(latestDate.AddDate(0, -9, 0)).Hours() / 24).Div(averageDaysInMonth)
	monthsBetweenFloat, _ := monthsBetween.Float64()
	expectedDrift := decimal.NewFromFloat(
		(math.Pow(expectedClampedLatestValue/100.0, 1/monthsBetweenFloat) - 1) * 100,
	).Round(6)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.True(t, decimal.RequireFromString("1000.00").Equal(upsert.LastKnownPrice))
	assert.True(t, expectedDrift.Equal(upsert.AverageMonthlyDrift))
}

func TestBuildPriceForecastUpsert_ShouldClampInternalPeaksByIQR(t *testing.T) {
	// Arrange
	latestDate := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	priceHistories := []domains.PriceHistory{
		{RecordedAt: latestDate, Value: decimal.RequireFromString("135.00")},
		{RecordedAt: latestDate.AddDate(0, -1, 0), Value: decimal.RequireFromString("1000.00")},
		{RecordedAt: latestDate.AddDate(0, -2, 0), Value: decimal.RequireFromString("120.00")},
		{RecordedAt: latestDate.AddDate(0, -3, 0), Value: decimal.RequireFromString("118.00")},
		{RecordedAt: latestDate.AddDate(0, -4, 0), Value: decimal.RequireFromString("115.00")},
		{RecordedAt: latestDate.AddDate(0, -5, 0), Value: decimal.RequireFromString("112.00")},
		{RecordedAt: latestDate.AddDate(0, -6, 0), Value: decimal.RequireFromString("110.00")},
		{RecordedAt: latestDate.AddDate(0, -7, 0), Value: decimal.RequireFromString("108.00")},
		{RecordedAt: latestDate.AddDate(0, -8, 0), Value: decimal.RequireFromString("105.00")},
		{RecordedAt: latestDate.AddDate(0, -9, 0), Value: decimal.RequireFromString("100.00")},
	}
	monthsBetween := decimal.NewFromFloat(latestDate.Sub(latestDate.AddDate(0, -9, 0)).Hours() / 24).Div(averageDaysInMonth)
	monthsBetweenFloat, _ := monthsBetween.Float64()
	expectedDrift := decimal.NewFromFloat((math.Pow(135.0/100.0, 1/monthsBetweenFloat) - 1) * 100).Round(6)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.True(t, decimal.RequireFromString("135.00").Equal(upsert.LastKnownPrice))
	assert.True(t, expectedDrift.Equal(upsert.AverageMonthlyDrift))
}

func TestBuildPriceForecastUpsert_ShouldNotMutateOriginalPriceHistories(t *testing.T) {
	// Arrange
	latestDate := time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)
	priceHistories := []domains.PriceHistory{
		{
			RecordedAt: latestDate,
			Value:      decimal.RequireFromString("140.00"),
		},
		{
			RecordedAt: latestDate.AddDate(0, -1, 0),
			Value:      decimal.RequireFromString("1000.00"),
		},
		{
			RecordedAt: latestDate.AddDate(0, -2, 0),
			Value:      decimal.RequireFromString("105.00"),
		},
		{
			RecordedAt: latestDate.AddDate(0, -3, 0),
			Value:      decimal.RequireFromString("102.00"),
		},
		{
			RecordedAt: latestDate.AddDate(0, -4, 0),
			Value:      decimal.RequireFromString("100.00"),
		},
	}
	originalPriceHistories := append([]domains.PriceHistory(nil), priceHistories...)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.Equal(t, originalPriceHistories, priceHistories)
}
