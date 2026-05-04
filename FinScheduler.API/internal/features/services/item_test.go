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

func TestItemsServiceGet_ShouldReturnErrorOnNilFilter(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var filter *domains.ItemFilter
	service := NewItemsService(uow, logger)

	// Act
	items, count, err := service.Get(ctx, filter)

	// Assert
	require.EqualError(t, err, "filter is nil")
	assert.Nil(t, items)
	assert.Zero(t, count)
}

func TestItemsServiceCreate_ShouldReturnErrorOnNilCreate(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	var create *domains.ItemCreate
	service := NewItemsService(uow, logger)

	// Act
	newID, err := service.Create(ctx, create)

	// Assert
	require.EqualError(t, err, "create is nil")
	assert.Equal(t, uuid.Nil, newID)
}

func TestItemsServiceUpdate_ShouldReturnErrorOnInvalidInput(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	nilID := uuid.Nil
	validID := uuid.New()
	update := &domains.ItemUpdate{}
	var nilUpdate *domains.ItemUpdate
	service := NewItemsService(uow, logger)

	// Act
	successOnNilID, errOnNilID := service.Update(ctx, nilID, update)
	successOnNilUpdate, errOnNilUpdate := service.Update(ctx, validID, nilUpdate)

	// Assert
	require.EqualError(t, errOnNilID, "itemID is nil")
	require.EqualError(t, errOnNilUpdate, "update is nil")
	assert.False(t, successOnNilID)
	assert.False(t, successOnNilUpdate)
}

func TestItemsServiceDelete_ShouldReturnErrorOnNilItemID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	logger := slog.Default()
	var uow *persistence.UnitOfWork
	itemID := uuid.Nil
	service := NewItemsService(uow, logger)

	// Act
	success, err := service.Delete(ctx, itemID)

	// Assert
	require.EqualError(t, err, "itemID is nil")
	assert.False(t, success)
}
