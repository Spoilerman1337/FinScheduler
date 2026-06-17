package domains

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestPriceForecastUpsertValidate(t *testing.T) {
	tests := []struct {
		name        string
		upsert      PriceForecastUpsert
		expectedErr string
	}{
		{
			name: "valid upsert",
			upsert: PriceForecastUpsert{
				CalculatedAt:        time.Date(2026, 6, 17, 10, 30, 0, 0, time.UTC),
				LastKnownPrice:      decimal.RequireFromString("123.45"),
				AverageMonthlyDrift: decimal.RequireFromString("-4.25"),
			},
		},
		{
			name: "missing calculatedAt",
			upsert: PriceForecastUpsert{
				LastKnownPrice:      decimal.RequireFromString("123.45"),
				AverageMonthlyDrift: decimal.RequireFromString("-4.25"),
			},
			expectedErr: "calculatedAt is required",
		},
		{
			name: "negative lastKnownPrice",
			upsert: PriceForecastUpsert{
				CalculatedAt:        time.Date(2026, 6, 17, 10, 30, 0, 0, time.UTC),
				LastKnownPrice:      decimal.RequireFromString("-1.00"),
				AverageMonthlyDrift: decimal.RequireFromString("-4.25"),
			},
			expectedErr: "lastKnownPrice must be zero or greater",
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
