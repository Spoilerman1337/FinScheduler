//go:build integration
// +build integration

package persistence_test

import (
	"context"
	"errors"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitOfWorkWithTx_ShouldCommitChangesWhenCallbackSucceeds(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	itemName := "Committed Item"
	tagName := "Committed Tag"
	itemPrice := decimal.NewFromFloat(12.50)
	itemCategory := "FoodDrinks"
	itemCountQuery := "SELECT COUNT(*) FROM items WHERE id = $1"
	tagCountQuery := "SELECT COUNT(*) FROM tags WHERE id = $1"
	linkCountQuery := "SELECT COUNT(*) FROM tag_to_item WHERE item_id = $1 AND tag_id = $2"

	itemCreate := &domains.ItemCreate{
		Name:     itemName,
		Price:    itemPrice,
		Category: itemCategory,
	}

	tagCreate := &domains.TagCreate{
		Name:     tagName,
		IsActive: true,
	}

	var createdItemID uuid.UUID
	var createdTagID uuid.UUID

	// Act
	err := uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var createErr error

		createdItemID, createErr = repositories.Items.Create(ctx, itemCreate)
		if createErr != nil {
			return createErr
		}

		createdTagID, createErr = repositories.Tags.Create(ctx, tagCreate)
		if createErr != nil {
			return createErr
		}

		tagIDPointer := &createdTagID
		linkCreate := &domains.TagToItemCreate{
			ItemId: &createdItemID,
			TagIds: []*uuid.UUID{tagIDPointer},
		}

		_, createErr = repositories.TagToItems.BulkInsert(ctx, linkCreate)
		if createErr != nil {
			return createErr
		}

		return nil
	})

	var itemCount int
	var tagCount int
	var linkCount int
	itemCountErr := testDB.Get(&itemCount, itemCountQuery, createdItemID)
	tagCountErr := testDB.Get(&tagCount, tagCountQuery, createdTagID)
	linkCountErr := testDB.Get(&linkCount, linkCountQuery, createdItemID, createdTagID)

	// Assert
	require.NoError(t, err)
	require.NoError(t, itemCountErr)
	require.NoError(t, tagCountErr)
	require.NoError(t, linkCountErr)
	assert.NotEqual(t, uuid.Nil, createdItemID)
	assert.NotEqual(t, uuid.Nil, createdTagID)
	assert.Equal(t, 1, itemCount)
	assert.Equal(t, 1, tagCount)
	assert.Equal(t, 1, linkCount)
}

func TestUnitOfWorkWithTx_ShouldRollbackChangesWhenCallbackReturnsError(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	expectedErr := errors.New("rollback requested")
	itemName := "Rolled Back Item"
	tagName := "Rolled Back Tag"
	itemPrice := decimal.NewFromFloat(15.75)
	itemCategory := "FoodDrinks"
	itemCountQuery := "SELECT COUNT(*) FROM items WHERE id = $1"
	tagCountQuery := "SELECT COUNT(*) FROM tags WHERE id = $1"
	linkCountQuery := "SELECT COUNT(*) FROM tag_to_item WHERE item_id = $1 AND tag_id = $2"

	itemCreate := &domains.ItemCreate{
		Name:     itemName,
		Price:    itemPrice,
		Category: itemCategory,
	}

	tagCreate := &domains.TagCreate{
		Name:     tagName,
		IsActive: true,
	}

	var createdItemID uuid.UUID
	var createdTagID uuid.UUID

	// Act
	err := uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var createErr error

		createdItemID, createErr = repositories.Items.Create(ctx, itemCreate)
		if createErr != nil {
			return createErr
		}

		createdTagID, createErr = repositories.Tags.Create(ctx, tagCreate)
		if createErr != nil {
			return createErr
		}

		tagIDPointer := &createdTagID
		linkCreate := &domains.TagToItemCreate{
			ItemId: &createdItemID,
			TagIds: []*uuid.UUID{tagIDPointer},
		}

		_, createErr = repositories.TagToItems.BulkInsert(ctx, linkCreate)
		if createErr != nil {
			return createErr
		}

		return expectedErr
	})

	var itemCount int
	var tagCount int
	var linkCount int
	itemCountErr := testDB.Get(&itemCount, itemCountQuery, createdItemID)
	tagCountErr := testDB.Get(&tagCount, tagCountQuery, createdTagID)
	linkCountErr := testDB.Get(&linkCount, linkCountQuery, createdItemID, createdTagID)

	// Assert
	require.ErrorIs(t, err, expectedErr)
	require.NoError(t, itemCountErr)
	require.NoError(t, tagCountErr)
	require.NoError(t, linkCountErr)
	assert.NotEqual(t, uuid.Nil, createdItemID)
	assert.NotEqual(t, uuid.Nil, createdTagID)
	assert.Zero(t, itemCount)
	assert.Zero(t, tagCount)
	assert.Zero(t, linkCount)
}

func TestUnitOfWorkWithTx_ShouldReturnErrorWhenTransactionCannotStart(t *testing.T) {
	// Arrange
	ctx := context.Background()
	invalidConnectionString := "postgres://test:secret@127.0.0.1:1/testdb?sslmode=disable&connect_timeout=1"
	invalidDriverName := "pgx"
	unavailableDB := sqlx.MustOpen(invalidDriverName, invalidConnectionString)
	t.Cleanup(func() {
		_ = unavailableDB.Close()
	})

	uow := persistence.NewUnitOfWork(unavailableDB, testLogger)
	called := false

	// Act
	err := uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		called = true
		return nil
	})

	// Assert
	require.Error(t, err)
	assert.False(t, called)
}
