//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

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

	tag, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Tag", *tag.Name)
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

	tag, err := svc.GetById(testContext, id)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Water", *tag.Name)
}
