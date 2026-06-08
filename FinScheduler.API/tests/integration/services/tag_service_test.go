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

func TestTagsServiceGetById_ShouldReturnCreatedTag(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewTagsService(uow, testLogger)
	tagName := "Groceries"
	tagIsActive := false
	create := &domains.TagCreate{
		Name:     tagName,
		IsActive: tagIsActive,
	}

	tagID, createErr := service.Create(ctx, create)

	// Act
	tag, getErr := service.GetById(ctx, tagID)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.NotNil(t, tag)
	assert.Equal(t, tagID, *tag.Id)
	assert.Equal(t, tagName, *tag.Name)
	assert.Equal(t, tagIsActive, *tag.IsActive)
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

func TestTagsServiceUpdate_ShouldRemoveItemLinksWhenTagBecomesInactive(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	tagsService := services.NewTagsService(uow, testLogger)
	itemsService := services.NewItemsService(uow, testLogger)
	tagName := "Groceries"
	itemName := "Milk"
	countQuery := `SELECT COUNT(*) FROM tag_to_item WHERE tag_id = $1`

	tagCreate := &domains.TagCreate{
		Name:     tagName,
		IsActive: true,
	}

	tagID, tagCreateErr := tagsService.Create(ctx, tagCreate)

	itemCreate := &domains.ItemCreate{
		Name:     itemName,
		Category: "FoodDrinks",
		TagIds:   []string{tagID.String()},
	}

	itemID, itemCreateErr := itemsService.Create(ctx, itemCreate)

	update := &domains.TagUpdate{
		Name:     tagName,
		IsActive: false,
	}

	// Act
	ok, updateErr := tagsService.Update(ctx, tagID, update)
	item, getItemErr := itemsService.GetById(ctx, itemID)

	var actualLinkCount int
	countErr := testDB.Get(&actualLinkCount, countQuery, tagID)

	// Assert
	require.NoError(t, tagCreateErr)
	require.NoError(t, itemCreateErr)
	require.NoError(t, updateErr)
	require.NoError(t, getItemErr)
	require.NoError(t, countErr)
	require.True(t, ok)
	require.NotNil(t, item)
	assert.Empty(t, item.Tags)
	assert.Zero(t, actualLinkCount)
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
