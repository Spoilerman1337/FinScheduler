//go:build integration
// +build integration

package items

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	repository := NewItemsRepository(testDB, testLogger)
	svc := NewItemsService(repository, testLogger)

	create := &ItemCreate{
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
	repository := NewItemsRepository(testDB, testLogger)
	svc := NewItemsService(repository, testLogger)

	create := &ItemCreate{
		Name:     "Ice",
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &ItemUpdate{
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
	repository := NewItemsRepository(testDB, testLogger)
	svc := NewItemsService(repository, testLogger)

	create := &ItemCreate{
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
