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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsHandler_GetListingInfo_ShouldReturnPaginatedTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	target := "/api/tags?page=0&pageSize=20"
	expectedName := "Groceries"
	expectedCount := int64(1)
	create := &domains.TagCreate{Name: expectedName}

	_, createErr := app.tagsService.Create(ctx, create)
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse domains.PaginatedList[domains.TagListingDto]
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Len(t, actualResponse.Data, 1)
	assert.Equal(t, expectedCount, actualResponse.Count)
	assert.Equal(t, expectedName, *actualResponse.Data[0].Name)
}

func Test_TagsHandler_GetListingInfo_ShouldReturnBadRequestOnInvalidQuery(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/tags?page=bad&pageSize=20"
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

func Test_TagsHandler_GetListingInfo_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodGet
	target := "/api/tags?page=0&pageSize=20"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_TagsHandler_GetDetailedInfo_ShouldReturnTag(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	expectedName := "Groceries"
	expectedIsActive := true
	create := &domains.TagCreate{Name: expectedName, IsActive: expectedIsActive}

	tagID, createErr := app.tagsService.Create(ctx, create)
	target := "/api/tags/" + tagID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse domains.TagDetailedDto
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, expectedName, *actualResponse.Name)
	assert.Equal(t, expectedIsActive, *actualResponse.IsActive)
}

func Test_TagsHandler_GetDetailedInfo_ShouldReturnBadRequestOnInvalidID(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/tags/bad-id"
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

func Test_TagsHandler_GetDetailedInfo_ShouldReturnNotFoundForMissingTag(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	missingID := uuid.New()
	target := "/api/tags/" + missingID.String()
	expectedBodyFragment := "tag not found"
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

func Test_TagsHandler_GetDetailedInfo_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodGet
	tagID := uuid.New()
	target := "/api/tags/" + tagID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_TagsHandler_GetLookup_ShouldReturnOnlyActiveTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	target := "/api/tags/lookup?page=0&pageSize=20"
	activeName := "Groceries"
	inactiveName := "Archived"
	activeCreate := &domains.TagCreate{Name: activeName, IsActive: true}
	inactiveCreate := &domains.TagCreate{Name: inactiveName, IsActive: false}

	_, activeCreateErr := app.tagsService.Create(ctx, activeCreate)
	_, inactiveCreateErr := app.tagsService.Create(ctx, inactiveCreate)
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse domains.PaginatedList[domains.Lookup]
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, activeCreateErr)
	require.NoError(t, inactiveCreateErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Len(t, actualResponse.Data, 1)
	assert.Equal(t, int64(1), actualResponse.Count)
	assert.Equal(t, activeName, *actualResponse.Data[0].Label)
}

func Test_TagsHandler_GetLookup_ShouldReturnBadRequestOnInvalidQuery(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/tags/lookup?page=bad&pageSize=20"
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

func Test_TagsHandler_GetLookup_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodGet
	target := "/api/tags/lookup?page=0&pageSize=20"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_TagsHandler_Create_ShouldReturnCreatedWithLocationAndBody(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	method := http.MethodPost
	target := "/api/tags"
	requestBody := `{"name":"Groceries","isActive":true}`
	locationHeaderName := "Location"
	locationPrefix := "/api/tags/"
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

func Test_TagsHandler_Create_ShouldReturnBadRequestOnMalformedJSON(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPost
	target := "/api/tags"
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

func Test_TagsHandler_Create_ShouldReturnBadRequestOnValidationError(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPost
	target := "/api/tags"
	requestBody := `{"name":"No","isActive":true}`
	expectedBodyFragment := "name too short"
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

func Test_TagsHandler_Create_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodPost
	target := "/api/tags"
	requestBody := `{"name":"Groceries","isActive":true}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_TagsHandler_Update_ShouldReturnNoContent(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB)
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPut
	originalName := "Groceries"
	updatedName := "Supermarket"
	create := &domains.TagCreate{Name: originalName}

	tagID, createErr := app.tagsService.Create(ctx, create)
	target := "/api/tags/" + tagID.String()
	requestBody := `{"name":"` + updatedName + `","isActive":true}`
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

func Test_TagsHandler_Update_ShouldReturnBadRequestOnMalformedJSON(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	tagID := uuid.New()
	target := "/api/tags/" + tagID.String()
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

func Test_TagsHandler_Update_ShouldReturnBadRequestOnInvalidID(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	target := "/api/tags/bad-id"
	requestBody := `{"name":"Supermarket","isActive":true}`
	expectedBodyFragment := "invalid UUID length"
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

func Test_TagsHandler_Update_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodPut
	tagID := uuid.New()
	target := "/api/tags/" + tagID.String()
	requestBody := `{"name":"Supermarket","isActive":true}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_TagsHandler_Update_ShouldReturnNotFoundForMissingTag(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	missingID := uuid.New()
	target := "/api/tags/" + missingID.String()
	requestBody := `{"name":"Missing","isActive":true}`
	expectedBodyFragment := "tag not found"
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
