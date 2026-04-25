//go:build integration
// +build integration

package repositories_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/pkg/dh"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagToItemsRepository_BulkInsertAndGetByItemIds_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewTagToItemsRepository(testDB, testLogger)

	itemID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Apple", Category: string(domains.FoodDrinks)})
	tagID1 := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Fruit", IsActive: true})
	tagID2 := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food", IsActive: true})

	// Act
	ok, err := repo.BulkInsert(testContext, &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{
			&tagID1,
			&tagID2,
		},
	})
	require.NoError(t, err)
	require.True(t, ok)

	tagToItems, err := repo.GetByItemIds(testContext, []uuid.UUID{itemID})
	require.NoError(t, err)

	// Assert
	actualTagIds := make([]uuid.UUID, 0, len(tagToItems))
	for _, tagToItem := range tagToItems {
		assert.Equal(t, itemID, tagToItem.ItemId)
		actualTagIds = append(actualTagIds, tagToItem.TagId)
	}

	assert.Len(t, actualTagIds, 2)
	assert.Contains(t, actualTagIds, tagID1)
	assert.Contains(t, actualTagIds, tagID2)
}

func Test_TagToItemsRepository_BulkDelete_ShouldRemoveOnlySelectedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewTagToItemsRepository(testDB, testLogger)

	itemID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Book", Category: string(domains.Entertainments)})
	tagID1 := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Paper", IsActive: true})
	tagID2 := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Gift", IsActive: true})
	testFixtures.MustLinkItemTags(t, itemID, tagID1, tagID2)

	// Act
	ok, err := repo.BulkDelete(testContext, &domains.TagToItemDelete{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{
			&tagID1,
		},
	})
	require.NoError(t, err)
	require.True(t, ok)

	tagToItems, err := repo.GetByItemIds(testContext, []uuid.UUID{itemID})
	require.NoError(t, err)

	// Assert
	require.Len(t, tagToItems, 1)
	assert.Equal(t, itemID, tagToItems[0].ItemId)
	assert.Equal(t, tagID2, tagToItems[0].TagId)
}

func Test_TagToItemsRepository_GetByItemIds_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)

	// Act
	tagToItems, err := repo.GetByItemIds(testContext, nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, tagToItems)
	assert.Contains(t, err.Error(), "itemId should not be nil")
}

func Test_TagToItemsRepository_GetByItemIds_ShouldReturnEmptySliceOnEmptyInput(t *testing.T) {
	// Arrange
	repo := repositories.NewTagToItemsRepository(testDB, testLogger)

	// Act
	tagToItems, err := repo.GetByItemIds(testContext, []uuid.UUID{})

	// Assert
	require.NoError(t, err)
	assert.Empty(t, tagToItems)
}

func Test_TagToItemsRepository_BulkInsert_ShouldReturnForeignKeyViolationForMissingReferences(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := uuid.New()
	tagID := uuid.New()

	// Act
	ok, err := repo.BulkInsert(testContext, &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{&tagID},
	})

	// Assert
	require.Error(t, err)
	assert.False(t, ok)
	assert.True(t, dh.IsPostgresForeignKeyViolation(err))
}

func Test_TagToItemsRepository_BulkDelete_ShouldReturnFalseWhenNothingWasDeleted(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	repo := repositories.NewTagToItemsRepository(testDB, testLogger)
	itemID := testFixtures.MustCreateItem(t, &domains.ItemCreate{Name: "Book", Category: string(domains.Entertainments)})
	tagID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Gift", IsActive: true})

	// Act
	ok, err := repo.BulkDelete(testContext, &domains.TagToItemDelete{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{&tagID},
	})

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}
