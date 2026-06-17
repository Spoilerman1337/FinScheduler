package domains

import (
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

type PriceForecastUpsert struct {
	LastKnownPrice      decimal.Decimal `json:"lastKnownPrice"`
	AverageMonthlyDrift decimal.Decimal `json:"averageMonthlyDrift"`
}

func (priceForecast *PriceForecastUpsert) Validate() error {
	if priceForecast.LastKnownPrice.IsNegative() {
		return fmt.Errorf("lastKnownPrice must be zero or greater")
	}

	return nil
}
