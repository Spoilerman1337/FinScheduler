//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/internal/features/services"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	itemsRepository := repositories.NewItemsRepository(testDB, testLogger)
	tagsRepository := repositories.NewTagsRepository(testDB, testLogger)
	tagToItemsRepository := repositories.NewTagToItemsRepository(testDB, testLogger)
	svc := services.NewItemsService(itemsRepository, tagsRepository, tagToItemsRepository, testLogger)

	create := &domains.ItemCreate{
		Name:     "Item",
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	item, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Item", *item.Name)
}

func Test_ItemsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	itemsRepository := repositories.NewItemsRepository(testDB, testLogger)
	tagsRepository := repositories.NewTagsRepository(testDB, testLogger)
	tagToItemsRepository := repositories.NewTagToItemsRepository(testDB, testLogger)
	svc := services.NewItemsService(itemsRepository, tagsRepository, tagToItemsRepository, testLogger)

	create := &domains.ItemCreate{
		Name:     "Ice",
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &domains.ItemUpdate{
		Name:     "Water",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Water", *item.Name)
	assert.Equal(t, 15.50, *item.Price)
}

func Test_ItemsService_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	itemsRepository := repositories.NewItemsRepository(testDB, testLogger)
	tagsRepository := repositories.NewTagsRepository(testDB, testLogger)
	tagToItemsRepository := repositories.NewTagToItemsRepository(testDB, testLogger)
	svc := services.NewItemsService(itemsRepository, tagsRepository, tagToItemsRepository, testLogger)

	create := &domains.ItemCreate{
		Name:     "Orange",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Delete(testContext, id)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := svc.GetById(testContext, id)

	// Assert
	assert.Nil(t, item)
	assert.Error(t, err)
}
