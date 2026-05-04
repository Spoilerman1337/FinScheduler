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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsServiceCreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewTagsService(uow, testLogger)
	tagName := "Tag"
	tagIsActive := true
	page := int32(0)
	pageSize := int32(20)

	create := &domains.TagCreate{
		Name:     tagName,
		IsActive: tagIsActive,
	}

	// Act
	tagID, createErr := service.Create(ctx, create)
	tagIDPointer := &tagID
	filter := &domains.TagFilter{
		Ids:      []*uuid.UUID{tagIDPointer},
		Page:     &page,
		PageSize: &pageSize,
	}

	tags, count, getErr := service.Get(ctx, filter)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.Len(t, tags, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, tagName, *tags[0].Name)
	assert.Equal(t, tagIsActive, *tags[0].IsActive)
}

func TestTagsServiceGetLookup_ShouldReturnOnlyActiveTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewTagsService(uow, testLogger)
	activeName := "Groceries"
	inactiveName := "Archived"
	page := int32(0)
	pageSize := int32(20)

	activeCreate := &domains.TagCreate{
		Name:     activeName,
		IsActive: true,
	}

	inactiveCreate := &domains.TagCreate{
		Name:     inactiveName,
		IsActive: false,
	}

	_, activeCreateErr := service.Create(ctx, activeCreate)
	_, inactiveCreateErr := service.Create(ctx, inactiveCreate)

	filter := &domains.TagLookupFilter{
		Page:     &page,
		PageSize: &pageSize,
	}

	// Act
	lookups, count, getErr := service.GetLookup(ctx, filter)

	// Assert
	require.NoError(t, activeCreateErr)
	require.NoError(t, inactiveCreateErr)
	require.NoError(t, getErr)
	require.Len(t, lookups, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, activeName, *lookups[0].Label)
}

func TestTagsServiceUpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewTagsService(uow, testLogger)
	originalName := "Ice"
	updatedName := "Water"
	originalIsActive := false
	updatedIsActive := true
	page := int32(0)
	pageSize := int32(20)

	create := &domains.TagCreate{
		Name:     originalName,
		IsActive: originalIsActive,
	}

	update := &domains.TagUpdate{
		Name:     updatedName,
		IsActive: updatedIsActive,
	}

	// Act
	tagID, createErr := service.Create(ctx, create)
	ok, updateErr := service.Update(ctx, tagID, update)
	tagIDPointer := &tagID
	filter := &domains.TagFilter{
		Ids:      []*uuid.UUID{tagIDPointer},
		Page:     &page,
		PageSize: &pageSize,
	}

	tags, count, getErr := service.Get(ctx, filter)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	require.True(t, ok)
	require.Len(t, tags, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, updatedName, *tags[0].Name)
	assert.Equal(t, updatedIsActive, *tags[0].IsActive)
}

func TestTagsServiceUpdateMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewTagsService(uow, testLogger)
	missingID := uuid.New()
	tagName := "Missing"
	tagIsActive := true

	update := &domains.TagUpdate{
		Name:     tagName,
		IsActive: tagIsActive,
	}

	// Act
	ok, err := service.Update(ctx, missingID, update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}
