package qh

import (
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEnum string

const (
	testEnumFirst  testEnum = "First"
	testEnumSecond testEnum = "Second"
)

func (e testEnum) IsValid() bool {
	return e == testEnumFirst || e == testEnumSecond
}

func TestParseUUIDs(t *testing.T) {
	t.Run("returns nil when key is absent", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseUUIDs(url.Values{}, "ids")

		// Assert
		require.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("parses repeated UUID values", func(t *testing.T) {
		// Arrange
		first := uuid.New()
		second := uuid.New()

		// Act
		res, err := ParseUUIDs(url.Values{
			"ids": []string{first.String(), second.String()},
		}, "ids")

		// Assert
		require.NoError(t, err)
		require.Len(t, res, 2)
		assert.Equal(t, first, *res[0])
		assert.Equal(t, second, *res[1])
	})

	t.Run("returns an error on invalid UUID", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseUUIDs(url.Values{
			"ids": []string{"not-a-uuid"},
		}, "ids")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid query parameter")
	})
}

func TestParseString(t *testing.T) {
	t.Run("returns nil when key is absent", func(t *testing.T) {
		// Arrange

		// Act
		res := ParseString(url.Values{}, "name")

		// Assert
		assert.Nil(t, res)
	})

	t.Run("parses present string value", func(t *testing.T) {
		// Arrange

		// Act
		res := ParseString(url.Values{"name": []string{"coffee"}}, "name")

		// Assert
		require.NotNil(t, res)
		assert.Equal(t, "coffee", *res)
	})
}

func TestParseDecimal(t *testing.T) {
	t.Run("parses valid decimal", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseDecimal(url.Values{"price": []string{"10.50"}}, "price")

		// Assert
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.True(t, decimal.NewFromFloat(10.50).Equal(*res))
	})

	t.Run("returns error for invalid decimal", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseDecimal(url.Values{"price": []string{"bad"}}, "price")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestParseInt32(t *testing.T) {
	t.Run("parses valid int32", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseInt32(url.Values{"page": []string{"42"}}, "page")

		// Assert
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, int32(42), *res)
	})

	t.Run("returns error for invalid int32", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseInt32(url.Values{"page": []string{"bad"}}, "page")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestParseBool(t *testing.T) {
	t.Run("parses valid bool", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseBool(url.Values{"isActive": []string{"true"}}, "isActive")

		// Assert
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.True(t, *res)
	})

	t.Run("returns error for invalid bool", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseBool(url.Values{"isActive": []string{"bad"}}, "isActive")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestParseTime(t *testing.T) {
	t.Run("parses valid time", func(t *testing.T) {
		// Arrange
		expected := time.Date(2026, time.April, 25, 12, 0, 0, 0, time.UTC)

		// Act
		res, err := ParseTime(url.Values{"createdFrom": []string{expected.Format(time.RFC3339)}}, "createdFrom")

		// Assert
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.True(t, expected.Equal(*res))
	})

	t.Run("returns error for invalid time", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseTime(url.Values{"createdFrom": []string{"2026-04-25"}}, "createdFrom")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestParseEnums(t *testing.T) {
	t.Run("returns nil when key is absent", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseEnums[testEnum](url.Values{}, "categories")

		// Assert
		require.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("parses repeated enum values", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseEnums[testEnum](url.Values{
			"categories": []string{string(testEnumFirst), string(testEnumSecond)},
		}, "categories")

		// Assert
		require.NoError(t, err)
		require.Len(t, res, 2)
		assert.Equal(t, testEnumFirst, *res[0])
		assert.Equal(t, testEnumSecond, *res[1])
	})

	t.Run("returns an error on invalid enum value", func(t *testing.T) {
		// Arrange

		// Act
		res, err := ParseEnums[testEnum](url.Values{
			"categories": []string{"Invalid"},
		}, "categories")

		// Assert
		require.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid query parameter")
	})
}
