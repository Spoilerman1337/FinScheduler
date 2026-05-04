//go:build integration
// +build integration

package repositories_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/tests/internal/testsupport"
	"testing"

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

	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	itemName := "Apple"
	itemPrice := decimal.NewFromFloat(15.50)
	itemCategory := "FoodDrinks"

	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    itemPrice,
		Category: itemCategory,
	}

	// Act
	id, createErr := repo.Create(ctx, create)
	item, getErr := repo.GetById(ctx, id)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.NotEqual(t, uuid.Nil, id)
	require.NotNil(t, item)
	assert.Equal(t, itemName, item.Name)
	assert.True(t, itemPrice.Equal(item.Price))
}

func Test_ItemsRepository_Get_ShouldFilterAndReturnCount(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	targetTagID := uuid.New()
	otherTagID := uuid.New()
	firstName := "Coffee Beans"
	secondName := "Coffee Syrup"
	thirdName := "Coffee Gift"
	fourthName := "Tea"
	foodCategory := "FoodDrinks"
	giftCategory := "Gifts"
	firstPrice := decimal.NewFromFloat(10.00)
	secondPrice := decimal.NewFromFloat(20.00)
	thirdPrice := decimal.NewFromFloat(30.00)
	fourthPrice := decimal.NewFromFloat(40.00)
	tagInsertQuery := "INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)"
	tagInsertArgs := []any{targetTagID, "Target", true, otherTagID, "Other", true}
	linkInsertQuery := "INSERT INTO tag_to_item (item_id, tag_id) VALUES ($1, $2), ($3, $4), ($5, $6), ($7, $8)"
	filterName := "Coffee"
	page := int32(0)
	pageSize := int32(1)
	category := domains.FoodDrinks
	tagIDPointer := uuidPointer(targetTagID)
	categoryPointerValue := itemCategoryPointer(category)

	firstCreate := &domains.ItemCreate{Name: firstName, Price: firstPrice, Category: foodCategory}
	secondCreate := &domains.ItemCreate{Name: secondName, Price: secondPrice, Category: foodCategory}
	thirdCreate := &domains.ItemCreate{Name: thirdName, Price: thirdPrice, Category: giftCategory}
	fourthCreate := &domains.ItemCreate{Name: fourthName, Price: fourthPrice, Category: foodCategory}

	firstID, firstCreateErr := repo.Create(ctx, firstCreate)
	secondID, secondCreateErr := repo.Create(ctx, secondCreate)
	thirdID, thirdCreateErr := repo.Create(ctx, thirdCreate)
	fourthID, fourthCreateErr := repo.Create(ctx, fourthCreate)
	_, tagInsertErr := testDB.Exec(tagInsertQuery, tagInsertArgs...)
	linkInsertArgs := []any{
		firstID, targetTagID,
		secondID, targetTagID,
		thirdID, targetTagID,
		fourthID, otherTagID,
	}
	_, linkInsertErr := testDB.Exec(linkInsertQuery, linkInsertArgs...)

	filter := &domains.ItemFilter{
		Name:       &filterName,
		Categories: []*domains.ItemCategory{categoryPointerValue},
		TagIds:     []*uuid.UUID{tagIDPointer},
		Page:       &page,
		PageSize:   &pageSize,
	}

	expectedNames := []string{firstName, secondName}

	// Act
	items, count, getErr := repo.Get(ctx, filter)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, thirdCreateErr)
	require.NoError(t, fourthCreateErr)
	require.NoError(t, tagInsertErr)
	require.NoError(t, linkInsertErr)
	require.NoError(t, getErr)
	require.Len(t, items, 1)
	assert.Equal(t, int64(2), count)
	assert.Contains(t, expectedNames, items[0].Name)
	assert.Equal(t, domains.FoodDrinks, items[0].Category)
}

func Test_ItemsRepository_GetById_ShouldReturnErrorOnNilID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	itemID := uuid.Nil

	// Act
	item, err := repo.GetById(ctx, itemID)

	// Assert
	require.EqualError(t, err, "id should not be nil")
	assert.Nil(t, item)
}

func Test_ItemsRepository_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	originalName := "Old"
	updatedName := "New"
	originalPrice := decimal.NewFromFloat(10.00)
	updatedPrice := decimal.NewFromFloat(15.50)
	itemCategory := "FoodDrinks"

	create := &domains.ItemCreate{
		Name:     originalName,
		Price:    originalPrice,
		Category: itemCategory,
	}

	update := &domains.ItemUpdate{
		Name:     updatedName,
		Price:    updatedPrice,
		Category: itemCategory,
	}

	// Act
	id, createErr := repo.Create(ctx, create)
	ok, updateErr := repo.Update(ctx, id, update)
	item, getErr := repo.GetById(ctx, id)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.NotNil(t, item)
	assert.Equal(t, updatedName, item.Name)
	assert.True(t, updatedPrice.Equal(item.Price))
}

func Test_ItemsRepository_Update_ShouldReturnFalseWhenItemDoesNotExist(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	itemName := "Missing"
	itemPrice := decimal.NewFromFloat(15.50)
	itemCategory := "FoodDrinks"

	update := &domains.ItemUpdate{
		Name:     itemName,
		Price:    itemPrice,
		Category: itemCategory,
	}

	// Act
	ok, err := repo.Update(ctx, itemID, update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsRepository_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "items")
	})

	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	itemName := "Book"
	itemPrice := decimal.NewFromFloat(10.00)
	itemCategory := "Entertainments"

	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    itemPrice,
		Category: itemCategory,
	}

	// Act
	id, createErr := repo.Create(ctx, create)
	ok, deleteErr := repo.Delete(ctx, id)
	item, getErr := repo.GetById(ctx, id)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, deleteErr)
	require.True(t, ok)
	assert.Nil(t, item)
	assert.Error(t, getErr)
}

func Test_ItemsRepository_Delete_ShouldReturnFalseWhenItemDoesNotExist(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewItemsRepository(testDB, testLogger)
	itemID := uuid.New()

	// Act
	ok, err := repo.Delete(ctx, itemID)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func uuidPointer(value uuid.UUID) *uuid.UUID {
	return &value
}

func itemCategoryPointer(value domains.ItemCategory) *domains.ItemCategory {
	return &value
}
