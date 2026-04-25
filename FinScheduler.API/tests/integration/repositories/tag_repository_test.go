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

	firstID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Apple"})
	secondID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Orange"})

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

func Test_TagsRepository_Get_ShouldApplyFiltersAndCount(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	repo := repositories.NewTagsRepository(testDB, testLogger)

	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food Delivery", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Travel", IsActive: false})

	page := int32(0)
	pageSize := int32(10)
	name := "Food"
	isActive := true

	// Act
	tags, count, err := repo.Get(testContext, &domains.TagFilter{
		Name:     &name,
		IsActive: &isActive,
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, tags, 2)
	assert.Equal(t, int64(2), count)
	assert.ElementsMatch(t, []string{"Food", "Food Delivery"}, []string{tags[0].Name, tags[1].Name})
}

func Test_TagsRepository_GetByIds_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	repo := repositories.NewTagsRepository(testDB, testLogger)

	// Act
	tags, err := repo.GetByIds(testContext, nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, tags)
	assert.Contains(t, err.Error(), "ids should not be nil")
}

func Test_TagsRepository_GetByIds_ShouldReturnEmptySliceOnEmptyInput(t *testing.T) {
	// Arrange
	repo := repositories.NewTagsRepository(testDB, testLogger)

	// Act
	tags, err := repo.GetByIds(testContext, []uuid.UUID{})

	// Assert
	require.NoError(t, err)
	assert.Empty(t, tags)
}

func Test_TagsRepository_GetLookup_ShouldReturnOnlyActiveTagsOrderedByLowerName(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	repo := repositories.NewTagsRepository(testDB, testLogger)

	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "banana", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Apple", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "carrot", IsActive: false})

	page := int32(0)
	pageSize := int32(10)

	// Act
	lookups, count, err := repo.GetLookup(testContext, &domains.TagLookupFilter{
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, lookups, 2)
	assert.Equal(t, int64(2), count)
	assert.Equal(t, "Apple", *lookups[0].Label)
	assert.Equal(t, "banana", *lookups[1].Label)
}

func Test_TagsRepository_Update_ShouldReturnFalseForMissingRow(t *testing.T) {
	// Arrange
	repo := repositories.NewTagsRepository(testDB, testLogger)

	// Act
	ok, err := repo.Update(testContext, uuid.New(), &domains.TagUpdate{
		Name:     "Missing",
		IsActive: true,
	})

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}
