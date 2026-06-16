package domains

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPriceHistoryPointDto_ShouldMapOnlyGraphFields(t *testing.T) {
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
	dto := NewPriceHistoryPointDto(priceHistory)

	// Assert
	require.NotNil(t, dto)
	assert.Equal(t, recordedAt, dto.Point)
	assert.True(t, expectedValue.Equal(dto.Value))
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
