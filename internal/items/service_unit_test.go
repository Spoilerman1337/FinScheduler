package items

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ItemsService_Get_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	svc := NewItemsService(nil)

	// Act
	_, _, err := svc.Get(nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_ItemsService_GetById_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	svc := NewItemsService(nil)
	nilUUID := uuid.Nil

	// Act
	_, err := svc.GetById(nilUUID)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")

}

func Test_ItemsService_Create_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	svc := NewItemsService(nil)

	// Act
	_, err := svc.Create(nil)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}

func Test_ItemsService_Update_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	svc := NewItemsService(nil)
	nilUUID := uuid.Nil
	notNilUUID, _ := uuid.NewV7()

	// Act
	_, err1 := svc.Update(nilUUID, &ItemUpdate{})
	_, err2 := svc.Update(notNilUUID, nil)

	// Assert
	assert.NotNilf(t, err1, "expected to get an error")
	assert.NotNilf(t, err2, "expected to get an error")
}

func Test_ItemsService_Delete_ValidateInputs_ShouldReturnErrorOnNilInput(t *testing.T) {
	// Arrange
	svc := NewItemsService(nil)
	nilUUID := uuid.Nil

	// Act
	_, err := svc.Delete(nilUUID)

	// Assert
	assert.NotNilf(t, err, "expected to get an error")
}
