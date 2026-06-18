package sh

import "github.com/shopspring/decimal"

var decimalPercentDivisor = decimal.NewFromInt(100)
var decimalOne = decimal.NewFromInt(1)

func WinsorizeValue(
	value decimal.Decimal,
	projectedValue decimal.Decimal,
	windowPercent decimal.Decimal,
) decimal.Decimal {
	if projectedValue.Cmp(decimal.Zero) <= 0 || windowPercent.IsNegative() {
		return value
	}

	windowFactor := windowPercent.Div(decimalPercentDivisor)
	lowerBound := projectedValue.Mul(decimalOne.Sub(windowFactor))
	upperBound := projectedValue.Mul(decimalOne.Add(windowFactor))

	if value.Cmp(lowerBound) < 0 {
		return lowerBound
	}

	if value.Cmp(upperBound) > 0 {
		return upperBound
	}

	return value
}
