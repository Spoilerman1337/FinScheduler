package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"log/slog"
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

func TestBuildPriceForecastUpsert_ShouldAverageMonthlyDriftsBetweenNeighborPoints(t *testing.T) {
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
	firstSegmentMonthsBetween := decimal.NewFromFloat(latestDate.Sub(middleDate).Hours() / 24).Div(averageDaysInMonth)
	secondSegmentMonthsBetween := decimal.NewFromFloat(middleDate.Sub(oldestDate).Hours() / 24).Div(averageDaysInMonth)
	firstSegmentDrift := decimal.RequireFromString("112.00").
		Sub(decimal.RequireFromString("111.00")).
		Div(decimal.RequireFromString("111.00")).
		Mul(decimal.NewFromInt(100)).
		Div(firstSegmentMonthsBetween)
	secondSegmentDrift := decimal.RequireFromString("111.00").
		Sub(decimal.RequireFromString("100.00")).
		Div(decimal.RequireFromString("100.00")).
		Mul(decimal.NewFromInt(100)).
		Div(secondSegmentMonthsBetween)
	expectedDrift := firstSegmentDrift.Add(secondSegmentDrift).Div(decimal.NewFromInt(2)).Round(2)

	// Act
	upsert := buildPriceForecastUpsert(priceHistories)

	// Assert
	require.NotNil(t, upsert)
	assert.True(t, decimal.RequireFromString("112.00").Equal(upsert.LastKnownPrice))
	assert.True(t, expectedDrift.Equal(upsert.AverageMonthlyDrift))
}
