package domains

import "github.com/shopspring/decimal"

func buildPriceChange(
	currentValue decimal.Decimal,
	previousValue *decimal.Decimal,
) (*decimal.Decimal, *decimal.Decimal) {
	if previousValue == nil {
		return nil, nil
	}

	absoluteChange := currentValue.Sub(*previousValue)
	if previousValue.IsZero() {
		return &absoluteChange, nil
	}

	percentChange := absoluteChange.Div(*previousValue).Mul(decimal.NewFromInt(100))
	return &absoluteChange, &percentChange
}
