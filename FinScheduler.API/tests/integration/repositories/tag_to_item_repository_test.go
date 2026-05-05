//go:build integration
// +build integration

package repositories_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagToItemsRepositoryGetByItemIDs_ShouldReturnErrorOnNilItemIDs(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	var itemIDs []uuid.UUID

	// Act
	tagToItems, err := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.EqualError(t, err, "itemId should not be nil")
	assert.Nil(t, tagToItems)
}

func TestTagToItemsRepositoryGetByItemIDs_ShouldReturnEmptySliceOnEmptyItemIDs(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemIDs := make([]uuid.UUID, 0)

	// Act
	tagToItems, err := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, tagToItems)
	assert.Empty(t, tagToItems)
}

func TestTagToItemsRepositoryGetByItemIDs_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewTagToItemsRepository(closedDB, testLogger)
	itemIDs := []uuid.UUID{uuid.New()}

	// Act
	tagToItems, err := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.Error(t, err)
	assert.Nil(t, tagToItems)
}

func TestTagToItemsRepositoryBulkInsertAndGetByItemIDs_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	firstTagID := uuid.New()
	secondTagID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Apple", "FoodDrinks"}
	tagInsertQuery := `INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)`
	tagInsertArgs := []any{firstTagID, "Fruit", true, secondTagID, "Food", true}
	firstTagIDPointer := &firstTagID
	secondTagIDPointer := &secondTagID
	create := &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{firstTagIDPointer, secondTagIDPointer},
	}
	itemIDs := []uuid.UUID{itemID}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, tagInsertErr := testDB.Exec(tagInsertQuery, tagInsertArgs...)

	// Act
	ok, insertErr := repo.BulkInsert(ctx, create)
	tagToItems, getErr := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, tagInsertErr)
	require.NoError(t, insertErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.Len(t, tagToItems, 2)

	actualTagIDs := make([]uuid.UUID, 0, len(tagToItems))
	for _, tagToItem := range tagToItems {
		assert.Equal(t, itemID, tagToItem.ItemId)
		actualTagIDs = append(actualTagIDs, tagToItem.TagId)
	}

	assert.Contains(t, actualTagIDs, firstTagID)
	assert.Contains(t, actualTagIDs, secondTagID)
}

func TestTagToItemsRepositoryGetByItemIDs_ShouldReturnOnlyRequestedItemLinks(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	requestedItemID := uuid.New()
	otherItemID := uuid.New()
	firstRequestedTagID := uuid.New()
	secondRequestedTagID := uuid.New()
	otherTagID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3), ($4, $5, $6)`
	itemInsertArgs := []any{requestedItemID, "Apple", "FoodDrinks", otherItemID, "Book", "Entertainments"}
	tagInsertQuery := `INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9)`
	tagInsertArgs := []any{
		firstRequestedTagID, "Fruit", true,
		secondRequestedTagID, "Fresh", true,
		otherTagID, "Gift", true,
	}
	linkInsertQuery := `INSERT INTO tag_to_item (item_id, tag_id) VALUES ($1, $2), ($3, $4), ($5, $6)`
	linkInsertArgs := []any{
		requestedItemID, firstRequestedTagID,
		requestedItemID, secondRequestedTagID,
		otherItemID, otherTagID,
	}
	itemIDs := []uuid.UUID{requestedItemID}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, tagInsertErr := testDB.Exec(tagInsertQuery, tagInsertArgs...)
	_, linkInsertErr := testDB.Exec(linkInsertQuery, linkInsertArgs...)

	// Act
	tagToItems, getErr := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, tagInsertErr)
	require.NoError(t, linkInsertErr)
	require.NoError(t, getErr)
	require.Len(t, tagToItems, 2)

	actualTagIDs := make([]uuid.UUID, 0, len(tagToItems))
	for _, tagToItem := range tagToItems {
		assert.Equal(t, requestedItemID, tagToItem.ItemId)
		actualTagIDs = append(actualTagIDs, tagToItem.TagId)
	}

	assert.Contains(t, actualTagIDs, firstRequestedTagID)
	assert.Contains(t, actualTagIDs, secondRequestedTagID)
	assert.NotContains(t, actualTagIDs, otherTagID)
}

func TestTagToItemsRepositoryBulkInsert_ShouldReturnErrorOnInvalidReference(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	missingTagID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Apple", "FoodDrinks"}
	countQuery := `SELECT COUNT(*) FROM tag_to_item WHERE item_id = $1`
	missingTagIDPointer := &missingTagID
	create := &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{missingTagIDPointer},
	}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)

	// Act
	ok, insertErr := repo.BulkInsert(ctx, create)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.Error(t, insertErr)
	require.NoError(t, countErr)
	assert.False(t, ok)
	assert.Zero(t, actualCount)
}

func TestTagToItemsRepositoryBulkInsert_ShouldReturnErrorOnEmptyTagIDs(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Apple", "FoodDrinks"}
	create := &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{},
	}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)

	// Act
	ok, err := repo.BulkInsert(ctx, create)

	// Assert
	require.NoError(t, itemInsertErr)
	require.Error(t, err)
	assert.False(t, ok)
}

func TestTagToItemsRepositoryBulkDelete_ShouldRemoveOnlySelectedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	firstTagID := uuid.New()
	secondTagID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Book", "Entertainments"}
	tagInsertQuery := `INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)`
	tagInsertArgs := []any{firstTagID, "Paper", true, secondTagID, "Gift", true}
	firstTagIDPointer := &firstTagID
	secondTagIDPointer := &secondTagID
	create := &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{firstTagIDPointer, secondTagIDPointer},
	}
	deleteInput := &domains.TagToItemDelete{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{firstTagIDPointer},
	}
	itemIDs := []uuid.UUID{itemID}

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, tagInsertErr := testDB.Exec(tagInsertQuery, tagInsertArgs...)
	_, insertErr := repo.BulkInsert(ctx, create)

	// Act
	ok, deleteErr := repo.BulkDelete(ctx, deleteInput)
	tagToItems, getErr := repo.GetByItemIds(ctx, itemIDs)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, tagInsertErr)
	require.NoError(t, insertErr)
	require.NoError(t, deleteErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.Len(t, tagToItems, 1)
	assert.Equal(t, itemID, tagToItems[0].ItemId)
	assert.Equal(t, secondTagID, tagToItems[0].TagId)
}

func TestTagToItemsRepositoryBulkDelete_ShouldReturnFalseWhenNothingDeleted(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	existingTagID := uuid.New()
	missingTagID := uuid.New()
	itemInsertQuery := `INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`
	itemInsertArgs := []any{itemID, "Milk", "FoodDrinks"}
	tagInsertQuery := `INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)`
	tagInsertArgs := []any{existingTagID, "Dairy", true, missingTagID, "Unused", true}
	linkInsertQuery := `INSERT INTO tag_to_item (item_id, tag_id) VALUES ($1, $2)`
	linkInsertArgs := []any{itemID, existingTagID}
	missingTagIDPointer := &missingTagID
	deleteInput := &domains.TagToItemDelete{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{missingTagIDPointer},
	}
	countQuery := `SELECT COUNT(*) FROM tag_to_item WHERE item_id = $1`

	_, itemInsertErr := testDB.Exec(itemInsertQuery, itemInsertArgs...)
	_, tagInsertErr := testDB.Exec(tagInsertQuery, tagInsertArgs...)
	_, linkInsertErr := testDB.Exec(linkInsertQuery, linkInsertArgs...)

	// Act
	ok, deleteErr := repo.BulkDelete(ctx, deleteInput)
	var actualCount int
	countErr := testDB.Get(&actualCount, countQuery, itemID)

	// Assert
	require.NoError(t, itemInsertErr)
	require.NoError(t, tagInsertErr)
	require.NoError(t, linkInsertErr)
	require.NoError(t, deleteErr)
	require.NoError(t, countErr)
	assert.False(t, ok)
	assert.Equal(t, 1, actualCount)
}

func TestTagToItemsRepositoryBulkDelete_ShouldReturnFalseOnEmptyTagIDs(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	deleteInput := &domains.TagToItemDelete{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{},
	}

	// Act
	ok, err := repo.BulkDelete(ctx, deleteInput)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestTagToItemsRepositoryBulkDelete_ShouldPanicOnNilItemID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	tagID := uuid.New()
	tagIDPointer := &tagID
	deleteInput := &domains.TagToItemDelete{
		ItemId: nil,
		TagIds: []*uuid.UUID{tagIDPointer},
	}

	// Act
	action := func() {
		_, _ = repo.BulkDelete(ctx, deleteInput)
	}

	// Assert
	assert.Panics(t, action)
}
