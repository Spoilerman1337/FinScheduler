package domains

import (
	"finscheduler/pkg/ph"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PriceHistory struct {
	Id         uuid.UUID       `db:"id"`
	ItemId     uuid.UUID       `db:"item_id"`
	RecordedAt time.Time       `db:"recorded_at"`
	Value      decimal.Decimal `db:"value"`
}

type PriceHistoryPointDto struct {
	Point          time.Time        `json:"point"`
	Value          decimal.Decimal  `json:"value"`
	AbsoluteChange *decimal.Decimal `json:"absoluteChange"`
	PercentChange  *decimal.Decimal `json:"percentChange"`
}

type PriceHistoryUpsert struct {
	Value decimal.Decimal `json:"value"`
}

func NewPriceHistoryPointDto(priceHistory PriceHistory, previousPriceHistory *PriceHistory) *PriceHistoryPointDto {
	dto := &PriceHistoryPointDto{
		Point: priceHistory.RecordedAt,
		Value: priceHistory.Value,
	}

	if previousPriceHistory == nil {
		return dto
	}

	dto.AbsoluteChange, dto.PercentChange = ph.BuildPriceChange(
		priceHistory.Value,
		&previousPriceHistory.Value,
	)

	return dto
}

func (priceHistory *PriceHistoryUpsert) Validate() error {
	if priceHistory.Value.IsNegative() {
		return fmt.Errorf("value must be zero or greater")
	}

	return nil
}
