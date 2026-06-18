package sh

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWinsorize_ShouldClampUpperOutlierByIQR(t *testing.T) {
	// Arrange
	values := []decimal.Decimal{
		decimal.RequireFromString("200"),
		decimal.RequireFromString("500"),
		decimal.RequireFromString("500"),
		decimal.RequireFromString("600"),
		decimal.RequireFromString("700"),
		decimal.RequireFromString("700"),
		decimal.RequireFromString("800"),
		decimal.RequireFromString("900"),
		decimal.RequireFromString("1000"),
		decimal.RequireFromString("50000"),
	}

	// Act
	winsorizedValues := Winsorize(values)

	// Assert
	require.Len(t, winsorizedValues, len(values))
	assert.True(t, decimal.RequireFromString("1500").Equal(winsorizedValues[9]))
	assert.True(t, decimal.RequireFromString("200").Equal(winsorizedValues[0]))
}

func TestWinsorize_ShouldClampLowerOutlierToZeroBound(t *testing.T) {
	// Arrange
	values := []decimal.Decimal{
		decimal.RequireFromString("-50"),
		decimal.RequireFromString("500"),
		decimal.RequireFromString("500"),
		decimal.RequireFromString("600"),
		decimal.RequireFromString("700"),
		decimal.RequireFromString("700"),
		decimal.RequireFromString("800"),
		decimal.RequireFromString("900"),
		decimal.RequireFromString("1000"),
		decimal.RequireFromString("1200"),
	}

	// Act
	winsorizedValues := Winsorize(values)

	// Assert
	require.Len(t, winsorizedValues, len(values))
	assert.True(t, decimal.Zero.Equal(winsorizedValues[0]))
}

func TestWinsorize_ShouldReturnCopyWithoutMutatingInput(t *testing.T) {
	// Arrange
	values := []decimal.Decimal{
		decimal.RequireFromString("200"),
		decimal.RequireFromString("500"),
		decimal.RequireFromString("50000"),
	}
	originalValues := append([]decimal.Decimal(nil), values...)

	// Act
	winsorizedValues := Winsorize(values)

	// Assert
	assert.Equal(t, originalValues, values)
	assert.NotSame(t, &values[0], &winsorizedValues[0])
}
