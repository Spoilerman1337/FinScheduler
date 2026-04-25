//go:build integration
// +build integration

package httpapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finscheduler/internal/features/domains"
	"finscheduler/tests/internal/testsupport"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsHandler_Create_ShouldReturnCreatedAndLocation(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	router := newTestRouter()
	body := `{"name":"Milk","price":10.5,"description":"Fresh milk","isActive":true,"cashback":2,"category":"FoodDrinks","tagIds":[]}`

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var createdID string
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&createdID))
	parsedID, err := uuid.Parse(createdID)
	require.NoError(t, err)
	assert.Equal(t, "/api/items/"+parsedID.String(), rec.Header().Get("Location"))

	var count int
	require.NoError(t, testDB.Get(&count, "SELECT COUNT(*) FROM items WHERE id = $1", parsedID))
	assert.Equal(t, 1, count)
}

func Test_ItemsHandler_Create_ShouldReturnBadRequestForInvalidReference(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	router := newTestRouter()
	body := fmt.Sprintf(`{"name":"Milk","price":10.5,"description":"Fresh milk","isActive":true,"cashback":2,"category":"FoodDrinks","tagIds":["%s"]}`, uuid.New())

	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), domains.ErrInvalidReference.Error())
}

func Test_ItemsHandler_Get_ShouldReturnPaginatedItems(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	router := newTestRouter()
	tagID := testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Dairy", IsActive: true})
	itemID := testFixtures.MustCreateItem(t, &domains.ItemCreate{
		Name:        "Milk",
		Price:       decimal.NewFromFloat(10.5),
		Description: "Fresh milk",
		IsActive:    true,
		Cashback:    2,
		Category:    string(domains.FoodDrinks),
	})
	testFixtures.MustLinkItemTags(t, itemID, tagID)

	req := httptest.NewRequest(http.MethodGet, "/api/items?ids="+itemID.String()+"&page=0&pageSize=20", nil)
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var payload domains.PaginatedList[domains.ItemDto]
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	require.Len(t, payload.Data, 1)
	assert.Equal(t, int64(1), payload.Count)
	require.NotNil(t, payload.Data[0].Id)
	assert.Equal(t, itemID, *payload.Data[0].Id)
	assert.Equal(t, "Milk", *payload.Data[0].Name)
	assert.ElementsMatch(t, []string{"Dairy"}, lookupLabelsHTTP(payload.Data[0].Tags))
}

func Test_ItemsHandler_Update_ShouldReturnNotFoundForMissingItem(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	router := newTestRouter()
	body := `{"name":"Updated","price":11.5,"description":"Updated","isActive":true,"cashback":3,"category":"FoodDrinks","tagIds":[]}`

	req := httptest.NewRequest(http.MethodPut, "/api/items/"+uuid.New().String(), bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, strings.ToLower(rec.Body.String()), "item not found")
}

func lookupLabelsHTTP(values []*domains.Lookup) []string {
	res := make([]string, 0, len(values))
	for _, value := range values {
		if value != nil && value.Label != nil {
			res = append(res, *value.Label)
		}
	}
	return res
}
