//go:build integration
// +build integration

package repositories_test

import (
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

func TestPriceHistoriesRepositoryGetByItemID_ShouldReturnErrorOnNilItemID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.Nil

	// Act
	priceHistories, err := repo.GetByItemID(ctx, itemID)

	// Assert
	require.EqualError(t, err, "itemID should not be nil")
	assert.Nil(t, priceHistories)
}

func TestPriceHistoriesRepositoryGetByItemID_ShouldReturnRecordsOrderedByRecordedAtDescending(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.New()
	firstHistoryID := uuid.New()
	secondHistoryID := uuid.New()
	olderDate := "2026-01-10"
	newerDate := "2026-01-15"
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Coffee", "FoodDrinks"}
	historyInsertQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`
	historyInsertArgs := []any{
		firstHistoryID, itemID, olderDate, decimal.RequireFromString("12.50"),
		secondHistoryID, itemID, newerDate, decimal.RequireFromString("15.00"),
	}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, historyInsertErr := testDB.Exec(historyInsertQuery, historyInsertArgs...)

	// Act
	priceHistories, getErr := repo.GetByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, historyInsertErr)
	require.NoError(t, getErr)
	require.Len(t, priceHistories, 2)
	assert.Equal(t, newerDate, priceHistories[0].RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("15.00").Equal(priceHistories[0].Value))
	assert.Equal(t, olderDate, priceHistories[1].RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, decimal.RequireFromString("12.50").Equal(priceHistories[1].Value))
}

func TestPriceHistoriesRepositoryGetByItemID_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewPriceHistoriesRepository(closedDB, testLogger)
	itemID := uuid.New()

	// Act
	priceHistories, err := repo.GetByItemID(ctx, itemID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, priceHistories)
}

func TestPriceHistoriesRepositoryUpsertToday_ShouldInsertNewRecordWhenTodayDoesNotExist(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.New()
	olderHistoryID := uuid.New()
	todayUTC := time.Now().UTC().Format("2006-01-02")
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Milk", "FoodDrinks"}
	historyInsertQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4)`
	historyInsertArgs := []any{olderHistoryID, itemID, "2026-01-10", decimal.RequireFromString("8.00")}
	upsert := &domains.PriceHistoryUpsert{
		Value: decimal.RequireFromString("9.50"),
	}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, historyInsertErr := testDB.Exec(historyInsertQuery, historyInsertArgs...)

	// Act
	priceHistory, upsertErr := repo.UpsertToday(ctx, itemID, upsert)
	priceHistories, getErr := repo.GetByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, historyInsertErr)
	require.NoError(t, upsertErr)
	require.NoError(t, getErr)
	require.NotNil(t, priceHistory)
	require.Len(t, priceHistories, 2)
	assert.NotEqual(t, uuid.Nil, priceHistory.Id)
	assert.NotEqual(t, olderHistoryID, priceHistory.Id)
	assert.Equal(t, itemID, priceHistory.ItemId)
	assert.Equal(t, todayUTC, priceHistory.RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.Value.Equal(priceHistory.Value))
	assert.Equal(t, todayUTC, priceHistories[0].RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.Value.Equal(priceHistories[0].Value))
}

func TestPriceHistoriesRepositoryUpsertToday_ShouldUpdateExistingRecordForToday(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.New()
	existingHistoryID := uuid.New()
	todayUTC := time.Now().UTC().Format("2006-01-02")
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Bread", "FoodDrinks"}
	historyInsertQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4)`
	historyInsertArgs := []any{existingHistoryID, itemID, todayUTC, decimal.RequireFromString("11.00")}
	upsert := &domains.PriceHistoryUpsert{
		Value: decimal.RequireFromString("13.25"),
	}
	countQuery := `SELECT COUNT(*) FROM price_history WHERE item_id = $1`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, historyInsertErr := testDB.Exec(historyInsertQuery, historyInsertArgs...)

	// Act
	priceHistory, upsertErr := repo.UpsertToday(ctx, itemID, upsert)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)
	priceHistories, getErr := repo.GetByItemID(ctx, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, historyInsertErr)
	require.NoError(t, upsertErr)
	require.NoError(t, countErr)
	require.NoError(t, getErr)
	require.NotNil(t, priceHistory)
	require.Len(t, priceHistories, 1)
	assert.Equal(t, 1, actualCount)
	assert.Equal(t, existingHistoryID, priceHistory.Id)
	assert.Equal(t, itemID, priceHistory.ItemId)
	assert.Equal(t, todayUTC, priceHistory.RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.Value.Equal(priceHistory.Value))
	assert.Equal(t, todayUTC, priceHistories[0].RecordedAt.UTC().Format("2006-01-02"))
	assert.True(t, upsert.Value.Equal(priceHistories[0].Value))
}

func TestPriceHistoriesRepositoryUpsertToday_ShouldReturnErrorOnNilItemID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.Nil
	upsert := &domains.PriceHistoryUpsert{
		Value: decimal.RequireFromString("10.00"),
	}

	// Act
	priceHistory, err := repo.UpsertToday(ctx, itemID, upsert)

	// Assert
	require.EqualError(t, err, "itemID should not be nil")
	assert.Nil(t, priceHistory)
}

func TestPriceHistoriesRepositoryUpsertToday_ShouldReturnErrorOnNilUpsert(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewPriceHistoriesRepository(testDB, testLogger)
	itemID := uuid.New()

	// Act
	priceHistory, err := repo.UpsertToday(ctx, itemID, nil)

	// Assert
	require.EqualError(t, err, "upsert should not be nil")
	assert.Nil(t, priceHistory)
}

func TestPriceHistoriesRepositoryUpsertToday_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewPriceHistoriesRepository(closedDB, testLogger)
	itemID := uuid.New()
	upsert := &domains.PriceHistoryUpsert{
		Value: decimal.RequireFromString("7.50"),
	}

	// Act
	priceHistory, err := repo.UpsertToday(ctx, itemID, upsert)

	// Assert
	require.Error(t, err)
	assert.Nil(t, priceHistory)
}
