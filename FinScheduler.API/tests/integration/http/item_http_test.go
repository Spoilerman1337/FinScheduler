//go:build integration
// +build integration

package featurehttp_test

import (
	"encoding/json"
	"finscheduler/internal/features/domains"
	"finscheduler/tests/internal/testsupport"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ItemsHandler_GetListingInfo_ShouldReturnPaginatedItems(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	target := "/api/items?page=0&pageSize=20"
	expectedName := "Milk"
	expectedCount := int64(1)
	create := &domains.ItemCreate{
		Name:     expectedName,
		Price:    decimal.NewFromFloat(12.50),
		Category: "FoodDrinks",
	}

	_, createErr := app.itemsService.Create(ctx, create)
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse domains.PaginatedList[domains.ItemListingDto]
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Len(t, actualResponse.Data, 1)
	assert.Equal(t, expectedCount, actualResponse.Count)
	assert.Equal(t, expectedName, actualResponse.Data[0].Name)
}

func Test_ItemsHandler_GetDetailedInfo_ShouldReturnItem(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	expectedName := "Milk"
	priceHistoryDate := "2026-01-15"
	priceHistoryValue := decimal.RequireFromString("13.75")
	insertHistoryQuery := `INSERT INTO price_history (id, item_id, recorded_at, value) VALUES ($1, $2, $3, $4)`
	create := &domains.ItemCreate{
		Name:     expectedName,
		Price:    decimal.NewFromFloat(12.50),
		Category: "FoodDrinks",
	}

	itemID, createErr := app.itemsService.Create(ctx, create)
	_, insertHistoryErr := testDB.Exec(insertHistoryQuery, uuid.New(), itemID, priceHistoryDate, priceHistoryValue)
	target := "/api/items/" + itemID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse domains.ItemDetailedDto
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, insertHistoryErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, expectedName, actualResponse.Name)
	assert.Equal(t, 12.5, actualResponse.Price)
	assert.Equal(t, domains.ItemCategory("FoodDrinks"), actualResponse.Category)
	require.Len(t, actualResponse.PriceHistory, 1)
	assert.Equal(t, priceHistoryDate, actualResponse.PriceHistory[0].Point.UTC().Format("2006-01-02"))
	assert.True(t, priceHistoryValue.Equal(actualResponse.PriceHistory[0].Value))
}

func Test_ItemsHandler_GetDetailedInfo_ShouldReturnBadRequestOnInvalidID(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/items/bad-id"
	expectedBodyFragment := "invalid UUID length"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_GetDetailedInfo_ShouldReturnNotFoundForMissingItem(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	missingID := uuid.New()
	target := "/api/items/" + missingID.String()
	expectedBodyFragment := "item not found"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_GetListingInfo_ShouldReturnBadRequestOnInvalidQuery(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/items?page=bad&pageSize=20"
	expectedBodyFragment := `invalid query parameter "page"`
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_GetListingInfo_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodGet
	target := "/api/items?page=0&pageSize=20"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_ItemsHandler_Create_ShouldReturnCreatedWithLocationAndBody(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	method := http.MethodPost
	target := "/api/items"
	requestBody := `{"name":"Coffee","price":15.5,"category":"FoodDrinks"}`
	locationHeaderName := "Location"
	locationPrefix := "/api/items/"
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualID uuid.UUID
	decodeErr := json.NewDecoder(response.Body).Decode(&actualID)
	actualLocation := response.Header.Get(locationHeaderName)

	// Assert
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotEqual(t, uuid.Nil, actualID)
	assert.Equal(t, locationPrefix+actualID.String(), actualLocation)
}

func Test_ItemsHandler_Create_ShouldReturnBadRequestOnMalformedJSON(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPost
	target := "/api/items"
	requestBody := `{"name":`
	expectedBodyFragment := "unexpected EOF"
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_Create_ShouldReturnBadRequestOnInvalidReference(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	method := http.MethodPost
	target := "/api/items"
	invalidTagID := uuid.New()
	expectedBodyFragment := domains.ErrInvalidReference.Error()
	requestBody := `{"name":"Coffee","price":15.5,"category":"FoodDrinks","tagIds":["` + invalidTagID.String() + `"]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_Create_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodPost
	target := "/api/items"
	requestBody := `{"name":"Coffee","price":15.5,"category":"FoodDrinks"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_ItemsHandler_Update_ShouldReturnNoContent(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPut
	itemName := "Coffee"
	updatedName := "Latte"
	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	itemID, createErr := app.itemsService.Create(ctx, create)
	target := "/api/items/" + itemID.String()
	requestBody := `{"name":"` + updatedName + `","price":12.5,"category":"FoodDrinks"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	require.NoError(t, createErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func Test_ItemsHandler_Update_ShouldReturnBadRequestOnMalformedJSON(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	itemID := uuid.New()
	target := "/api/items/" + itemID.String()
	requestBody := `{"name":`
	expectedBodyFragment := "unexpected EOF"
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_Update_ShouldReturnBadRequestOnInvalidReference(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPut
	itemName := "Coffee"
	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	itemID, createErr := app.itemsService.Create(ctx, create)
	invalidTagID := uuid.New()
	target := "/api/items/" + itemID.String()
	expectedBodyFragment := domains.ErrInvalidReference.Error()
	requestBody := `{"name":"Coffee updated","price":12.5,"category":"FoodDrinks","tagIds":["` + invalidTagID.String() + `"]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	require.NoError(t, createErr)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_Update_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodPut
	itemID := uuid.New()
	target := "/api/items/" + itemID.String()
	requestBody := `{"name":"Coffee updated","price":12.5,"category":"FoodDrinks"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_ItemsHandler_Update_ShouldReturnNotFoundForMissingItem(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	missingID := uuid.New()
	target := "/api/items/" + missingID.String()
	expectedBodyFragment := "item not found"
	requestBody := `{"name":"Missing","price":12.5,"category":"FoodDrinks"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_UpdateCashbackByTag_ShouldReturnNoContentAndUpdateTaggedItems(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPatch
	tagCreate := &domains.TagCreate{Name: "Groceries"}

	tagID, tagCreateErr := app.tagsService.Create(ctx, tagCreate)
	firstTaggedID, firstTaggedCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Milk",
		Category: "FoodDrinks",
		TagIds:   []string{tagID.String()},
	})
	secondTaggedID, secondTaggedCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Bread",
		Category: "FoodDrinks",
		TagIds:   []string{tagID.String()},
	})
	untouchedID, untouchedCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Taxi",
		Category: "Transport",
	})
	target := "/api/items/cashback/tag"
	requestBody := `{"cashback":7,"tagId":"` + tagID.String() + `"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	firstTaggedItem, firstTaggedGetErr := app.itemsService.GetDetailedInfo(ctx, firstTaggedID)
	secondTaggedItem, secondTaggedGetErr := app.itemsService.GetDetailedInfo(ctx, secondTaggedID)
	untouchedItem, untouchedGetErr := app.itemsService.GetDetailedInfo(ctx, untouchedID)

	// Assert
	require.NoError(t, tagCreateErr)
	require.NoError(t, firstTaggedCreateErr)
	require.NoError(t, secondTaggedCreateErr)
	require.NoError(t, untouchedCreateErr)
	require.NoError(t, firstTaggedGetErr)
	require.NoError(t, secondTaggedGetErr)
	require.NoError(t, untouchedGetErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NotNil(t, firstTaggedItem)
	require.NotNil(t, secondTaggedItem)
	require.NotNil(t, untouchedItem)
	assert.Equal(t, int32(7), firstTaggedItem.Cashback)
	assert.Equal(t, int32(7), secondTaggedItem.Cashback)
	assert.Equal(t, int32(0), untouchedItem.Cashback)
}

func Test_ItemsHandler_UpdateCashbackByItems_ShouldReturnNoContentAndUpdateSelectedItems(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPatch

	firstItemID, firstCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Milk",
		Category: "FoodDrinks",
	})
	secondItemID, secondCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Taxi",
		Category: "Transport",
	})
	thirdItemID, thirdCreateErr := app.itemsService.Create(ctx, &domains.ItemCreate{
		Name:     "Gym",
		Category: "Sports",
	})
	target := "/api/items/cashback/items"
	requestBody := `{"cashback":11,"itemIds":["` + firstItemID.String() + `","` + thirdItemID.String() + `"]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	firstItem, firstGetErr := app.itemsService.GetDetailedInfo(ctx, firstItemID)
	secondItem, secondGetErr := app.itemsService.GetDetailedInfo(ctx, secondItemID)
	thirdItem, thirdGetErr := app.itemsService.GetDetailedInfo(ctx, thirdItemID)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, thirdCreateErr)
	require.NoError(t, firstGetErr)
	require.NoError(t, secondGetErr)
	require.NoError(t, thirdGetErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NotNil(t, firstItem)
	require.NotNil(t, secondItem)
	require.NotNil(t, thirdItem)
	assert.Equal(t, int32(11), firstItem.Cashback)
	assert.Equal(t, int32(0), secondItem.Cashback)
	assert.Equal(t, int32(11), thirdItem.Cashback)
}

func Test_ItemsHandler_UpdateCashbackByItems_ShouldReturnBadRequestOnInvalidPayload(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPatch
	target := "/api/items/cashback/items"
	expectedBodyFragment := "itemId is invalid: bad-uuid"
	requestBody := `{"cashback":11,"itemIds":["bad-uuid"]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_UpdateCashbackByTag_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodPatch
	target := "/api/items/cashback/tag"
	requestBody := `{"cashback":11,"tagId":"` + uuid.New().String() + `"}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_ItemsHandler_Delete_ShouldReturnNoContent(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodDelete
	itemName := "Coffee"
	create := &domains.ItemCreate{
		Name:     itemName,
		Price:    decimal.NewFromFloat(10.00),
		Category: "FoodDrinks",
	}

	itemID, createErr := app.itemsService.Create(ctx, create)
	target := "/api/items/" + itemID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	require.NoError(t, createErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
}

func Test_ItemsHandler_Delete_ShouldReturnBadRequestOnInvalidID(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodDelete
	target := "/api/items/bad-id"
	expectedBodyFragment := "invalid UUID length"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}

func Test_ItemsHandler_Delete_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodDelete
	itemID := uuid.New()
	target := "/api/items/" + itemID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_ItemsHandler_Delete_ShouldReturnNotFoundForMissingItem(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodDelete
	missingID := uuid.New()
	target := "/api/items/" + missingID.String()
	expectedBodyFragment := "item not found"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	actualBody := recorder.Body.String()

	// Assert
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Contains(t, actualBody, expectedBodyFragment)
}
