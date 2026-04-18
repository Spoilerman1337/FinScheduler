//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewItemsService(uow, testLogger)

	create := &domains.ItemCreate{
		Name:     "Item",
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	items, count, err := svc.Get(testContext, newItemFilterByID(id))
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Item", *items[0].Name)
}

func Test_ItemsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewItemsService(uow, testLogger)

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

	items, count, err := svc.Get(testContext, newItemFilterByID(id))
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Water", *items[0].Name)
	assert.Equal(t, 15.50, *items[0].Price)
}

func Test_ItemsService_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewItemsService(uow, testLogger)

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

	page := int32(0)
	pageSize := int32(20)

	var filter = domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}
	items, count, err := svc.Get(testContext, &filter)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, int64(0), count)
}

func Test_ItemsService_Create_ShouldRollbackItemWhenTagInsertFails(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewItemsService(uow, testLogger)

	create := &domains.ItemCreate{
		Name:     "Rollback",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
		TagIds:   []string{uuid.New().String()},
	}

	// Act
	id, err := svc.Create(testContext, create)

	// Assert
	require.Error(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE name = $1", create.Name))
	assert.Equal(t, 0, count)
}
