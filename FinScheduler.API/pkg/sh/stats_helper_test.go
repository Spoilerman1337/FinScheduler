package sh

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWinsorizeValue_ShouldClampToUpperBound(t *testing.T) {
	// Arrange
	value := decimal.RequireFromString("160.00")
	projectedValue := decimal.RequireFromString("100.00")
	windowPercent := decimal.RequireFromString("25.00")

	// Act
	winsorizedValue := WinsorizeValue(value, projectedValue, windowPercent)

	// Assert
	assert.True(t, decimal.RequireFromString("125.00").Equal(winsorizedValue))
}

func TestWinsorizeValue_ShouldClampToLowerBound(t *testing.T) {
	// Arrange
	value := decimal.RequireFromString("70.00")
	projectedValue := decimal.RequireFromString("100.00")
	windowPercent := decimal.RequireFromString("25.00")

	// Act
	winsorizedValue := WinsorizeValue(value, projectedValue, windowPercent)

	// Assert
	assert.True(t, decimal.RequireFromString("75.00").Equal(winsorizedValue))
}

func TestWinsorizeValue_ShouldReturnValueWithinCorridor(t *testing.T) {
	// Arrange
	value := decimal.RequireFromString("118.00")
	projectedValue := decimal.RequireFromString("100.00")
	windowPercent := decimal.RequireFromString("25.00")

	// Act
	winsorizedValue := WinsorizeValue(value, projectedValue, windowPercent)

	// Assert
	assert.True(t, decimal.RequireFromString("118.00").Equal(winsorizedValue))
}
