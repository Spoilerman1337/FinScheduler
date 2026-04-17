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
