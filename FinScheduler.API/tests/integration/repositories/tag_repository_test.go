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

func TestTagsRepositoryCreateAndGetByIDs_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	tagName := "Apple"
	tagIsActive := true

	create := &domains.TagCreate{
		Name:     tagName,
		IsActive: tagIsActive,
	}

	// Act
	tagID, createErr := repo.Create(ctx, create)
	ids := []uuid.UUID{tagID}
	tags, getErr := repo.GetByIds(ctx, ids)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.NotEqual(t, uuid.Nil, tagID)
	require.Len(t, tags, 1)
	assert.Equal(t, tagName, tags[0].Name)
	assert.Equal(t, tagIsActive, tags[0].IsActive)
}

func TestTagsRepositoryGet_ShouldFilterAndReturnCount(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	firstName := "Groceries"
	secondName := "Ground Coffee"
	thirdName := "Grocery Archive"
	fourthName := "Travel"
	activeValue := true
	inactiveValue := false
	filterName := "Gro"
	page := int32(0)
	pageSize := int32(1)

	firstCreate := &domains.TagCreate{Name: firstName, IsActive: activeValue}
	secondCreate := &domains.TagCreate{Name: secondName, IsActive: activeValue}
	thirdCreate := &domains.TagCreate{Name: thirdName, IsActive: inactiveValue}
	fourthCreate := &domains.TagCreate{Name: fourthName, IsActive: activeValue}

	_, firstCreateErr := repo.Create(ctx, firstCreate)
	_, secondCreateErr := repo.Create(ctx, secondCreate)
	_, thirdCreateErr := repo.Create(ctx, thirdCreate)
	_, fourthCreateErr := repo.Create(ctx, fourthCreate)

	filter := &domains.TagFilter{
		Name:     &filterName,
		IsActive: &activeValue,
		Page:     &page,
		PageSize: &pageSize,
	}

	expectedNames := []string{firstName, secondName}

	// Act
	tags, count, getErr := repo.Get(ctx, filter)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, thirdCreateErr)
	require.NoError(t, fourthCreateErr)
	require.NoError(t, getErr)
	require.Len(t, tags, 1)
	assert.Equal(t, int64(2), count)
	assert.Contains(t, expectedNames, tags[0].Name)
	assert.True(t, tags[0].IsActive)
}

func TestTagsRepositoryGetByIDs_ShouldReturnOnlyRequestedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	firstName := "Apple"
	secondName := "Orange"

	firstCreate := &domains.TagCreate{Name: firstName, IsActive: true}
	secondCreate := &domains.TagCreate{Name: secondName, IsActive: true}

	firstID, firstCreateErr := repo.Create(ctx, firstCreate)
	secondID, secondCreateErr := repo.Create(ctx, secondCreate)
	ids := []uuid.UUID{firstID}

	// Act
	tags, getErr := repo.GetByIds(ctx, ids)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, getErr)
	require.Len(t, tags, 1)
	assert.Equal(t, firstID, tags[0].Id)
	assert.NotEqual(t, secondID, tags[0].Id)
}

func TestTagsRepositoryGetByIDs_ShouldReturnErrorOnNilIDs(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	var ids []uuid.UUID

	// Act
	tags, err := repo.GetByIds(ctx, ids)

	// Assert
	require.EqualError(t, err, "ids should not be nil")
	assert.Nil(t, tags)
}

func TestTagsRepositoryGetByIDs_ShouldReturnEmptySliceOnEmptyIDs(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	ids := make([]uuid.UUID, 0)

	// Act
	tags, err := repo.GetByIds(ctx, ids)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, tags)
	assert.Empty(t, tags)
}

func TestTagsRepositoryGetLookup_ShouldReturnActiveMatchesOrderedAndCount(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	firstName := "Alpha"
	secondName := "Alpine"
	thirdName := "Alphabet Soup"
	fourthName := "Travel"
	filterName := "Alp"
	page := int32(0)
	pageSize := int32(20)
	activeValue := true
	inactiveValue := false

	firstCreate := &domains.TagCreate{Name: firstName, IsActive: activeValue}
	secondCreate := &domains.TagCreate{Name: secondName, IsActive: activeValue}
	thirdCreate := &domains.TagCreate{Name: thirdName, IsActive: inactiveValue}
	fourthCreate := &domains.TagCreate{Name: fourthName, IsActive: activeValue}

	_, firstCreateErr := repo.Create(ctx, firstCreate)
	_, secondCreateErr := repo.Create(ctx, secondCreate)
	_, thirdCreateErr := repo.Create(ctx, thirdCreate)
	_, fourthCreateErr := repo.Create(ctx, fourthCreate)

	filter := &domains.TagLookupFilter{
		Name:     &filterName,
		Page:     &page,
		PageSize: &pageSize,
	}

	expectedLabels := []string{firstName, secondName}

	// Act
	lookups, count, getErr := repo.GetLookup(ctx, filter)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, thirdCreateErr)
	require.NoError(t, fourthCreateErr)
	require.NoError(t, getErr)
	require.Len(t, lookups, 2)
	assert.Equal(t, int64(2), count)
	assert.Equal(t, expectedLabels[0], *lookups[0].Label)
	assert.Equal(t, expectedLabels[1], *lookups[1].Label)
}

func TestTagsRepositoryUpdateAndGetByIDs_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	originalName := "Old"
	updatedName := "New"
	originalIsActive := false
	updatedIsActive := true

	create := &domains.TagCreate{
		Name:     originalName,
		IsActive: originalIsActive,
	}

	update := &domains.TagUpdate{
		Name:     updatedName,
		IsActive: updatedIsActive,
	}

	// Act
	tagID, createErr := repo.Create(ctx, create)
	ok, updateErr := repo.Update(ctx, tagID, update)
	ids := []uuid.UUID{tagID}
	tags, getErr := repo.GetByIds(ctx, ids)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.Len(t, tags, 1)
	assert.Equal(t, updatedName, tags[0].Name)
	assert.Equal(t, updatedIsActive, tags[0].IsActive)
}

func TestTagsRepositoryUpdate_ShouldReturnFalseWhenTagDoesNotExist(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewTagsRepository(testDB, testLogger)
	tagID := uuid.New()
	tagName := "Missing"
	tagIsActive := true

	update := &domains.TagUpdate{
		Name:     tagName,
		IsActive: tagIsActive,
	}

	// Act
	ok, err := repo.Update(ctx, tagID, update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}
