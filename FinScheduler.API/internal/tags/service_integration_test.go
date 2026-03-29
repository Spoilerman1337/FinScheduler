//go:build integration
// +build integration

package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	repository := NewTagsRepository(testDB, testLogger)
	svc := NewTagsService(repository, testLogger)

	create := &TagCreate{
		Name: "Tag",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	tag, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Tag", *tag.Name)
}

func Test_TagsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	repository := NewTagsRepository(testDB, testLogger)
	svc := NewTagsService(repository, testLogger)

	create := &TagCreate{
		Name: "Ice",
	}

	update := &TagUpdate{
		Name: "Water",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	tag, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Water", *tag.Name)
}
