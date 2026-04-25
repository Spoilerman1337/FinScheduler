//go:build integration
// +build integration

package services_test

import (
	"errors"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsService_Flow_CreateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	create := &domains.ItemCreate{
		Name:     "Item",
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	page := int32(0)
	pageSize := int32(20)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, err := svc.Get(testContext, &filter)
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Item", *items[0].Name)
}

func Test_ItemsService_UpdateAndGet_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	create := &domains.ItemCreate{
		Name:     "Ice",
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	update := &domains.ItemUpdate{
		Name:     "Water",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Update(testContext, id, update)
	require.NoError(t, err)
	require.True(t, ok)

	page := int32(0)
	pageSize := int32(20)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, err := svc.Get(testContext, &filter)
	require.NoError(t, err)
	require.Len(t, items, 1)

	// Assert
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Water", *items[0].Name)
	assert.Equal(t, 15.50, *items[0].Price)
}

func Test_ItemsService_DeleteAndGet_ShouldErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	create := &domains.ItemCreate{
		Name:     "Orange",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
	}

	// Act
	id, err := svc.Create(testContext, create)
	require.NoError(t, err)

	ok, err := svc.Delete(testContext, id)
	require.NoError(t, err)
	require.True(t, ok)

	page := int32(0)
	pageSize := int32(20)
	filter := domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	}

	items, count, err := svc.Get(testContext, &filter)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Equal(t, int64(0), count)
}

func Test_ItemsService_UpdateMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	update := &domains.ItemUpdate{
		Name:     "Missing",
		Category: "FoodDrinks",
	}

	// Act
	ok, err := svc.Update(testContext, uuid.New(), update)

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsService_DeleteMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	// Act
	ok, err := svc.Delete(testContext, uuid.New())

	// Assert
	require.NoError(t, err)
	assert.False(t, ok)
}

func Test_ItemsService_Create_ShouldRollbackItemWhenTagInsertFails(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)

	create := &domains.ItemCreate{
		Name:     "Rollback",
		Price:    decimal.NewFromFloat(15.50),
		Category: "FoodDrinks",
		TagIds:   []string{uuid.New().String()},
	}

	// Act
	id, err := svc.Create(testContext, create)

	// Assert
	require.Error(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE name = $1", create.Name))
	assert.Equal(t, 0, count)
}

func Test_ItemsService_Create_ShouldAttachTagsAndPopulateGet(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)
	tagOne := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food", IsActive: true})
	tagTwo := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Morning", IsActive: true})

	// Act
	id, err := svc.Create(testContext, &domains.ItemCreate{
		Name:        "Latte",
		Price:       decimal.NewFromFloat(8.5),
		Description: "Coffee with milk",
		IsActive:    true,
		Cashback:    2,
		Category:    string(domains.FoodDrinks),
		TagIds:      []string{tagOne.String(), tagTwo.String()},
	})
	page := int32(0)
	pageSize := int32(20)
	items, count, getErr := svc.Get(testContext, &domains.ItemFilter{
		Ids:      []*uuid.UUID{&id},
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, id)
	require.NoError(t, getErr)
	require.Len(t, items, 1)
	assert.Equal(t, int64(1), count)
	assert.ElementsMatch(t, []string{"Food", "Morning"}, lookupLabels(items[0].Tags))
	assert.ElementsMatch(t, []string{tagOne.String(), tagTwo.String()}, lookupValues(items[0].Tags))
}

func Test_ItemsService_Create_ShouldMapForeignKeyViolationToInvalidReference(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)
	missingTagID := uuid.New()

	// Act
	id, err := svc.Create(testContext, &domains.ItemCreate{
		Name:     "Invalid References",
		Category: "FoodDrinks",
		TagIds:   []string{missingTagID.String()},
	})

	// Assert
	require.Error(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	assert.True(t, errors.Is(err, domains.ErrInvalidReference))

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE name = $1", "Invalid References"))
	assert.Equal(t, 0, count)
}

func Test_ItemsService_Update_ShouldReconcileTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)
	linkRepo := repositories.NewTagToItemsRepository(testDB, testLogger)

	tagOne := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Food", IsActive: true})
	tagTwo := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Coffee", IsActive: true})
	tagThree := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Office", IsActive: true})

	itemID, err := svc.Create(testContext, &domains.ItemCreate{
		Name:     "Espresso",
		Category: "FoodDrinks",
		TagIds:   []string{tagOne.String(), tagTwo.String()},
	})
	require.NoError(t, err)

	// Act
	ok, err := svc.Update(testContext, itemID, &domains.ItemUpdate{
		Name:        "Espresso Updated",
		Price:       decimal.NewFromFloat(9.25),
		Description: "Updated description",
		IsActive:    true,
		Cashback:    4,
		Category:    "FoodDrinks",
		TagIds:      []string{tagTwo.String(), tagThree.String()},
	})
	relations, relationErr := linkRepo.GetByItemIds(testContext, []uuid.UUID{itemID})
	page := int32(0)
	pageSize := int32(20)
	items, count, getErr := svc.Get(testContext, &domains.ItemFilter{
		Ids:      []*uuid.UUID{&itemID},
		Page:     &page,
		PageSize: &pageSize,
	})

	// Assert
	require.NoError(t, err)
	assert.True(t, ok)
	require.NoError(t, relationErr)
	assert.ElementsMatch(t, []uuid.UUID{tagTwo, tagThree}, tagIDs(relations))
	require.NoError(t, getErr)
	require.Len(t, items, 1)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, "Espresso Updated", *items[0].Name)
	assert.ElementsMatch(t, []string{"Coffee", "Office"}, lookupLabels(items[0].Tags))
}

func Test_ItemsService_Delete_ShouldCascadeRelations(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	svc := services.NewItemsService(persistence.NewUnitOfWork(testDB, testLogger), testLogger)
	linkRepo := repositories.NewTagToItemsRepository(testDB, testLogger)

	tagID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Attached", IsActive: true})
	itemID, err := svc.Create(testContext, &domains.ItemCreate{
		Name:     "Delete Cascade",
		Category: "FoodDrinks",
		TagIds:   []string{tagID.String()},
	})
	require.NoError(t, err)

	// Act
	ok, err := svc.Delete(testContext, itemID)
	relations, relationErr := linkRepo.GetByItemIds(testContext, []uuid.UUID{itemID})

	// Assert
	require.NoError(t, err)
	assert.True(t, ok)
	require.NoError(t, relationErr)
	assert.Empty(t, relations)
}

func lookupLabels(values []*domains.Lookup) []string {
	res := make([]string, 0, len(values))
	for _, value := range values {
		if value != nil && value.Label != nil {
			res = append(res, *value.Label)
		}
	}
	return res
}

func lookupValues(values []*domains.Lookup) []string {
	res := make([]string, 0, len(values))
	for _, value := range values {
		if value != nil && value.Value != nil {
			res = append(res, *value.Value)
		}
	}
	return res
}

func tagIDs(relations []domains.TagToItem) []uuid.UUID {
	res := make([]uuid.UUID, 0, len(relations))
	for _, relation := range relations {
		res = append(res, relation.TagId)
	}
	return res
}
