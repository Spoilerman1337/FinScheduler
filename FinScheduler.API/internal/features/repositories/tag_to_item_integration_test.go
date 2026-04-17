//go:build integration
// +build integration

package repositories

import (
	"finscheduler/internal/features/domains"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagToItemsRepository_BulkInsertAndGetByItemIds_ShouldNotErr(t *testing.T) {
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE items, tags, tag_to_item CASCADE")
	})

	repo := NewTagToItemsRepository(testDB, testLogger)

	itemID := uuid.New()
	tagID1 := uuid.New()
	tagID2 := uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`,
		itemID, "Apple", "FoodDrinks",
	)
	require.NoError(t, err)

	_, err = testDB.Exec(
		`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)`,
		tagID1, "Fruit", true, tagID2, "Food", true,
	)
	require.NoError(t, err)

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
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE items, tags, tag_to_item CASCADE")
	})

	repo := NewTagToItemsRepository(testDB, testLogger)

	itemID := uuid.New()
	tagID1 := uuid.New()
	tagID2 := uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO items (id, name, category) VALUES ($1, $2, $3)`,
		itemID, "Book", "Entertainments",
	)
	require.NoError(t, err)

	_, err = testDB.Exec(
		`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3), ($4, $5, $6)`,
		tagID1, "Paper", true, tagID2, "Gift", true,
	)
	require.NoError(t, err)

	_, err = repo.BulkInsert(testContext, &domains.TagToItemCreate{
		ItemId: &itemID,
		TagIds: []*uuid.UUID{
			&tagID1,
			&tagID2,
		},
	})
	require.NoError(t, err)

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

	require.Len(t, tagToItems, 1)
	assert.Equal(t, itemID, tagToItems[0].ItemId)
	assert.Equal(t, tagID2, tagToItems[0].TagId)
}
