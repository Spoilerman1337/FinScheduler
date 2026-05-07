package qh_test

import (
	"net/url"
	"testing"
	"time"

	"finscheduler/pkg/qh"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCategory string

const (
	testCategoryFood   testCategory = "Food"
	testCategoryTravel testCategory = "Travel"
)

func (category testCategory) IsValid() bool {
	switch category {
	case testCategoryFood, testCategoryTravel:
		return true
	default:
		return false
	}
}

func TestParseUUIDs_ShouldReturnNilWhenKeyIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "ids"

	// Act
	result, err := qh.ParseUUIDs(queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseUUIDs_ShouldParseAllValuesInOrder(t *testing.T) {
	// Arrange
	firstID := uuid.New()
	secondID := uuid.New()
	key := "ids"
	queryParams := url.Values{
		key: []string{firstID.String(), secondID.String()},
	}

	// Act
	result, err := qh.ParseUUIDs(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, firstID, *result[0])
	assert.Equal(t, secondID, *result[1])
}

func TestParseUUIDs_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "ids"
	invalidValue := "broken-uuid"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseUUIDs(queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), `invalid query parameter "ids" value "broken-uuid"`)
}

func TestParseString_ShouldReturnNilWhenValueIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "name"

	// Act
	result := qh.ParseString(queryParams, key)

	// Assert
	assert.Nil(t, result)
}

func TestParseString_ShouldReturnNilWhenValueIsEmpty(t *testing.T) {
	// Arrange
	key := "name"
	queryParams := url.Values{
		key: []string{""},
	}

	// Act
	result := qh.ParseString(queryParams, key)

	// Assert
	assert.Nil(t, result)
}

func TestParseString_ShouldReturnPointerWhenValueIsPresent(t *testing.T) {
	// Arrange
	key := "name"
	expectedValue := "coffee"
	queryParams := url.Values{
		key: []string{expectedValue},
	}

	// Act
	result := qh.ParseString(queryParams, key)

	// Assert
	require.NotNil(t, result)
	assert.Equal(t, expectedValue, *result)
}

func TestParseString_ShouldUseFirstValueWhenMultipleValuesArePresent(t *testing.T) {
	// Arrange
	key := "name"
	firstValue := "coffee"
	secondValue := "tea"
	queryParams := url.Values{
		key: []string{firstValue, secondValue},
	}

	// Act
	result := qh.ParseString(queryParams, key)

	// Assert
	require.NotNil(t, result)
	assert.Equal(t, firstValue, *result)
}

func TestParseDecimal_ShouldReturnNilWhenValueIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "price"

	// Act
	result, err := qh.ParseDecimal(queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseDecimal_ShouldReturnParsedDecimalWhenValueIsValid(t *testing.T) {
	// Arrange
	key := "price"
	expectedValue := decimal.RequireFromString("10.55")
	queryParams := url.Values{
		key: []string{"10.55"},
	}

	// Act
	result, err := qh.ParseDecimal(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, expectedValue.Equal(*result))
}

func TestParseDecimal_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "price"
	invalidValue := "bad-decimal"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseDecimal(queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), `invalid query parameter "price" value "bad-decimal"`)
}

func TestParseDecimal_ShouldUseFirstValueWhenMultipleValuesArePresent(t *testing.T) {
	// Arrange
	key := "price"
	firstValue := "10.55"
	secondValue := "99.99"
	expectedValue := decimal.RequireFromString(firstValue)
	queryParams := url.Values{
		key: []string{firstValue, secondValue},
	}

	// Act
	result, err := qh.ParseDecimal(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, expectedValue.Equal(*result))
}

func TestParseInt32_ShouldReturnNilWhenValueIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "page"

	// Act
	result, err := qh.ParseInt32(queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseInt32_ShouldReturnParsedInt32WhenValueIsValid(t *testing.T) {
	// Arrange
	key := "page"
	expectedValue := int32(42)
	queryParams := url.Values{
		key: []string{"42"},
	}

	// Act
	result, err := qh.ParseInt32(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedValue, *result)
}

func TestParseInt32_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "page"
	invalidValue := "bad-int"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseInt32(queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), `invalid query parameter "page" value "bad-int"`)
}

func TestParseInt32_ShouldUseFirstValueWhenMultipleValuesArePresent(t *testing.T) {
	// Arrange
	key := "page"
	firstValue := "42"
	secondValue := "99"
	expectedValue := int32(42)
	queryParams := url.Values{
		key: []string{firstValue, secondValue},
	}

	// Act
	result, err := qh.ParseInt32(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedValue, *result)
}

func TestParseBool_ShouldReturnNilWhenValueIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "isActive"

	// Act
	result, err := qh.ParseBool(queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseBool_ShouldReturnParsedBoolWhenValueIsValid(t *testing.T) {
	// Arrange
	key := "isActive"
	expectedValue := true
	queryParams := url.Values{
		key: []string{"true"},
	}

	// Act
	result, err := qh.ParseBool(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedValue, *result)
}

func TestParseBool_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "isActive"
	invalidValue := "not-bool"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseBool(queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), `invalid query parameter "isActive" value "not-bool"`)
}

func TestParseBool_ShouldUseFirstValueWhenMultipleValuesArePresent(t *testing.T) {
	// Arrange
	key := "isActive"
	firstValue := "true"
	secondValue := "false"
	expectedValue := true
	queryParams := url.Values{
		key: []string{firstValue, secondValue},
	}

	// Act
	result, err := qh.ParseBool(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, expectedValue, *result)
}

func TestParseTime_ShouldReturnNilWhenValueIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "createdAt"

	// Act
	result, err := qh.ParseTime(queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseTime_ShouldReturnParsedTimeWhenValueIsValid(t *testing.T) {
	// Arrange
	key := "createdAt"
	expectedValue := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	timeValue := expectedValue.Format(time.RFC3339)
	queryParams := url.Values{
		key: []string{timeValue},
	}

	// Act
	result, err := qh.ParseTime(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, expectedValue.Equal(*result))
}

func TestParseTime_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "createdAt"
	invalidValue := "not-time"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseTime(queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), `invalid query parameter "createdAt" value "not-time"`)
}

func TestParseTime_ShouldUseFirstValueWhenMultipleValuesArePresent(t *testing.T) {
	// Arrange
	key := "createdAt"
	expectedValue := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	firstValue := expectedValue.Format(time.RFC3339)
	secondValue := time.Date(2027, 2, 11, 13, 0, 0, 0, time.UTC).Format(time.RFC3339)
	queryParams := url.Values{
		key: []string{firstValue, secondValue},
	}

	// Act
	result, err := qh.ParseTime(queryParams, key)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, expectedValue.Equal(*result))
}

func TestParseEnums_ShouldReturnNilWhenKeyIsMissing(t *testing.T) {
	// Arrange
	queryParams := url.Values{}
	key := "categories"

	// Act
	result, err := qh.ParseEnums[testCategory](queryParams, key)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestParseEnums_ShouldReturnParsedEnumsInOrder(t *testing.T) {
	// Arrange
	key := "categories"
	firstCategory := string(testCategoryFood)
	secondCategory := string(testCategoryTravel)
	queryParams := url.Values{
		key: []string{firstCategory, secondCategory},
	}

	// Act
	result, err := qh.ParseEnums[testCategory](queryParams, key)

	// Assert
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, testCategoryFood, *result[0])
	assert.Equal(t, testCategoryTravel, *result[1])
}

func TestParseEnums_ShouldReturnErrorWhenValueIsInvalid(t *testing.T) {
	// Arrange
	key := "categories"
	invalidValue := "UnknownCategory"
	queryParams := url.Values{
		key: []string{invalidValue},
	}

	// Act
	result, err := qh.ParseEnums[testCategory](queryParams, key)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, `invalid query parameter "categories" value "UnknownCategory"`)
}
