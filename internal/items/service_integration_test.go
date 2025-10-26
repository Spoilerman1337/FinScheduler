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
	repo := NewItemsRepository(testDB)
	svc := NewItemsService(repo)

	// Act
	id, err := svc.Create(&ItemCreate{Name: "Apple", Price: decimal.NewFromFloat(3.2)})
	require.NoError(t, err)

	got, err := svc.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Apple", *got.Name)
}

func Test_ItemsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	svc := NewItemsService(NewItemsRepository(testDB))

	// Act
	id, err := svc.Create(&ItemCreate{
		Name:  "Book",
		Price: decimal.NewFromFloat(10),
	})
	require.NoError(t, err)

	ok, err := svc.Update(id, &ItemUpdate{
		Name:  "Notebook",
		Price: decimal.NewFromFloat(15.5),
	})
	require.NoError(t, err)
	require.True(t, ok)

	got, err := svc.GetById(id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Notebook", *got.Name)
}

func Test_ItemsService_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	svc := NewItemsService(NewItemsRepository(testDB))

	// Act
	id, err := svc.Create(&ItemCreate{
		Name:  "Book",
		Price: decimal.NewFromFloat(10.0),
	})
	require.NoError(t, err)

	ok, err := svc.Delete(id)
	require.NoError(t, err)
	require.True(t, ok)

	got, err := svc.GetById(id)

	// Assert
	assert.Nil(t, got)
	assert.Error(t, err)
}
