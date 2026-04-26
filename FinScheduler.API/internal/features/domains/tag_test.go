package domains

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTagsFilter_ShouldParseAllSupportedFields(t *testing.T) {
	// Arrange
	firstID := uuid.New()
	secondID := uuid.New()
	requestURL := "/tags?ids=" + firstID.String() +
		"&ids=" + secondID.String() +
		"&name=groceries" +
		"&isActive=true" +
		"&page=2" +
		"&pageSize=25"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewTagsFilter(request)

	// Assert
	require.NoError(t, err)
	require.Len(t, filter.Ids, 2)
	require.NotNil(t, filter.Name)
	require.NotNil(t, filter.IsActive)
	require.NotNil(t, filter.Page)
	require.NotNil(t, filter.PageSize)

	assert.Equal(t, firstID, *filter.Ids[0])
	assert.Equal(t, secondID, *filter.Ids[1])
	assert.Equal(t, "groceries", *filter.Name)
	assert.True(t, *filter.IsActive)
	assert.Equal(t, int32(2), *filter.Page)
	assert.Equal(t, int32(25), *filter.PageSize)
}

func TestNewTagsFilter_ShouldReturnErrorOnInvalidQueryParam(t *testing.T) {
	// Arrange
	requestURL := "/tags?isActive=not-a-bool"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewTagsFilter(request)

	// Assert
	require.Error(t, err)
	assert.Equal(t, TagFilter{}, filter)
	assert.Contains(t, err.Error(), `invalid query parameter "isActive"`)
}

func TestNewTagLookupFilter_ShouldParseSupportedFields(t *testing.T) {
	// Arrange
	requestURL := "/tags/lookup?name=rec&page=1&pageSize=10"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewTagLookupFilter(request)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, filter.Name)
	require.NotNil(t, filter.Page)
	require.NotNil(t, filter.PageSize)

	assert.Equal(t, "rec", *filter.Name)
	assert.Equal(t, int32(1), *filter.Page)
	assert.Equal(t, int32(10), *filter.PageSize)
}

func TestNewTagLookupFilter_ShouldReturnErrorOnInvalidQueryParam(t *testing.T) {
	// Arrange
	requestURL := "/tags/lookup?page=wrong"
	request := httptest.NewRequest("GET", requestURL, nil)

	// Act
	filter, err := NewTagLookupFilter(request)

	// Assert
	require.Error(t, err)
	assert.Equal(t, TagLookupFilter{}, filter)
	assert.Contains(t, err.Error(), `invalid query parameter "page"`)
}

func TestNewTagDto_ShouldMapTagFields(t *testing.T) {
	// Arrange
	tagID := uuid.New()
	tag := Tag{
		Id:       tagID,
		Name:     "Recurring",
		IsActive: true,
	}

	// Act
	dto := NewTagDto(tag)

	// Assert
	require.NotNil(t, dto)
	require.NotNil(t, dto.Id)
	require.NotNil(t, dto.Name)
	require.NotNil(t, dto.IsActive)

	assert.Equal(t, tagID, *dto.Id)
	assert.Equal(t, "Recurring", *dto.Name)
	assert.True(t, *dto.IsActive)
}

func TestTagCreateValidate(t *testing.T) {
	tests := []struct {
		name          string
		tagCreate     TagCreate
		expectedError string
	}{
		{
			name: "valid tag create",
			tagCreate: TagCreate{
				Name:     "Groceries",
				IsActive: true,
			},
			expectedError: "",
		},
		{
			name: "name too short",
			tagCreate: TagCreate{
				Name:     "No",
				IsActive: true,
			},
			expectedError: "name too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tagCreate := tt.tagCreate

			// Act
			err := tagCreate.Validate()

			// Assert
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestTagUpdateValidate(t *testing.T) {
	tests := []struct {
		name          string
		tagUpdate     TagUpdate
		expectedError string
	}{
		{
			name: "valid tag update",
			tagUpdate: TagUpdate{
				Name:     "Transport",
				IsActive: false,
			},
			expectedError: "",
		},
		{
			name: "name too short",
			tagUpdate: TagUpdate{
				Name:     "No",
				IsActive: false,
			},
			expectedError: "name too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tagUpdate := tt.tagUpdate

			// Act
			err := tagUpdate.Validate()

			// Assert
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestTagFilterValidate(t *testing.T) {
	validPage := int32(0)
	validPageSize := int32(20)
	negativePage := int32(-1)
	zeroPageSize := int32(0)

	tests := []struct {
		name          string
		filter        TagFilter
		expectedError string
	}{
		{
			name: "valid filter",
			filter: TagFilter{
				Page:     &validPage,
				PageSize: &validPageSize,
			},
			expectedError: "",
		},
		{
			name: "page is nil",
			filter: TagFilter{
				Page:     nil,
				PageSize: &validPageSize,
			},
			expectedError: "page must be zero or greater",
		},
		{
			name: "page is negative",
			filter: TagFilter{
				Page:     &negativePage,
				PageSize: &validPageSize,
			},
			expectedError: "page must be zero or greater",
		},
		{
			name: "page size is nil",
			filter: TagFilter{
				Page:     &validPage,
				PageSize: nil,
			},
			expectedError: "pageSize must be positive",
		},
		{
			name: "page size is zero",
			filter: TagFilter{
				Page:     &validPage,
				PageSize: &zeroPageSize,
			},
			expectedError: "pageSize must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			filter := tt.filter

			// Act
			err := filter.Validate()

			// Assert
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

func TestTagLookupFilterValidate(t *testing.T) {
	validPage := int32(1)
	validPageSize := int32(15)
	negativePageSize := int32(-1)

	tests := []struct {
		name          string
		filter        TagLookupFilter
		expectedError string
	}{
		{
			name: "valid lookup filter",
			filter: TagLookupFilter{
				Page:     &validPage,
				PageSize: &validPageSize,
			},
			expectedError: "",
		},
		{
			name: "page is nil",
			filter: TagLookupFilter{
				Page:     nil,
				PageSize: &validPageSize,
			},
			expectedError: "page must be zero or greater",
		},
		{
			name: "page size is nil",
			filter: TagLookupFilter{
				Page:     &validPage,
				PageSize: nil,
			},
			expectedError: "pageSize must be positive",
		},
		{
			name: "page size is negative",
			filter: TagLookupFilter{
				Page:     &validPage,
				PageSize: &negativePageSize,
			},
			expectedError: "pageSize must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			filter := tt.filter

			// Act
			err := filter.Validate()

			// Assert
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
