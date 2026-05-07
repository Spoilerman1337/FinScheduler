package domains

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaginatedList_ShouldReturnListWithProvidedDataAndCount(t *testing.T) {
	// Arrange
	data := []string{"first", "second"}
	count := int64(42)

	// Act
	result := NewPaginatedList(data, count)

	// Assert
	require.NotNil(t, result)
	assert.Equal(t, data, result.Data)
	assert.Equal(t, count, result.Count)
}
