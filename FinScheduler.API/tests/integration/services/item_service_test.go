//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGetListingInfo_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	expectedName := "Item"
	page := int32(0)
	pageSize := int32(20)

	create := &domains.ItemCreate{
		Name:     expectedName,
		Category: "FoodDrinks",
	}

	// Act
	id, createErr := service.Create(ctx, create)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, getErr := service.GetListingInfo(ctx, &filter)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.Len(t, items, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, expectedName, items[0].Name)
}

func Test_ItemsService_CreateAndGetDetailedInfo_ShouldReturnAssignedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	tagsService := services.NewTagsService(uow, testLogger)
	tagName := "Groceries"
	itemName := "Milk"
	tagCreate := &domains.TagCreate{Name: tagName}

	tagID, tagCreateErr := tagsService.Create(ctx, tagCreate)

	create := &domains.ItemCreate{
		Name:     itemName,
		Category: "FoodDrinks",
		TagIds:   []string{tagID.String()},
	}

	// Act
	itemID, itemCreateErr := itemsService.Create(ctx, create)
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)

	// Assert
	require.NoError(t, tagCreateErr)
	require.NoError(t, itemCreateErr)
	require.NoError(t, getErr)
	require.NotNil(t, item)
	require.Len(t, item.Tags, 1)
	require.NotNil(t, item.PriceHistory)
	require.NotNil(t, item.PriceForecast)
	assert.Equal(t, tagName, item.Tags[0].Label)
	assert.Equal(t, tagID.String(), item.Tags[0].Value)
	assert.Empty(t, item.PriceHistory)
	assert.Empty(t, item.PriceForecast)
}

func Test_ItemsService_GetDetailedInfo_ShouldReturnPriceHistoryOrderedByDateDescending(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	itemName := "Milk"
	olderDate := "2026-01-10"
	newerDate := "2026-01-15"
	insertHistoryQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`
	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.RequireFromString("10.00"),
		Category: "FoodDrinks",
	}

	itemID, itemCreateErr := itemsService.Create(ctx, create)
	_, insertHistoryErr := testDB.Exec(
		insertHistoryQuery,
		uuid.New(), itemID, olderDate, decimal.RequireFromString("9.50"),
		uuid.New(), itemID, newerDate, decimal.RequireFromString("11.25"),
	)

	// Act
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)

	// Assert
	require.NoError(t, itemCreateErr)
	require.NoError(t, insertHistoryErr)
	require.NoError(t, getErr)
	require.NotNil(t, item)
	require.Len(t, item.PriceHistory, 2)
	require.Empty(t, item.PriceForecast)
	assert.Equal(t, newerDate, item.PriceHistory[0].Point.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("11.25").Equal(item.PriceHistory[0].Value))
	require.NotNil(t, item.PriceHistory[0].AbsoluteChange)
	require.NotNil(t, item.PriceHistory[0].PercentChange)
	assert.True(t, decimal.RequireFromString("1.75").Equal(*item.PriceHistory[0].AbsoluteChange))
	assert.True(t, decimal.RequireFromString("18.42105263157895").Equal(*item.PriceHistory[0].PercentChange))
	assert.Equal(t, olderDate, item.PriceHistory[1].Point.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("9.50").Equal(item.PriceHistory[1].Value))
	assert.Nil(t, item.PriceHistory[1].AbsoluteChange)
	assert.Nil(t, item.PriceHistory[1].PercentChange)
}

func Test_ItemsService_UpdateAndGetListingInfo_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	originalName := "Ice"
	updatedName := "Water"
	updatedPrice := 15.50
	page := int32(0)
	pageSize := int32(20)

	create := &domains.ItemCreate{
		Name:     originalName,
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &domains.ItemUpdate{
		Name:     updatedName,
		Price:    decimal.NewFromFloat(updatedPrice),
		Category: "FoodDrinks",
	}

	// Act
	id, createErr := service.Create(ctx, create)
	ok, updateErr := service.Update(ctx, id, update)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, getErr := service.GetListingInfo(ctx, &filter)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, updatedName, items[0].Name)
	assert.Equal(t, updatedPrice, items[0].Price)
}

func Test_ItemsService_Update_ShouldReconcileTagLinks(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	tagsService := services.NewTagsService(uow, testLogger)
	firstTagName := "Old Tag"
	secondTagName := "New Tag"
	itemName := "Tagged item"
	updatedItemName := "Tagged item updated"
	firstTagCreate := &domains.TagCreate{Name: firstTagName}
	secondTagCreate := &domains.TagCreate{Name: secondTagName}

	firstTagID, firstTagCreateErr := tagsService.Create(ctx, firstTagCreate)
	secondTagID, secondTagCreateErr := tagsService.Create(ctx, secondTagCreate)

	itemCreate := &domains.ItemCreate{
		Name:     itemName,
		Category: "FoodDrinks",
		TagIds:   []string{firstTagID.String()},
	}
	itemID, itemCreateErr := itemsService.Create(ctx, itemCreate)

	update := &domains.ItemUpdate{
		Name:     updatedItemName,
		Category: "FoodDrinks",
		TagIds:   []string{secondTagID.String()},
	}
	query := "SELECT tag_id FROM tag_to_item WHERE item_id = $1 ORDER BY tag_id"

	// Act
	ok, updateErr := itemsService.Update(ctx, itemID, update)
	var actualTagIDs []uuid.UUID
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)
	selectErr := testDB.Select(&actualTagIDs, query, itemID)

	// Assert
	require.NoError(t, firstTagCreateErr)
	require.NoError(t, secondTagCreateErr)
	require.NoError(t, itemCreateErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.NoError(t, selectErr)
	require.True(t, ok)
	require.NotNil(t, item)
	require.Len(t, item.Tags, 1)
	assert.Equal(t, updatedItemName, item.Name)
	assert.Equal(t, secondTagName, item.Tags[0].Label)
	assert.Equal(t, secondTagID.String(), item.Tags[0].Value)
	assert.Equal(t, []uuid.UUID{secondTagID}, actualTagIDs)
}

func Test_ItemsService_Update_ShouldUpsertPriceHistoryWhenPriceChanged(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	todayUTC := time.Now().UTC().Format("2006-01-02")
	countQuery := "SELECT COUNT(*) FROM price_history WHERE item_id = $1"
	forecastCountQuery := "SELECT COUNT(*) FROM price_forecast WHERE item_id = $1"
	create := &domains.ItemCreate{
		Name:     "Coffee",
		Price:    decimal.RequireFromString("10.00"),
		Category: "FoodDrinks",
	}
	update := &domains.ItemUpdate{
		Name:     "Coffee",
		Price:    decimal.RequireFromString("12.50"),
		Category: "FoodDrinks",
	}

	itemID, itemCreateErr := itemsService.Create(ctx, create)

	// Act
	ok, updateErr := itemsService.Update(ctx, itemID, update)
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)
	var actualForecastCount int
	forecastCountErr := testDB.Get(&actualForecastCount, forecastCountQuery, itemID)

	// Assert
	require.NoError(t, itemCreateErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.NoError(t, countErr)
	require.NoError(t, forecastCountErr)
	require.True(t, ok)
	require.NotNil(t, item)
	require.Len(t, item.PriceHistory, 1)
	require.Empty(t, item.PriceForecast)
	assert.Equal(t, 1, actualCount)
	assert.Equal(t, 0, actualForecastCount)
	assert.Equal(t, todayUTC, item.PriceHistory[0].Point.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("12.50").Equal(item.PriceHistory[0].Value))
	assert.Nil(t, item.PriceHistory[0].AbsoluteChange)
	assert.Nil(t, item.PriceHistory[0].PercentChange)
}

func Test_ItemsService_Update_ShouldNotUpsertPriceHistoryWhenPriceIsUnchanged(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	countQuery := "SELECT COUNT(*) FROM price_history WHERE item_id = $1"
	create := &domains.ItemCreate{
		Name:     "Tea",
		Price:    decimal.RequireFromString("7.50"),
		Category: "FoodDrinks",
	}
	update := &domains.ItemUpdate{
		Name:     "Tea updated",
		Price:    decimal.RequireFromString("7.50"),
		Category: "FoodDrinks",
	}

	itemID, itemCreateErr := itemsService.Create(ctx, create)

	// Act
	ok, updateErr := itemsService.Update(ctx, itemID, update)
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)

	// Assert
	require.NoError(t, itemCreateErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.NoError(t, countErr)
	require.True(t, ok)
	require.NotNil(t, item)
	assert.Equal(t, 0, actualCount)
	assert.Empty(t, item.PriceHistory)
	assert.Empty(t, item.PriceForecast)
}

func Test_ItemsService_Update_ShouldBuildPriceForecastWhenHistoryHasAtLeastTwoPoints(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	todayUTC := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
	middleDate := todayUTC.AddDate(0, 0, -15)
	oldestDate := todayUTC.AddDate(0, -3, 0)
	oldestValue := decimal.RequireFromString("10.00")
	middleValue := decimal.RequireFromString("11.00")
	newValue := decimal.RequireFromString("12.50")
	insertHistoryQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`
	forecastCountQuery := "SELECT COUNT(*) FROM price_forecast WHERE item_id = $1"
	create := &domains.ItemCreate{
		Name:     "Forecast coffee",
		Price:    oldestValue,
		Category: "FoodDrinks",
	}
	update := &domains.ItemUpdate{
		Name:     "Forecast coffee",
		Price:    newValue,
		Category: "FoodDrinks",
	}

	itemID, itemCreateErr := itemsService.Create(ctx, create)
	_, insertHistoryErr := testDB.Exec(
		insertHistoryQuery,
		uuid.New(), itemID, oldestDate.Format("2006-01-02"), oldestValue,
		uuid.New(), itemID, middleDate.Format("2006-01-02"), middleValue,
	)
	monthsBetween := decimal.NewFromFloat(todayUTC.Sub(oldestDate).Hours() / 24).Div(decimal.RequireFromString("30.4375"))
	monthsBetweenFloat, _ := monthsBetween.Float64()
	expectedDrift := decimal.NewFromFloat((math.Pow(12.5/10.0, 1/monthsBetweenFloat) - 1) * 100).Round(6)
	monthlyGrowthFactor := decimal.NewFromInt(1).Add(expectedDrift.Div(decimal.NewFromInt(100)))
	expectedFirstForecastValue := newValue.Mul(monthlyGrowthFactor).Round(2)
	expectedLastForecastValue := newValue
	for monthOffset := 0; monthOffset < 12; monthOffset++ {
		expectedLastForecastValue = expectedLastForecastValue.Mul(monthlyGrowthFactor)
	}
	expectedLastForecastValue = expectedLastForecastValue.Round(2)

	// Act
	ok, updateErr := itemsService.Update(ctx, itemID, update)
	item, getErr := itemsService.GetDetailedInfo(ctx, itemID)
	var actualForecastCount int
	forecastCountErr := testDB.Get(&actualForecastCount, forecastCountQuery, itemID)

	// Assert
	require.NoError(t, itemCreateErr)
	require.NoError(t, insertHistoryErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.NoError(t, forecastCountErr)
	require.True(t, ok)
	require.NotNil(t, item)
	require.Len(t, item.PriceForecast, 12)
	assert.Equal(t, 1, actualForecastCount)
	assert.Equal(t, todayUTC.AddDate(0, 1, 0), item.PriceForecast[0].Point)
	assert.True(t, expectedFirstForecastValue.Equal(item.PriceForecast[0].Value))
	assert.Equal(t, todayUTC.AddDate(0, 12, 0), item.PriceForecast[11].Point)
	assert.True(t, expectedLastForecastValue.Equal(item.PriceForecast[11].Value))
}

func Test_ItemsService_DeleteAndGetListingInfo_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	itemName := "Orange"
	page := int32(0)
	pageSize := int32(20)

	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, createErr := service.Create(ctx, create)
	ok, deleteErr := service.Delete(ctx, id)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, getErr := service.GetListingInfo(ctx, &filter)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, deleteErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	assert.Empty(t, items)
	assert.Equal(t, int64(0), count)
}

func Test_ItemsService_UpdateMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	missingID := uuid.New()

	update := &domains.ItemUpdate{
		Name:     "Missing",
		Category: "FoodDrinks",
	}

	// Act
	ok, err := service.Update(ctx, missingID, update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsService_DeleteMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	missingID := uuid.New()

	// Act
	ok, err := service.Delete(ctx, missingID)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsService_Create_ShouldRollbackItemWhenTagInsertFails(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewItemsService(uow, testLogger)
	itemName := "Rollback"
	expectedCount := 0
	invalidTagID := uuid.New()
	countQuery := "SELECT COUNT(*) FROM items WHERE name = $1"

	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
		TagIds:   []string{invalidTagID.String()},
	}

	// Act
	id, createErr := service.Create(ctx, create)
	var count int
	countErr := testDB.Get(&count, countQuery, create.Name)

	// Assert
	require.ErrorIs(t, createErr, domains.ErrInvalidReference)
	require.NoError(t, countErr)
	assert.NotEqual(t, uuid.Nil, id)
	assert.Equal(t, expectedCount, count)
}
