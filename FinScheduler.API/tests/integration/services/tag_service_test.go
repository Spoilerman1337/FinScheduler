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

	svc := services.NewTagsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

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

	svc := services.NewTagsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

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

func Test_TagsService_UpdateMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	svc := services.NewTagsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	update := &domains.TagUpdate{
		Name: "Missing",
	}

	// Act
	ok, err := svc.Update(testContext, uuid.New(), update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_TagsService_GetLookup_ShouldReturnOnlyActiveTagsOrderedByLowerName(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	svc := services.NewTagsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "banana", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Apple", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "carrot", IsActive: false})

	page := int32(0)
	pageSize := int32(10)

	// Act
	lookups, count, err := svc.GetLookup(testContext, &domains.TagLookupFilter{
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.Len(t, lookups, 2)
	assert.Equal(t, int64(2), count)
	assert.Equal(t, "Apple", *lookups[0].Label)
	assert.Equal(t, "banana", *lookups[1].Label)
}
