//go:build integration
// +build integration

package repositories

import (
	"finscheduler/internal/features/domains"
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsRepository_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE tags CASCADE")
	})

	repo := NewTagsRepository(testDB, testLogger)

	create := &domains.TagCreate{
		Name: "Apple",
	}

	// Act
	id, err := repo.Create(testContext, create)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)

	tag, err := repo.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Apple", tag.Name)
}

func Test_TagsRepository_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testDB.Exec("TRUNCATE tags CASCADE")
	})

	repo := NewTagsRepository(testDB, testLogger)

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

	tag, err := repo.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "New", tag.Name)
}
