package domains

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemFilter(t *testing.T) {
	// Arrange
	idOne := uuid.New()
	idTwo := uuid.New()
	tagOne := uuid.New()
	tagTwo := uuid.New()
	createdFrom := time.Date(2026, time.April, 20, 10, 0, 0, 0, time.UTC)
	createdTo := time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC)
	updatedFrom := time.Date(2026, time.April, 22, 10, 0, 0, 0, time.UTC)
	updatedTo := time.Date(2026, time.April, 23, 10, 0, 0, 0, time.UTC)

	req := httptest.NewRequest("GET", "/api/items?"+
		"ids="+idOne.String()+"&ids="+idTwo.String()+
		"&name=coffee"+
		"&priceFrom=10.50"+
		"&priceTo=15.75"+
		"&description=fresh"+
		"&isActive=true"+
		"&createdFrom="+createdFrom.Format(time.RFC3339)+
		"&createdTo="+createdTo.Format(time.RFC3339)+
		"&updatedFrom="+updatedFrom.Format(time.RFC3339)+
		"&updatedTo="+updatedTo.Format(time.RFC3339)+
		"&cashbackFrom=1"+
		"&cashbackTo=5"+
		"&categories="+string(FoodDrinks)+
		"&categories="+string(Travel)+
		"&tagIds="+tagOne.String()+
		"&tagIds="+tagTwo.String()+
		"&page=2"+
		"&pageSize=25", nil)

	// Act
	filter, err := NewItemFilter(req)

	// Assert
	require.NoError(t, err)
	require.Len(t, filter.Ids, 2)
	assert.Equal(t, idOne, *filter.Ids[0])
	assert.Equal(t, idTwo, *filter.Ids[1])
	require.NotNil(t, filter.Name)
	assert.Equal(t, "coffee", *filter.Name)
	require.NotNil(t, filter.PriceFrom)
	assert.True(t, decimal.NewFromFloat(10.50).Equal(*filter.PriceFrom))
	require.NotNil(t, filter.PriceTo)
	assert.True(t, decimal.NewFromFloat(15.75).Equal(*filter.PriceTo))
	require.NotNil(t, filter.Description)
	assert.Equal(t, "fresh", *filter.Description)
	require.NotNil(t, filter.IsActive)
	assert.True(t, *filter.IsActive)
	require.NotNil(t, filter.CreatedFrom)
	assert.True(t, createdFrom.Equal(*filter.CreatedFrom))
	require.NotNil(t, filter.CreatedTo)
	assert.True(t, createdTo.Equal(*filter.CreatedTo))
	require.NotNil(t, filter.UpdatedFrom)
	assert.True(t, updatedFrom.Equal(*filter.UpdatedFrom))
	require.NotNil(t, filter.UpdatedTo)
	assert.True(t, updatedTo.Equal(*filter.UpdatedTo))
	require.NotNil(t, filter.CashbackFrom)
	assert.Equal(t, int32(1), *filter.CashbackFrom)
	require.NotNil(t, filter.CashbackTo)
	assert.Equal(t, int32(5), *filter.CashbackTo)
	require.Len(t, filter.Categories, 2)
	assert.Equal(t, FoodDrinks, *filter.Categories[0])
	assert.Equal(t, Travel, *filter.Categories[1])
	require.Len(t, filter.TagIds, 2)
	assert.Equal(t, tagOne, *filter.TagIds[0])
	assert.Equal(t, tagTwo, *filter.TagIds[1])
	require.NotNil(t, filter.Page)
	assert.Equal(t, int32(2), *filter.Page)
	require.NotNil(t, filter.PageSize)
	assert.Equal(t, int32(25), *filter.PageSize)
}

func TestNewItemFilterReturnsErrorOnInvalidQuery(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/api/items?priceFrom=bad", nil)

	// Act
	_, err := NewItemFilter(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid query parameter")
}

func TestItemCreateValidate(t *testing.T) {
	t.Run("accepts valid item", func(t *testing.T) {
		// Arrange
		valid := &ItemCreate{
			Name:     "Coffee",
			Price:    decimal.NewFromFloat(10.50),
			Cashback: 3,
			Category: string(FoodDrinks),
			TagIds:   []string{uuid.New().String()},
		}

		// Act
		err := valid.Validate()

		// Assert
		require.NoError(t, err)
	})

	cases := []struct {
		name   string
		input  ItemCreate
		errMsg string
	}{
		{
			name: "short name",
			input: ItemCreate{
				Name:     "ab",
				Category: string(FoodDrinks),
			},
			errMsg: "name too short",
		},
		{
			name: "negative price",
			input: ItemCreate{
				Name:     "Coffee",
				Price:    decimal.NewFromInt(-1),
				Category: string(FoodDrinks),
			},
			errMsg: "price is negative",
		},
		{
			name: "negative cashback",
			input: ItemCreate{
				Name:     "Coffee",
				Cashback: -1,
				Category: string(FoodDrinks),
			},
			errMsg: "cashback is negative",
		},
		{
			name: "invalid category",
			input: ItemCreate{
				Name:     "Coffee",
				Category: "Invalid",
			},
			errMsg: "category is invalid",
		},
		{
			name: "invalid tag id",
			input: ItemCreate{
				Name:     "Coffee",
				Category: string(FoodDrinks),
				TagIds:   []string{"bad"},
			},
			errMsg: "tagId is invalid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			err := tc.input.Validate()

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestItemUpdateValidate(t *testing.T) {
	t.Run("accepts valid item", func(t *testing.T) {
		// Arrange
		valid := &ItemUpdate{
			Name:     "Coffee",
			Price:    decimal.NewFromFloat(10.50),
			Cashback: 3,
			Category: string(FoodDrinks),
			TagIds:   []string{uuid.New().String()},
		}

		// Act
		err := valid.Validate()

		// Assert
		require.NoError(t, err)
	})

	cases := []struct {
		name   string
		input  ItemUpdate
		errMsg string
	}{
		{
			name: "short name",
			input: ItemUpdate{
				Name:     "ab",
				Category: string(FoodDrinks),
			},
			errMsg: "name too short",
		},
		{
			name: "negative price",
			input: ItemUpdate{
				Name:     "Coffee",
				Price:    decimal.NewFromInt(-1),
				Category: string(FoodDrinks),
			},
			errMsg: "price is negative",
		},
		{
			name: "negative cashback",
			input: ItemUpdate{
				Name:     "Coffee",
				Cashback: -1,
				Category: string(FoodDrinks),
			},
			errMsg: "cashback is negative",
		},
		{
			name: "invalid category",
			input: ItemUpdate{
				Name:     "Coffee",
				Category: "Invalid",
			},
			errMsg: "category is invalid",
		},
		{
			name: "invalid tag id",
			input: ItemUpdate{
				Name:     "Coffee",
				Category: string(FoodDrinks),
				TagIds:   []string{"bad"},
			},
			errMsg: "tagId is invalid",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			err := tc.input.Validate()

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestItemFilterValidate(t *testing.T) {
	t.Run("accepts valid filter", func(t *testing.T) {
		// Arrange
		page := int32(0)
		pageSize := int32(20)
		priceFrom := decimal.NewFromFloat(10)
		priceTo := decimal.NewFromFloat(12)
		createdFrom := time.Date(2026, time.April, 20, 10, 0, 0, 0, time.UTC)
		createdTo := time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC)
		updatedFrom := time.Date(2026, time.April, 22, 10, 0, 0, 0, time.UTC)
		updatedTo := time.Date(2026, time.April, 23, 10, 0, 0, 0, time.UTC)
		cashbackFrom := int32(1)
		cashbackTo := int32(3)

		filter := &ItemFilter{
			Page:         &page,
			PageSize:     &pageSize,
			PriceFrom:    &priceFrom,
			PriceTo:      &priceTo,
			CreatedFrom:  &createdFrom,
			CreatedTo:    &createdTo,
			UpdatedFrom:  &updatedFrom,
			UpdatedTo:    &updatedTo,
			CashbackFrom: &cashbackFrom,
			CashbackTo:   &cashbackTo,
		}

		// Act
		err := filter.Validate()

		// Assert
		require.NoError(t, err)
	})

	cases := []struct {
		name   string
		build  func() ItemFilter
		errMsg string
	}{
		{
			name: "nil page",
			build: func() ItemFilter {
				pageSize := int32(20)
				return ItemFilter{PageSize: &pageSize}
			},
			errMsg: "page must be zero or greater",
		},
		{
			name: "negative page",
			build: func() ItemFilter {
				page := int32(-1)
				pageSize := int32(20)
				return ItemFilter{Page: &page, PageSize: &pageSize}
			},
			errMsg: "page must be zero or greater",
		},
		{
			name: "nil page size",
			build: func() ItemFilter {
				page := int32(0)
				return ItemFilter{Page: &page}
			},
			errMsg: "pageSize must be positive",
		},
		{
			name: "non-positive page size",
			build: func() ItemFilter {
				page := int32(0)
				pageSize := int32(0)
				return ItemFilter{Page: &page, PageSize: &pageSize}
			},
			errMsg: "pageSize must be positive",
		},
		{
			name: "invalid price range",
			build: func() ItemFilter {
				page := int32(0)
				pageSize := int32(20)
				priceFrom := decimal.NewFromFloat(20)
				priceTo := decimal.NewFromFloat(10)
				return ItemFilter{Page: &page, PageSize: &pageSize, PriceFrom: &priceFrom, PriceTo: &priceTo}
			},
			errMsg: "priceTo cannot be lesser than priceFrom",
		},
		{
			name: "invalid created range",
			build: func() ItemFilter {
				page := int32(0)
				pageSize := int32(20)
				createdFrom := time.Date(2026, time.April, 21, 10, 0, 0, 0, time.UTC)
				createdTo := time.Date(2026, time.April, 20, 10, 0, 0, 0, time.UTC)
				return ItemFilter{Page: &page, PageSize: &pageSize, CreatedFrom: &createdFrom, CreatedTo: &createdTo}
			},
			errMsg: "createTo cannot be earlier than createFrom",
		},
		{
			name: "invalid updated range",
			build: func() ItemFilter {
				page := int32(0)
				pageSize := int32(20)
				updatedFrom := time.Date(2026, time.April, 23, 10, 0, 0, 0, time.UTC)
				updatedTo := time.Date(2026, time.April, 22, 10, 0, 0, 0, time.UTC)
				return ItemFilter{Page: &page, PageSize: &pageSize, UpdatedFrom: &updatedFrom, UpdatedTo: &updatedTo}
			},
			errMsg: "updateTo cannot be earlier than updateFrom",
		},
		{
			name: "invalid cashback range",
			build: func() ItemFilter {
				page := int32(0)
				pageSize := int32(20)
				cashbackFrom := int32(3)
				cashbackTo := int32(1)
				return ItemFilter{Page: &page, PageSize: &pageSize, CashbackFrom: &cashbackFrom, CashbackTo: &cashbackTo}
			},
			errMsg: "cashbackFrom cannot be lesser than cashbackTO",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			filter := tc.build()

			// Act
			err := filter.Validate()

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestItemCategoryIsValidReturnsTrueForKnownCategory(t *testing.T) {
	// Arrange
	category := FoodDrinks

	// Act
	isValid := category.IsValid()

	// Assert
	assert.True(t, isValid)
}

func TestItemCategoryIsValidReturnsFalseForUnknownCategory(t *testing.T) {
	// Arrange
	category := ItemCategory("Invalid")

	// Act
	isValid := category.IsValid()

	// Assert
	assert.False(t, isValid)
}
