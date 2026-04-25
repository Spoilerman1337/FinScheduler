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

func Test_ItemsRepository_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	create := &domains.ItemCreate{
		Name:     "Apple",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	item, err := repo.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	require.NotNil(t, item)
	assert.Equal(t, "Apple", item.Name)
	assert.True(t, decimal.NewFromFloat(15.50).Equal(item.Price))
}

func Test_ItemsRepository_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	create := &domains.ItemCreate{
		Name:     "Old",
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &domains.ItemUpdate{
		Name:     "New",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)

	ok, err := repo.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := repo.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "New", item.Name)
	assert.True(t, decimal.NewFromFloat(15.50).Equal(item.Price))
}

func Test_ItemsRepository_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	create := &domains.ItemCreate{
		Name:     "Book",
		Price:    decimal.NewFromFloat(10.00),
		Category: "Entertainments",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)

	ok, err := repo.Delete(testContext, id)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := repo.GetById(testContext, id)

	// Assert
	assert.Nil(t, item)
	assert.Error(t, err)
}

func Test_ItemsRepository_Get_ShouldApplyFiltersAndCount(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	foodTag := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food", IsActive: true})
	travelTag := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Travel", IsActive: true})

	coffeeID := testFixtures.MustCreateItem(t, &domains.ItemCreate{
		Name:        "Coffee Beans",
		Price:       decimal.NewFromFloat(12.50),
		Description: "Fresh arabica beans",
		IsActive:    true,
		Cashback:    5,
		Category:    string(domains.FoodDrinks),
	})
	taxiID := testFixtures.MustCreateItem(t, &domains.ItemCreate{
		Name:        "Taxi Ride",
		Price:       decimal.NewFromFloat(35.00),
		Description: "Airport transfer",
		IsActive:    false,
		Cashback:    0,
		Category:    string(domains.Transport),
	})
	testFixtures.MustLinkItemTags(t, coffeeID, foodTag)
	testFixtures.MustLinkItemTags(t, taxiID, travelTag)

	coffeeCreated := time.Date(2026, time.April, 20, 10, 0, 0, 0, time.UTC)
	coffeeUpdated := time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC)
	_, err := testDB.Exec(
		`UPDATE items SET created_at = $2, updated_at = $3 WHERE id = $1`,
		coffeeID, coffeeCreated, coffeeUpdated,
	)
	require.NoError(t, err)

	page := int32(0)
	pageSize := int32(10)
	name := "Coffee"
	description := "arabica"
	isActive := true
	priceFrom := decimal.NewFromFloat(10)
	priceTo := decimal.NewFromFloat(15)
	cashbackFrom := int32(1)
	cashbackTo := int32(5)
	createdFrom := coffeeCreated.Add(-time.Hour)
	createdTo := coffeeCreated.Add(time.Hour)
	updatedFrom := coffeeUpdated.Add(-time.Hour)
	updatedTo := coffeeUpdated.Add(time.Hour)
	category := domains.FoodDrinks

	// Act
	items, count, err := repo.Get(testContext, &domains.ItemFilter{
		Name:         &name,
		Description:  &description,
		IsActive:     &isActive,
		PriceFrom:    &priceFrom,
		PriceTo:      &priceTo,
		CreatedFrom:  &createdFrom,
		CreatedTo:    &createdTo,
		UpdatedFrom:  &updatedFrom,
		UpdatedTo:    &updatedTo,
		CashbackFrom: &cashbackFrom,
		CashbackTo:   &cashbackTo,
		Categories:   []*domains.ItemCategory{&category},
		TagIds:       []*uuid.UUID{&foodTag},
		Page:         &page,
		PageSize:     &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, coffeeID, items[0].Id)
	assert.Equal(t, "Coffee Beans", items[0].Name)
}

func Test_ItemsRepository_Get_ShouldPaginateAndOrderByCreatedAtDesc(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	olderID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Older", Category: string(domains.FoodDrinks)})
	middleID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Middle", Category: string(domains.FoodDrinks)})
	newestID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Newest", Category: string(domains.FoodDrinks)})

	_, err := testDB.Exec(`UPDATE items SET created_at = $2 WHERE id = $1`,
		olderID, time.Date(2026, time.April, 20, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	_, err = testDB.Exec(`UPDATE items SET created_at = $2 WHERE id = $1`,
		middleID, time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	_, err = testDB.Exec(`UPDATE items SET created_at = $2 WHERE id = $1`,
		newestID, time.Date(2026, time.April, 22, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	page := int32(1)
	pageSize := int32(1)

	// Act
	items, count, err := repo.Get(testContext, &domains.ItemFilter{
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int64(3), count)
	assert.Equal(t, middleID, items[0].Id)
}

func Test_ItemsRepository_GetById_ShouldReturnErrorOnNilUUID(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	// Act
	item, err := repo.GetById(testContext, uuid.Nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, item)
	assert.Contains(t, err.Error(), "id should not be nil")
}

func Test_ItemsRepository_Update_ShouldReturnFalseForMissingRow(t *testing.T) {
	// Arrange
	repo := repositories.NewItemsRepository(testDB, testLogger)

	// Act
	ok, err := repo.Update(testContext, uuid.New(), &domains.ItemUpdate{
		Name:     "Missing",
		Price:    decimal.NewFromFloat(10),
		Category: string(domains.FoodDrinks),
	})

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsRepository_Delete_ShouldCascadeTagRelations(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewItemsRepository(testDB, testLogger)

	itemID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Delete Me", Category: string(domains.FoodDrinks)})
	tagID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Attached", IsActive: true})
	testFixtures.MustLinkItemTags(t, itemID, tagID)

	// Act
	ok, err := repo.Delete(testContext, itemID)

	// Assert
	require.NoError(t, err)
	assert.True(t, ok)

	var relationCount int
	require.NoError(t, testDB.Get(&relationCount, "SELECT COUNT(*) FROM tag_to_item WHERE item_id = $1", itemID))
	assert.Equal(t, 0, relationCount)
}
