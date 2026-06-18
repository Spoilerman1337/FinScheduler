//go:build integration
// +build integration

package repositories_test

import (
	"database/sql"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/tests/internal/testsupport"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceForecastsRepositoryGetLatestByItemID_ShouldReturnErrorOnNilItemID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)

	// Act
	priceForecast, err := repo.GetLatestByItemID(ctx, uuid.Nil)

	// Assert
	require.EqualError(t, err, "itemID should not be nil")
	assert.Nil(t, priceForecast)
}

func TestPriceForecastsRepositoryGetLatestByItemID_ShouldReturnStoredForecast(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	itemID := uuid.New()
	olderForecastID := uuid.New()
	newerForecastID := uuid.New()
	olderCalculatedAt := "2026-06-17"
	newerCalculatedAt := "2026-06-18"
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	forecastInsertQuery := `INSERT INTO price_forecast (id, item_id, calculated_at, last_known_price, average_monthly_drift) VALUES ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemID, "Coffee", "FoodDrinks")
	_, forecastInsertErr := testDB.Exec(
		forecastInsertQuery,
		olderForecastID,
		itemID,
		olderCalculatedAt,
		decimal.RequireFromString("100.50"),
		decimal.RequireFromString("-2.250123"),
		newerForecastID,
		itemID,
		newerCalculatedAt,
		decimal.RequireFromString("101.75"),
		decimal.RequireFromString("1.500456"),
	)

	// Act
	priceForecast, getErr := repo.GetLatestByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, forecastInsertErr)
	require.NoError(t, getErr)
	require.NotNil(t, priceForecast)
	assert.Equal(t, newerForecastID, priceForecast.Id)
	assert.Equal(t, itemID, priceForecast.ItemId)
	assert.Equal(t, newerCalculatedAt, priceForecast.CalculatedAt.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("101.75").Equal(priceForecast.LastKnownPrice))
	assert.True(t, decimal.RequireFromString("1.500456").Equal(priceForecast.AverageMonthlyDrift))
}

func TestPriceForecastsRepositoryGetLatestByItemID_ShouldReturnNoRowsWhenForecastMissing(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	itemID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemID, "Tea", "FoodDrinks")

	// Act
	priceForecast, err := repo.GetLatestByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, priceForecast)
}

func TestPriceForecastsRepositoryUpsertByItemID_ShouldInsertNewRecordWhenForecastDoesNotExist(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	itemID := uuid.New()
	todayUTC := time.Now().UTC().Format("2006-01-02")
	upsert := &domains.PriceForecastUpsert{
		LastKnownPrice:      decimal.RequireFromString("120.00"),
		AverageMonthlyDrift: decimal.RequireFromString("3.750123"),
	}
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemID, "Bread", "FoodDrinks")

	// Act
	priceForecast, upsertErr := repo.UpsertByItemID(ctx, itemID, upsert)
	loadedPriceForecast, getErr := repo.GetLatestByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, upsertErr)
	require.NoError(t, getErr)
	require.NotNil(t, priceForecast)
	require.NotNil(t, loadedPriceForecast)
	assert.NotEqual(t, uuid.Nil, priceForecast.Id)
	assert.Equal(t, itemID, priceForecast.ItemId)
	assert.Equal(t, todayUTC, priceForecast.CalculatedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.LastKnownPrice.Equal(priceForecast.LastKnownPrice))
	assert.True(t, upsert.AverageMonthlyDrift.Equal(priceForecast.AverageMonthlyDrift))
	assert.Equal(t, priceForecast.Id, loadedPriceForecast.Id)
}

func TestPriceForecastsRepositoryUpsertByItemID_ShouldInsertNewRecordWhenItemHasForecastForAnotherDate(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	itemID := uuid.New()
	existingForecastID := uuid.New()
	todayUTC := time.Now().UTC().Format("2006-01-02")
	initialCalculatedAt := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")
	upsert := &domains.PriceForecastUpsert{
		LastKnownPrice:      decimal.RequireFromString("150.00"),
		AverageMonthlyDrift: decimal.RequireFromString("-1.500321"),
	}
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	forecastInsertQuery := `INSERT INTO price_forecast (id, item_id, calculated_at, last_known_price, average_monthly_drift) VALUES ($1, $2, $3, $4, $5)`
	countQuery := `SELECT COUNT(*) FROM price_forecast WHERE item_id = $1`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemID, "Milk", "FoodDrinks")
	_, forecastInsertErr := testDB.Exec(
		forecastInsertQuery,
		existingForecastID,
		itemID,
		initialCalculatedAt,
		decimal.RequireFromString("100.00"),
		decimal.RequireFromString("2.000111"),
	)

	// Act
	priceForecast, upsertErr := repo.UpsertByItemID(ctx, itemID, upsert)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, forecastInsertErr)
	require.NoError(t, upsertErr)
	require.NoError(t, countErr)
	require.NotNil(t, priceForecast)
	assert.Equal(t, 2, actualCount)
	assert.NotEqual(t, existingForecastID, priceForecast.Id)
	assert.Equal(t, itemID, priceForecast.ItemId)
	assert.Equal(t, todayUTC, priceForecast.CalculatedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.LastKnownPrice.Equal(priceForecast.LastKnownPrice))
	assert.True(t, upsert.AverageMonthlyDrift.Equal(priceForecast.AverageMonthlyDrift))
}

func TestPriceForecastsRepositoryUpsertByItemID_ShouldUpdateExistingRecordForSameItemAndDate(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	itemID := uuid.New()
	existingForecastID := uuid.New()
	todayUTC := time.Now().UTC().Format("2006-01-02")
	upsert := &domains.PriceForecastUpsert{
		LastKnownPrice:      decimal.RequireFromString("150.00"),
		AverageMonthlyDrift: decimal.RequireFromString("-1.500321"),
	}
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	forecastInsertQuery := `INSERT INTO price_forecast (id, item_id, calculated_at, last_known_price, average_monthly_drift) VALUES ($1, $2, $3, $4, $5)`
	countQuery := `SELECT COUNT(*) FROM price_forecast WHERE item_id = $1`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemID, "Milk", "FoodDrinks")
	_, forecastInsertErr := testDB.Exec(
		forecastInsertQuery,
		existingForecastID,
		itemID,
		todayUTC,
		decimal.RequireFromString("100.00"),
		decimal.RequireFromString("2.000111"),
	)

	// Act
	priceForecast, upsertErr := repo.UpsertByItemID(ctx, itemID, upsert)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, forecastInsertErr)
	require.NoError(t, upsertErr)
	require.NoError(t, countErr)
	require.NotNil(t, priceForecast)
	assert.Equal(t, 1, actualCount)
	assert.Equal(t, existingForecastID, priceForecast.Id)
	assert.Equal(t, itemID, priceForecast.ItemId)
	assert.Equal(t, todayUTC, priceForecast.CalculatedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.LastKnownPrice.Equal(priceForecast.LastKnownPrice))
	assert.True(t, upsert.AverageMonthlyDrift.Equal(priceForecast.AverageMonthlyDrift))
}

func TestPriceForecastsRepositoryUpsertByItemID_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewPriceForecastsRepository(testDB, testLogger)
	upsert := &domains.PriceForecastUpsert{
		LastKnownPrice:      decimal.RequireFromString("10.00"),
		AverageMonthlyDrift: decimal.RequireFromString("1.00"),
	}

	// Act
	priceForecastOnNilItemID, errOnNilItemID := repo.UpsertByItemID(ctx, uuid.Nil, upsert)
	priceForecastOnNilUpsert, errOnNilUpsert := repo.UpsertByItemID(ctx, uuid.New(), nil)

	// Assert
	require.EqualError(t, errOnNilItemID, "itemID should not be nil")
	require.EqualError(t, errOnNilUpsert, "upsert should not be nil")
	assert.Nil(t, priceForecastOnNilItemID)
	assert.Nil(t, priceForecastOnNilUpsert)
}

func TestPriceForecastsRepositoryUpsertByItemID_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewPriceForecastsRepository(closedDB, testLogger)
	upsert := &domains.PriceForecastUpsert{
		LastKnownPrice:      decimal.RequireFromString("10.00"),
		AverageMonthlyDrift: decimal.RequireFromString("1.00"),
	}

	// Act
	priceForecast, err := repo.UpsertByItemID(ctx, uuid.New(), upsert)

	// Assert
	require.Error(t, err)
	assert.Nil(t, priceForecast)
}
