package domains

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPriceForecastPoints(t *testing.T) {
	// Arrange
	calculatedAt := mustDate(t, "2026-06-17")
	priceForecast := PriceForecast{
		Id:                  uuid.New(),
		ItemId:              uuid.New(),
		CalculatedAt:        calculatedAt,
		LastKnownPrice:      decimal.RequireFromString("200.00"),
		AverageMonthlyDrift: decimal.RequireFromString("2.50"),
	}

	// Act
	points := BuildPriceForecastPoints(priceForecast, 3)

	// Assert
	require.Len(t, points, 3)
	assert.Equal(t, mustDate(t, "2026-07-17"), points[0].Point)
	assert.True(t, decimal.RequireFromString("205.00").Equal(points[0].Value))
	assert.Equal(t, mustDate(t, "2026-08-17"), points[1].Point)
	assert.True(t, decimal.RequireFromString("210.00").Equal(points[1].Value))
	assert.Equal(t, mustDate(t, "2026-09-17"), points[2].Point)
	assert.True(t, decimal.RequireFromString("215.00").Equal(points[2].Value))
}

func TestBuildPriceForecastPoints_ShouldClampToLastDayOfTargetMonth(t *testing.T) {
	tests := []struct {
		name           string
		calculatedAt   string
		monthsAhead    int
		expectedPoints []string
	}{
		{
			name:         "non leap year from january 31",
			calculatedAt: "2026-01-31",
			monthsAhead:  3,
			expectedPoints: []string{
				"2026-02-28",
				"2026-03-31",
				"2026-04-30",
			},
		},
		{
			name:         "non leap year from january 30",
			calculatedAt: "2026-01-30",
			monthsAhead:  3,
			expectedPoints: []string{
				"2026-02-28",
				"2026-03-30",
				"2026-04-30",
			},
		},
		{
			name:         "leap year from january 31",
			calculatedAt: "2024-01-31",
			monthsAhead:  3,
			expectedPoints: []string{
				"2024-02-29",
				"2024-03-31",
				"2024-04-30",
			},
		},
		{
			name:         "month with 31 days to month with 30 days",
			calculatedAt: "2026-05-31",
			monthsAhead:  3,
			expectedPoints: []string{
				"2026-06-30",
				"2026-07-31",
				"2026-08-31",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			priceForecast := PriceForecast{
				Id:                  uuid.New(),
				ItemId:              uuid.New(),
				CalculatedAt:        mustDate(t, tt.calculatedAt),
				LastKnownPrice:      decimal.RequireFromString("100.00"),
				AverageMonthlyDrift: decimal.RequireFromString("1.00"),
			}

			// Act
			points := BuildPriceForecastPoints(priceForecast, tt.monthsAhead)

			// Assert
			require.Len(t, points, len(tt.expectedPoints))
			for index, expectedPoint := range tt.expectedPoints {
				assert.Equal(t, mustDate(t, expectedPoint), points[index].Point)
			}
		})
	}
}

func TestPriceForecastUpsertValidate(t *testing.T) {
	tests := []struct {
		name        string
		upsert      PriceForecastUpsert
		expectedErr string
	}{
		{
			name: "valid upsert",
			upsert: PriceForecastUpsert{
				LastKnownPrice:      decimal.RequireFromString("123.45"),
				AverageMonthlyDrift: decimal.RequireFromString("-4.25"),
			},
		},
		{
			name: "negative lastKnownPrice",
			upsert: PriceForecastUpsert{
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

func mustDate(t *testing.T, value string) time.Time {
	t.Helper()

	result, err := time.Parse("2006-01-02", value)
	require.NoError(t, err)
	return result.UTC()
}
