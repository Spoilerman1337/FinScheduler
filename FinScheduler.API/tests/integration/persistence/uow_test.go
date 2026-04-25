//go:build integration
// +build integration

package persistence_test

import (
	"errors"
	"testing"

	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UnitOfWork_WithoutTx_ShouldExposeRepositories(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)

	// Act
	err := uow.WithoutTx(func(repositories persistence.Repositories) error {
		_, err := repositories.Items.Create(testContext, &domains.ItemCreate{
			Name:     "WithoutTx",
			Category: string(domains.FoodDrinks),
		})
		return err
	})

	// Assert
	require.NoError(t, err)

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE name = $1", "WithoutTx"))
	assert.Equal(t, 1, count)
}

func Test_UnitOfWork_WithTx_ShouldCommitOnSuccess(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)

	// Act
	err := uow.WithTx(testContext, func(repositories persistence.Repositories) error {
		_, err := repositories.Tags.Create(testContext, &domains.TagCreate{
			Name:     "Committed",
			IsActive: true,
		})
		return err
	})

	// Assert
	require.NoError(t, err)

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM tags WHERE name = $1", "Committed"))
	assert.Equal(t, 1, count)
}

func Test_UnitOfWork_WithTx_ShouldRollbackOnCallbackError(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	expectedErr := errors.New("force rollback")

	// Act
	err := uow.WithTx(testContext, func(repositories persistence.Repositories) error {
		_, createErr := repositories.Items.Create(testContext, &domains.ItemCreate{
			Name:     "RolledBack",
			Category: string(domains.FoodDrinks),
		})
		require.NoError(t, createErr)
		return expectedErr
	})

	// Assert
	require.ErrorIs(t, err, expectedErr)

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE name = $1", "RolledBack"))
	assert.Equal(t, 0, count)
}
