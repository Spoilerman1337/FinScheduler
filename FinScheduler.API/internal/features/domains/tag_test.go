package domains

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTagsFilter(t *testing.T) {
	// Arrange
	idOne := uuid.New()
	idTwo := uuid.New()

	req := httptest.NewRequest("GET", "/api/tags?"+
		"ids="+idOne.String()+"&ids="+idTwo.String()+
		"&name=groceries"+
		"&isActive=true"+
		"&page=1"+
		"&pageSize=10", nil)

	// Act
	filter, err := NewTagsFilter(req)

	// Assert
	require.NoError(t, err)
	require.Len(t, filter.Ids, 2)
	assert.Equal(t, idOne, *filter.Ids[0])
	assert.Equal(t, idTwo, *filter.Ids[1])
	require.NotNil(t, filter.Name)
	assert.Equal(t, "groceries", *filter.Name)
	require.NotNil(t, filter.IsActive)
	assert.True(t, *filter.IsActive)
	require.NotNil(t, filter.Page)
	assert.Equal(t, int32(1), *filter.Page)
	require.NotNil(t, filter.PageSize)
	assert.Equal(t, int32(10), *filter.PageSize)
}

func TestNewTagsFilterReturnsErrorOnInvalidQuery(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/api/tags?isActive=not-bool", nil)

	// Act
	_, err := NewTagsFilter(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid query parameter")
}

func TestNewTagLookupFilter(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/api/tags/lookup?name=gro&page=0&pageSize=20", nil)

	// Act
	filter, err := NewTagLookupFilter(req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, filter.Name)
	assert.Equal(t, "gro", *filter.Name)
	require.NotNil(t, filter.Page)
	assert.Equal(t, int32(0), *filter.Page)
	require.NotNil(t, filter.PageSize)
	assert.Equal(t, int32(20), *filter.PageSize)
}

func TestTagCreateValidate(t *testing.T) {
	t.Run("accepts valid tag", func(t *testing.T) {
		// Arrange
		create := &TagCreate{Name: "Groceries"}

		// Act
		err := create.Validate()

		// Assert
		require.NoError(t, err)
	})

	t.Run("returns error for short name", func(t *testing.T) {
		// Arrange
		create := &TagCreate{Name: "ab"}

		// Act
		err := create.Validate()

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name too short")
	})
}

func TestTagUpdateValidate(t *testing.T) {
	t.Run("accepts valid tag", func(t *testing.T) {
		// Arrange
		update := &TagUpdate{Name: "Groceries"}

		// Act
		err := update.Validate()

		// Assert
		require.NoError(t, err)
	})

	t.Run("returns error for short name", func(t *testing.T) {
		// Arrange
		update := &TagUpdate{Name: "ab"}

		// Act
		err := update.Validate()

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name too short")
	})
}

func TestTagFilterValidate(t *testing.T) {
	t.Run("accepts valid filter", func(t *testing.T) {
		// Arrange
		page := int32(0)
		pageSize := int32(20)
		filter := &TagFilter{Page: &page, PageSize: &pageSize}

		// Act
		err := filter.Validate()

		// Assert
		require.NoError(t, err)
	})

	cases := []struct {
		name   string
		build  func() *TagFilter
		errMsg string
	}{
		{
			name: "nil page",
			build: func() *TagFilter {
				pageSize := int32(20)
				return &TagFilter{PageSize: &pageSize}
			},
			errMsg: "page must be zero or greater",
		},
		{
			name: "invalid page size",
			build: func() *TagFilter {
				page := int32(0)
				pageSize := int32(0)
				return &TagFilter{Page: &page, PageSize: &pageSize}
			},
			errMsg: "pageSize must be positive",
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

func TestTagLookupFilterValidate(t *testing.T) {
	t.Run("accepts valid filter", func(t *testing.T) {
		// Arrange
		page := int32(0)
		pageSize := int32(20)
		filter := &TagLookupFilter{Page: &page, PageSize: &pageSize}

		// Act
		err := filter.Validate()

		// Assert
		require.NoError(t, err)
	})

	cases := []struct {
		name   string
		build  func() *TagLookupFilter
		errMsg string
	}{
		{
			name: "invalid page",
			build: func() *TagLookupFilter {
				page := int32(-1)
				pageSize := int32(20)
				return &TagLookupFilter{Page: &page, PageSize: &pageSize}
			},
			errMsg: "page must be zero or greater",
		},
		{
			name: "invalid page size",
			build: func() *TagLookupFilter {
				page := int32(0)
				pageSize := int32(0)
				return &TagLookupFilter{Page: &page, PageSize: &pageSize}
			},
			errMsg: "pageSize must be positive",
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
