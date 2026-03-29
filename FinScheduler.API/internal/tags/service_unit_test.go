package tags

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_TagsService_Get_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewTagsService(nil, logger)
	ctx := context.Background()

	// Act
	_, _, err := svc.Get(ctx, nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_TagsService_GetById_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewTagsService(nil, logger)
	ctx := context.Background()
	nilUUID := uuid.Nil

	// Act
	_, err := svc.GetById(ctx, nilUUID)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")

}

func Test_TagsService_Create_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewTagsService(nil, logger)
	ctx := context.Background()

	// Act
	_, err := svc.Create(ctx, nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_TagsService_Update_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	logger := slog.Default()
	svc := NewTagsService(nil, logger)
	nilUUID := uuid.Nil
	notNilUUID, _ := uuid.NewV7()
	ctx := context.Background()

	// Act
	_, err1 := svc.Update(ctx, nilUUID, &TagUpdate{})
	_, err2 := svc.Update(ctx, notNilUUID, nil)

	// Assert
	assert.NotNilf(t, err1, "expected to get an error")
	assert.NotNilf(t, err2, "expected to get an error")
}
