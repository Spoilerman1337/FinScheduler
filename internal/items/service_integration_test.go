//go:build integration
// +build integration

package items

import (
	"github.com/shopspring/decimal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	repository := NewItemsRepository(testDB)
	svc := NewItemsService(repository)

	create := &ItemCreate{
		Name:  "Item",
		Price: decimal.NewFromFloat(15.50),
	}

	// Act
	id, err := svc.Create(create)
	require.NoError(t, err)

	item, err := svc.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Item", *item.Name)
	assert.Equal(t, 15.50, *item.Price)
}

func Test_ItemsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	svc := NewItemsService(NewItemsRepository(testDB))

	create := &ItemCreate{
		Name:  "Ice",
		Price: decimal.NewFromFloat(10.00),
	}

	update := &ItemUpdate{
		Name:  "Water",
		Price: decimal.NewFromFloat(15.50),
	}

	// Act
	id, err := svc.Create(create)
	require.NoError(t, err)

	ok, err := svc.Update(id, update)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := svc.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Water", *item.Name)
	assert.Equal(t, 15.50, *item.Price)
}

func Test_ItemsService_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	svc := NewItemsService(NewItemsRepository(testDB))

	create := &ItemCreate{
		Name:  "Orange",
		Price: decimal.NewFromFloat(15.50),
	}

	// Act
	id, err := svc.Create(create)
	require.NoError(t, err)

	ok, err := svc.Delete(id)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := svc.GetById(id)

	// Assert
	assert.Nil(t, item)
	assert.Error(t, err)
}
