package sh

import (
	"sort"

	"github.com/shopspring/decimal"
)

var decimalHalf = decimal.RequireFromString("0.5")
var winsorizationIQRMultiplier = decimal.RequireFromString("1.5")

func Winsorize(values []decimal.Decimal) []decimal.Decimal {
	result := make([]decimal.Decimal, len(values))
	copy(result, values)

	if len(result) < 2 {
		return result
	}

	sortedValues := make([]decimal.Decimal, len(values))
	copy(sortedValues, values)
	sort.Slice(sortedValues, func(leftIndex int, rightIndex int) bool {
		return sortedValues[leftIndex].LessThan(sortedValues[rightIndex])
	})

	lowerQuartile, upperQuartile := quartiles(sortedValues)
	interquartileRange := upperQuartile.Sub(lowerQuartile)
	lowerBound := lowerQuartile.Sub(interquartileRange.Mul(winsorizationIQRMultiplier))
	if lowerBound.IsNegative() {
		lowerBound = decimal.Zero
	}

	upperBound := upperQuartile.Add(interquartileRange.Mul(winsorizationIQRMultiplier))
	for index := range result {
		if result[index].Cmp(lowerBound) < 0 {
			result[index] = lowerBound
		} else if result[index].Cmp(upperBound) > 0 {
			result[index] = upperBound
		}
	}

	return result
}

func quartiles(sortedValues []decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	middleIndex := len(sortedValues) / 2
	lowerHalf := sortedValues[:middleIndex]
	upperHalf := sortedValues[middleIndex:]
	if len(sortedValues)%2 != 0 {
		upperHalf = sortedValues[middleIndex+1:]
	}

	return median(lowerHalf), median(upperHalf)
}

func median(values []decimal.Decimal) decimal.Decimal {
	middleIndex := len(values) / 2
	if len(values)%2 != 0 {
		return values[middleIndex]
	}

	return values[middleIndex-1].Add(values[middleIndex]).Mul(decimalHalf)
}
