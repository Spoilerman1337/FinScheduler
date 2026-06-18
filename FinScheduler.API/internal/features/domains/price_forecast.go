package domains

import (
	"finscheduler/pkg/ph"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PriceForecast struct {
	Id                  uuid.UUID       `db:"id"`
	ItemId              uuid.UUID       `db:"item_id"`
	CalculatedAt        time.Time       `db:"calculated_at"`
	LastKnownPrice      decimal.Decimal `db:"last_known_price"`
	AverageMonthlyDrift decimal.Decimal `db:"average_monthly_drift"`
}

type PriceForecastPointDto struct {
	Point          time.Time        `json:"point"`
	Value          decimal.Decimal  `json:"value"`
	AbsoluteChange *decimal.Decimal `json:"absoluteChange"`
	PercentChange  *decimal.Decimal `json:"percentChange"`
}

type PriceForecastUpsert struct {
	LastKnownPrice      decimal.Decimal `json:"lastKnownPrice"`
	AverageMonthlyDrift decimal.Decimal `json:"averageMonthlyDrift"`
}

func BuildPriceForecastPoints(priceForecast PriceForecast, monthsAhead int) []PriceForecastPointDto {
	if monthsAhead <= 0 {
		return make([]PriceForecastPointDto, 0)
	}

	points := make([]PriceForecastPointDto, 0, monthsAhead)
	monthlyDriftFactor := priceForecast.AverageMonthlyDrift.Div(decimal.NewFromInt(100))
	monthlyGrowthFactor := decimal.NewFromInt(1).Add(monthlyDriftFactor)
	currentValue := priceForecast.LastKnownPrice
	previousValue := priceForecast.LastKnownPrice
	for monthOffset := 1; monthOffset <= monthsAhead; monthOffset++ {
		currentValue = currentValue.Mul(monthlyGrowthFactor)
		value := currentValue.Round(2)
		absoluteChange, percentChange := ph.BuildPriceChange(value, &previousValue)
		points = append(points, PriceForecastPointDto{
			Point:          normalizeDate(priceForecast.CalculatedAt, monthOffset),
			Value:          value,
			AbsoluteChange: absoluteChange,
			PercentChange:  percentChange,
		})
		previousValue = value
	}

	return points
}

func normalizeDate(value time.Time, months int) time.Time {
	targetMonthStart := time.Date(
		value.Year(),
		value.Month(),
		1,
		value.Hour(),
		value.Minute(),
		value.Second(),
		value.Nanosecond(),
		value.Location(),
	).AddDate(0, months, 0)

	lastDayOfTargetMonth := time.Date(
		targetMonthStart.Year(),
		targetMonthStart.Month()+1,
		0,
		value.Hour(),
		value.Minute(),
		value.Second(),
		value.Nanosecond(),
		value.Location(),
	).Day()

	targetDay := value.Day()
	if targetDay > lastDayOfTargetMonth {
		targetDay = lastDayOfTargetMonth
	}

	return time.Date(
		targetMonthStart.Year(),
		targetMonthStart.Month(),
		targetDay,
		value.Hour(),
		value.Minute(),
		value.Second(),
		value.Nanosecond(),
		value.Location(),
	)
}

func (priceForecast *PriceForecastUpsert) Validate() error {
	if priceForecast.LastKnownPrice.IsNegative() {
		return fmt.Errorf("lastKnownPrice must be zero or greater")
	}

	return nil
}
