//go:build integration
// +build integration

package items

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsRepository_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE items CASCADE")
	})

	repo := NewItemsRepository(testDB, testLogger)

	create := &ItemCreate{
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
		testDB.Exec("TRUNCATE items CASCADE")
	})

	repo := NewItemsRepository(testDB, testLogger)

	create := &ItemCreate{
		Name:     "Old",
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &ItemUpdate{
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
		testDB.Exec("TRUNCATE items CASCADE")
	})

	repo := NewItemsRepository(testDB, testLogger)

	create := &ItemCreate{
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

func Test_ItemsRepository_CreateAndGet_ShouldReturnItemWithTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE items, tags, tag_to_item CASCADE")
	})

	repo := NewItemsRepository(testDB, testLogger)

	var tagId = uuid.New()

	_, err := testDB.Exec(
		`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3)`,
		tagId, "Fruit", true,
	)
	require.NoError(t, err)

	create := &ItemCreate{
		Name:     "Apple",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
		TagIds:   []string{tagId.String()},
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

	require.NotNil(t, item.Tags)
	assert.Len(t, item.Tags, 1)

	assert.Equal(t, tagId.String(), *item.Tags[0].Value)
	assert.Equal(t, "Fruit", *item.Tags[0].Label)
}

func Test_ItemsRepository_Update_ShouldReconcileTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE items, tags, tag_to_item CASCADE")
	})

	repo := NewItemsRepository(testDB, testLogger)

	tag1 := uuid.New()
	tag2 := uuid.New()
	tag3 := uuid.New()

	_, err := testDB.Exec(`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3)`, tag1, "Tag1", true)
	require.NoError(t, err)

	_, err = testDB.Exec(`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3)`, tag2, "Tag2", true)
	require.NoError(t, err)

	_, err = testDB.Exec(`INSERT INTO tags (id, name, is_active) VALUES ($1, $2, $3)`, tag3, "Tag3", true)
	require.NoError(t, err)

	create := &ItemCreate{
		Name:   "Old",
		Price:  decimal.NewFromFloat(10.00),
		TagIds: []string{tag1.String()},
	}

	update := &ItemUpdate{
		Name:   "New",
		Price:  decimal.NewFromFloat(15.50),
		TagIds: []string{tag2.String(), tag3.String()},
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
	actualTagIds := make([]string, len(item.Tags))
	for idx, tag := range item.Tags {
		actualTagIds[idx] = *tag.Value
	}

	assert.Len(t, actualTagIds, 2)
	assert.Contains(t, actualTagIds, tag2.String())
	assert.Contains(t, actualTagIds, tag3.String())
	assert.NotContains(t, actualTagIds, tag1.String())
}
