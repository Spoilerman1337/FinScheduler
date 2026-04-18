//go:build integration
// +build integration

package repositories_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsRepository_CreateAndGetByIds_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	repo := repositories.NewTagsRepository(testDB, testLogger)

	create := &domains.TagCreate{
		Name: "Apple",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	tags, err := repo.GetByIds(testContext, []uuid.UUID{id})
	require.NoError(t, err)
	require.Len(t, tags, 1)

	// Assert
	assert.Equal(t, "Apple", tags[0].Name)
}

func Test_TagsRepository_GetByIds_ShouldReturnOnlyRequestedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	repo := repositories.NewTagsRepository(testDB, testLogger)

	firstID, err := repo.Create(testContext, &domains.TagCreate{Name: "Apple"})
	require.NoError(t, err)
	secondID, err := repo.Create(testContext, &domains.TagCreate{Name: "Orange"})
	require.NoError(t, err)

	// Act
	tags, err := repo.GetByIds(testContext, []uuid.UUID{firstID})

	// Assert
	require.NoError(t, err)
	require.Len(t, tags, 1)
	assert.Equal(t, firstID, tags[0].Id)
	assert.NotEqual(t, secondID, tags[0].Id)
}

func Test_TagsRepository_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	repo := repositories.NewTagsRepository(testDB, testLogger)

	create := &domains.TagCreate{
		Name: "Old",
	}

	update := &domains.TagUpdate{
		Name: "New",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)

	ok, err := repo.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	tags, err := repo.GetByIds(testContext, []uuid.UUID{id})
	require.NoError(t, err)
	require.Len(t, tags, 1)

	// Assert
	assert.Equal(t, "New", tags[0].Name)
}
