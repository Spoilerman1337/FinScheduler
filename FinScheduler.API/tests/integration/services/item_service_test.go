//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

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
	assert.Equal(t, tagName, item.Tags[0].Label)
	assert.Equal(t, tagID.String(), item.Tags[0].Value)
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
