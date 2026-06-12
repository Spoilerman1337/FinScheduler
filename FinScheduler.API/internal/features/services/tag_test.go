package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/persistence"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsServiceGetListingInfo_ShouldReturnErrorOnNilFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var filter *domains.TagFilter
	service := NewTagsService(uow, logger)

	// Act
	tags, count, err := service.GetListingInfo(ctx, filter)

	// Assert
	require.EqualError(t, err, "filter is nil")
	assert.Nil(t, tags)
	assert.Zero(t, count)
}

func TestTagsServiceGetLookup_ShouldReturnErrorOnNilFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var filter *domains.TagLookupFilter
	service := NewTagsService(uow, logger)

	// Act
	lookups, count, err := service.GetLookup(ctx, filter)

	// Assert
	require.EqualError(t, err, "filter is nil")
	assert.Nil(t, lookups)
	assert.Zero(t, count)
}

func TestTagsServiceCreate_ShouldReturnErrorOnNilCreate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var create *domains.TagCreate
	service := NewTagsService(uow, logger)

	// Act
	newID, err := service.Create(ctx, create)

	// Assert
	require.EqualError(t, err, "create is nil")
	assert.Equal(t, uuid.Nil, newID)
}

func TestTagsServiceUpdate_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	nilID := uuid.Nil
	validID := uuid.New()
	update := &domains.TagUpdate{}
	var nilUpdate *domains.TagUpdate
	service := NewTagsService(uow, logger)

	// Act
	successOnNilID, errOnNilID := service.Update(ctx, nilID, update)
	successOnNilUpdate, errOnNilUpdate := service.Update(ctx, validID, nilUpdate)

	// Assert
	require.EqualError(t, errOnNilID, "tagID is nil")
	require.EqualError(t, errOnNilUpdate, "update is nil")
	assert.False(t, successOnNilID)
	assert.False(t, successOnNilUpdate)
}
