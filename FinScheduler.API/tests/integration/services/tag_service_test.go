//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewTagsService(uow, testLogger)

	create := &domains.TagCreate{
		Name: "Tag",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	tags, count, err := svc.Get(testContext, &domains.TagFilter{Ids: []*uuid.UUID{&id}})
	require.NoError(t, err)
	require.Len(t, tags, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Tag", *tags[0].Name)
}

func Test_TagsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	uow := persistence.NewUnitOfWork(testDB, testLogger)
	svc := services.NewTagsService(uow, testLogger)

	create := &domains.TagCreate{
		Name: "Ice",
	}

	update := &domains.TagUpdate{
		Name: "Water",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	tags, count, err := svc.Get(testContext, &domains.TagFilter{Ids: []*uuid.UUID{&id}})
	require.NoError(t, err)
	require.Len(t, tags, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Water", *tags[0].Name)
}
