package domains

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPriceHistoryPointDto_ShouldMapGraphFieldsAndChangesAgainstPreviousPoint(t *testing.T) {
	// Arrange
	recordedAt := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedValue := decimal.RequireFromString("15.75")
	previousValue := decimal.RequireFromString("12.60")
	expectedAbsoluteChange := decimal.RequireFromString("3.15")
	expectedPercentChange := decimal.RequireFromString("25")
	priceHistory := PriceHistory{
		Id:         uuid.New(),
		ItemId:     uuid.New(),
		RecordedAt: recordedAt,
		Value:      expectedValue,
	}
	previousPriceHistory := &PriceHistory{
		Id:         uuid.New(),
		ItemId:     priceHistory.ItemId,
		RecordedAt: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		Value:      previousValue,
	}

	// Act
	dto := NewPriceHistoryPointDto(priceHistory, previousPriceHistory)

	// Assert
	require.NotNil(t, dto)
	assert.Equal(t, recordedAt, dto.Point)
	assert.True(t, expectedValue.Equal(dto.Value))
	require.NotNil(t, dto.AbsoluteChange)
	require.NotNil(t, dto.PercentChange)
	assert.True(t, expectedAbsoluteChange.Equal(*dto.AbsoluteChange))
	assert.True(t, expectedPercentChange.Equal(*dto.PercentChange))
}

func TestNewPriceHistoryPointDto_ShouldLeaveChangesNilWhenPreviousPointIsMissing(t *testing.T) {
	// Arrange
	recordedAt := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	expectedValue := decimal.RequireFromString("15.75")
	priceHistory := PriceHistory{
		Id:         uuid.New(),
		ItemId:     uuid.New(),
		RecordedAt: recordedAt,
		Value:      expectedValue,
	}

	// Act
	dto := NewPriceHistoryPointDto(priceHistory, nil)

	// Assert
	require.NotNil(t, dto)
	assert.Equal(t, recordedAt, dto.Point)
	assert.True(t, expectedValue.Equal(dto.Value))
	assert.Nil(t, dto.AbsoluteChange)
	assert.Nil(t, dto.PercentChange)
}

func TestNewPriceHistoryPointDto_ShouldLeavePercentChangeNilWhenPreviousValueIsZero(t *testing.T) {
	// Arrange
	priceHistory := PriceHistory{
		Id:         uuid.New(),
		ItemId:     uuid.New(),
		RecordedAt: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		Value:      decimal.RequireFromString("15.75"),
	}
	previousPriceHistory := &PriceHistory{
		Id:         uuid.New(),
		ItemId:     priceHistory.ItemId,
		RecordedAt: time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		Value:      decimal.Zero,
	}
	expectedAbsoluteChange := decimal.RequireFromString("15.75")

	// Act
	dto := NewPriceHistoryPointDto(priceHistory, previousPriceHistory)

	// Assert
	require.NotNil(t, dto)
	require.NotNil(t, dto.AbsoluteChange)
	assert.True(t, expectedAbsoluteChange.Equal(*dto.AbsoluteChange))
	assert.Nil(t, dto.PercentChange)
}

func TestPriceHistoryUpsertValidate(t *testing.T) {
	tests := []struct {
		name        string
		upsert      PriceHistoryUpsert
		expectedErr string
	}{
		{
			name: "valid value",
			upsert: PriceHistoryUpsert{
				Value: decimal.RequireFromString("10.50"),
			},
			expectedErr: "",
		},
		{
			name: "negative value",
			upsert: PriceHistoryUpsert{
				Value: decimal.RequireFromString("-1.00"),
			},
			expectedErr: "value must be zero or greater",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			upsert := tt.upsert

			// Act
			err := upsert.Validate()

			// Assert
			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
