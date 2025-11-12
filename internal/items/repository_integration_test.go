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
	repo := NewItemsRepository(testDB)

	create := &ItemCreate{
		Name:  "Apple",
		Price: decimal.NewFromFloat(15.50),
	}

	// Act
	id, err := repo.Create(create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	item, err := repo.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Apple", item.Name)
	assert.True(t, decimal.NewFromFloat(15.50).Equal(item.Price))
}

func Test_ItemsRepository_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	repo := NewItemsRepository(testDB)

	create := &ItemCreate{
		Name:  "Old",
		Price: decimal.NewFromFloat(10.00),
	}

	update := &ItemUpdate{
		Name:  "New",
		Price: decimal.NewFromFloat(15.50),
	}

	// Act
	id, err := repo.Create(create)
	require.NoError(t, err)

	ok, err := repo.Update(id, update)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := repo.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "New", item.Name)
	assert.True(t, decimal.NewFromFloat(15.50).Equal(item.Price))
}

func Test_ItemsRepository_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	repo := NewItemsRepository(testDB)

	create := &ItemCreate{
		Name:  "Book",
		Price: decimal.NewFromFloat(10.00),
	}

	// Act
	id, err := repo.Create(create)
	require.NoError(t, err)

	ok, err := repo.Delete(id)
	require.NoError(t, err)
	require.True(t, ok)

	item, err := repo.GetById(id)

	// Assert
	assert.Nil(t, item)
	assert.Error(t, err)
}
