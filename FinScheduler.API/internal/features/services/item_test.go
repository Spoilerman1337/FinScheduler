package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ItemsService_Get_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewItemsService(nil, logger)
	ctx := context.Background()

	// Act
	_, _, err := svc.Get(ctx, nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_ItemsService_Create_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewItemsService(nil, logger)
	ctx := context.Background()

	// Act
	_, err := svc.Create(ctx, nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_ItemsService_Update_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	t.Run("nil id", func(t *testing.T) {
		// Arrange
		logger := slog.Default()
		svc := NewItemsService(nil, logger)
		ctx := context.Background()

		// Act
		_, err := svc.Update(ctx, uuid.Nil, &domains.ItemUpdate{})

		// Assert
		assert.NotNilf(t, err, "expected to get an error")
	})

	t.Run("nil payload", func(t *testing.T) {
		// Arrange
		logger := slog.Default()
		svc := NewItemsService(nil, logger)
		notNilUUID, _ := uuid.NewV7()
		ctx := context.Background()

		// Act
		_, err := svc.Update(ctx, notNilUUID, nil)

		// Assert
		assert.NotNilf(t, err, "expected to get an error")
	})
}

func Test_ItemsService_Delete_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewItemsService(nil, logger)
	nilUUID := uuid.Nil
	ctx := context.Background()

	// Act
	_, err := svc.Delete(ctx, nilUUID)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}
