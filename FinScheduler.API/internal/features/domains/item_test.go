package domains

import (
	"database/sql"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItemFilter_ShouldParseAllSupportedFields(t *testing.T) {
	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()
	tagID := uuid.New()
	createdFrom := "2026-01-01T10:00:00Z"
	createdTo := "2026-01-02T10:00:00Z"
	updatedFrom := "2026-01-03T10:00:00Z"
	updatedTo := "2026-01-04T10:00:00Z"
	requestURL := "/items?ids=" + id1.String() + "&ids=" + id2.String() +
		"&name=coffee" +
		"&priceFrom=10.50" +
		"&priceTo=20.75" +
		"&description=shop" +
		"&isActive=true" +
		"&createdFrom=" + createdFrom +
		"&createdTo=" + createdTo +
		"&updatedFrom=" + updatedFrom +
		"&updatedTo=" + updatedTo +
		"&cashbackFrom=5" +
		"&cashbackTo=10" +
		"&categories=" + string(FoodDrinks) +
		"&categories=" + string(Travel) +
		"&tagIds=" + tagID.String() +
		"&page=2" +
		"&pageSize=50"
	req := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewItemFilter(req)

	// Assert
	require.NoError(t, err)
	require.Len(t, filter.Ids, 2)
	require.NotNil(t, filter.Name)
	require.NotNil(t, filter.PriceFrom)
	require.NotNil(t, filter.PriceTo)
	require.NotNil(t, filter.Description)
	require.NotNil(t, filter.IsActive)
	require.NotNil(t, filter.CreatedFrom)
	require.NotNil(t, filter.CreatedTo)
	require.NotNil(t, filter.UpdatedFrom)
	require.NotNil(t, filter.UpdatedTo)
	require.NotNil(t, filter.CashbackFrom)
	require.NotNil(t, filter.CashbackTo)
	require.Len(t, filter.Categories, 2)
	require.Len(t, filter.TagIds, 1)
	require.NotNil(t, filter.Page)
	require.NotNil(t, filter.PageSize)

	assert.Equal(t, id1, *filter.Ids[0])
	assert.Equal(t, id2, *filter.Ids[1])
	assert.Equal(t, "coffee", *filter.Name)
	assert.True(t, decimal.RequireFromString("10.50").Equal(*filter.PriceFrom))
	assert.True(t, decimal.RequireFromString("20.75").Equal(*filter.PriceTo))
	assert.Equal(t, "shop", *filter.Description)
	assert.True(t, *filter.IsActive)
	assert.Equal(t, int32(5), *filter.CashbackFrom)
	assert.Equal(t, int32(10), *filter.CashbackTo)
	assert.Equal(t, FoodDrinks, *filter.Categories[0])
	assert.Equal(t, Travel, *filter.Categories[1])
	assert.Equal(t, tagID, *filter.TagIds[0])
	assert.Equal(t, int32(2), *filter.Page)
	assert.Equal(t, int32(50), *filter.PageSize)
	assert.Equal(t, createdFrom, filter.CreatedFrom.UTC().Format(time.RFC3339))
	assert.Equal(t, createdTo, filter.CreatedTo.UTC().Format(time.RFC3339))
	assert.Equal(t, updatedFrom, filter.UpdatedFrom.UTC().Format(time.RFC3339))
	assert.Equal(t, updatedTo, filter.UpdatedTo.UTC().Format(time.RFC3339))
}

func TestNewItemFilter_ShouldReturnErrorOnInvalidQueryParam(t *testing.T) {
	// Arrange
	requestURL := "/items?page=bad"
	req := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewItemFilter(req)

	// Assert
	require.Error(t, err)
	assert.Equal(t, ItemFilter{}, filter)
	assert.Contains(t, err.Error(), `invalid query parameter "page"`)
}

func TestNewItemDto_ShouldMapUpdatedAtAndTags(t *testing.T) {
	// Arrange
	itemID := uuid.New()
	tagID := uuid.New()
	createdAt := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 11, 13, 0, 0, 0, time.UTC)
	price := decimal.RequireFromString("99.95")

	item := Item{
		Id:          itemID,
		Name:        "Subscription",
		Price:       price,
		Description: "Monthly",
		IsActive:    true,
		CreatedAt:   createdAt,
		UpdatedAt:   sql.NullTime{Time: updatedAt, Valid: true},
		Cashback:    7,
		Category:    Subscriptions,
	}
	tags := []Tag{
		{
			Id:       tagID,
			Name:     "Recurring",
			IsActive: true,
		},
	}

	// Act
	dto := NewItemDto(item, tags)

	// Assert
	require.NotNil(t, dto)
	require.NotNil(t, dto.UpdatedAt)
	require.Len(t, dto.Tags, 1)
	require.NotNil(t, dto.Price)
	require.NotNil(t, dto.Category)

	assert.Equal(t, itemID, *dto.Id)
	assert.Equal(t, "Subscription", *dto.Name)
	assert.Equal(t, 99.95, *dto.Price)
	assert.Equal(t, "Monthly", *dto.Description)
	assert.True(t, *dto.IsActive)
	assert.Equal(t, createdAt, *dto.CreatedAt)
	assert.Equal(t, updatedAt, *dto.UpdatedAt)
	assert.Equal(t, int32(7), *dto.Cashback)
	assert.Equal(t, Subscriptions, *dto.Category)
	assert.Equal(t, "Recurring", *dto.Tags[0].Label)
	assert.Equal(t, tagID.String(), *dto.Tags[0].Value)
}

func TestNewItemDto_ShouldNilUpdatedAtAndTags(t *testing.T) {
	// Arrange
	item := Item{
		Id:          uuid.New(),
		Name:        "Gift",
		Price:       decimal.RequireFromString("15.00"),
		Description: "Box",
		IsActive:    false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   sql.NullTime{},
		Cashback:    0,
		Category:    Gifts,
	}

	// Act
	dto := NewItemDto(item, nil)

	// Assert
	require.NotNil(t, dto)
	assert.Nil(t, dto.UpdatedAt)
	assert.Nil(t, dto.Tags)
}

func TestItemCreateValidate(t *testing.T) {
	valid := ItemCreate{
		Name:        "Coffee",
		Price:       decimal.RequireFromString("10.50"),
		Description: "Latte",
		IsActive:    true,
		Cashback:    5,
		Category:    string(FoodDrinks),
		TagIds:      []string{uuid.New().String()},
	}

	tests := []struct {
		name        string
		mutate      func(item *ItemCreate)
		expectedErr string
	}{
		{
			name:        "valid",
			mutate:      func(item *ItemCreate) {},
			expectedErr: "",
		},
		{
			name: "name too short",
			mutate: func(item *ItemCreate) {
				item.Name = "No"
			},
			expectedErr: "name too short",
		},
		{
			name: "price is negative",
			mutate: func(item *ItemCreate) {
				item.Price = decimal.RequireFromString("-1")
			},
			expectedErr: "price is negative",
		},
		{
			name: "cashback is negative",
			mutate: func(item *ItemCreate) {
				item.Cashback = -1
			},
			expectedErr: "cashback is negative",
		},
		{
			name: "category is invalid",
			mutate: func(item *ItemCreate) {
				item.Category = "Unknown"
			},
			expectedErr: "category is invalid",
		},
		{
			name: "tag id is invalid",
			mutate: func(item *ItemCreate) {
				item.TagIds = []string{"bad-uuid"}
			},
			expectedErr: "tagId is invalid: bad-uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			item := valid
			item.TagIds = append([]string(nil), valid.TagIds...)
			tt.mutate(&item)

			// Act
			err := item.Validate()

			// Assert
			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}

func TestItemUpdateValidate(t *testing.T) {
	valid := ItemUpdate{
		Name:        "Coffee",
		Price:       decimal.RequireFromString("10.50"),
		Description: "Latte",
		IsActive:    true,
		Cashback:    5,
		Category:    string(FoodDrinks),
		TagIds:      []string{uuid.New().String()},
	}

	tests := []struct {
		name        string
		mutate      func(item *ItemUpdate)
		expectedErr string
	}{
		{
			name:        "valid",
			mutate:      func(item *ItemUpdate) {},
			expectedErr: "",
		},
		{
			name: "name too short",
			mutate: func(item *ItemUpdate) {
				item.Name = "No"
			},
			expectedErr: "name too short",
		},
		{
			name: "price is negative",
			mutate: func(item *ItemUpdate) {
				item.Price = decimal.RequireFromString("-1")
			},
			expectedErr: "price is negative",
		},
		{
			name: "cashback is negative",
			mutate: func(item *ItemUpdate) {
				item.Cashback = -1
			},
			expectedErr: "cashback is negative",
		},
		{
			name: "category is invalid",
			mutate: func(item *ItemUpdate) {
				item.Category = "Unknown"
			},
			expectedErr: "category is invalid",
		},
		{
			name: "tag id is invalid",
			mutate: func(item *ItemUpdate) {
				item.TagIds = []string{"bad-uuid"}
			},
			expectedErr: "tagId is invalid: bad-uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			item := valid
			item.TagIds = append([]string(nil), valid.TagIds...)
			tt.mutate(&item)

			// Act
			err := item.Validate()

			// Assert
			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}

func TestItemFilterValidate(t *testing.T) {
	page := int32(0)
	pageSize := int32(20)
	priceFrom := decimal.RequireFromString("10")
	priceTo := decimal.RequireFromString("20")
	createdFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	createdTo := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	updatedFrom := time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)
	updatedTo := time.Date(2026, 1, 4, 0, 0, 0, 0, time.UTC)
	cashbackFrom := int32(1)
	cashbackTo := int32(5)

	valid := ItemFilter{
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

	tests := []struct {
		name        string
		mutate      func(filter *ItemFilter)
		expectedErr string
	}{
		{
			name:        "valid",
			mutate:      func(filter *ItemFilter) {},
			expectedErr: "",
		},
		{
			name: "page is nil",
			mutate: func(filter *ItemFilter) {
				filter.Page = nil
			},
			expectedErr: "page must be zero or greater",
		},
		{
			name: "page size is nil",
			mutate: func(filter *ItemFilter) {
				filter.PageSize = nil
			},
			expectedErr: "pageSize must be positive",
		},
		{
			name: "price range is reversed",
			mutate: func(filter *ItemFilter) {
				reversed := decimal.RequireFromString("5")
				filter.PriceTo = &reversed
			},
			expectedErr: "priceTo cannot be lesser than priceFrom",
		},
		{
			name: "created range is reversed",
			mutate: func(filter *ItemFilter) {
				reversed := createdFrom.Add(-time.Hour)
				filter.CreatedTo = &reversed
			},
			expectedErr: "createTo cannot be earlier than createFrom",
		},
		{
			name: "updated range is reversed",
			mutate: func(filter *ItemFilter) {
				reversed := updatedFrom.Add(-time.Hour)
				filter.UpdatedTo = &reversed
			},
			expectedErr: "updateTo cannot be earlier than updateFrom",
		},
		{
			name: "cashback range is reversed",
			mutate: func(filter *ItemFilter) {
				reversed := int32(0)
				filter.CashbackTo = &reversed
			},
			expectedErr: "cashbackFrom cannot be lesser than cashbackTO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			filter := valid
			tt.mutate(&filter)

			// Act
			err := filter.Validate()

			// Assert
			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}

func TestItemCategoryIsValid(t *testing.T) {
	// Arrange
	validCategoryOne := FoodDrinks
	validCategoryTwo := Education
	invalidCategoryOne := ItemCategory("Unknown")
	invalidCategoryTwo := ItemCategory("")

	// Act
	isValidOne := validCategoryOne.IsValid()
	isValidTwo := validCategoryTwo.IsValid()
	isInvalidOne := invalidCategoryOne.IsValid()
	isInvalidTwo := invalidCategoryTwo.IsValid()

	// Assert
	assert.True(t, isValidOne)
	assert.True(t, isValidTwo)
	assert.False(t, isInvalidOne)
	assert.False(t, isInvalidTwo)
}

func TestValidateTagIds(t *testing.T) {
	t.Run("valid ids", func(t *testing.T) {
		// Arrange
		tagIDs := []string{uuid.New().String(), uuid.New().String()}

		// Act
		err := validateTagIds(tagIDs)

		// Assert
		require.NoError(t, err)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange
		tagIDs := []string{uuid.New().String(), "broken"}
		expectedError := "tagId is invalid: broken"

		// Act
		err := validateTagIds(tagIDs)

		// Assert
		require.EqualError(t, err, expectedError)
	})
}
