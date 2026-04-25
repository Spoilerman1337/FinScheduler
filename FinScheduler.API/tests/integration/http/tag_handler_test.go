//go:build integration
// +build integration

package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finscheduler/internal/features/domains"
	"finscheduler/tests/internal/testsupport"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TagsHandler_GetLookup_ShouldReturnOnlyActiveTags(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	router := newTestRouter()
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "banana", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "Apple", IsActive: true})
	testFixtures.MustCreateTag(t, &domains.TagCreate{Name: "carrot", IsActive: false})

	req := httptest.NewRequest(http.MethodGet, "/api/tags/lookup?page=0&pageSize=20", nil)
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

	var payload domains.PaginatedList[domains.Lookup]
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&payload))
	require.Len(t, payload.Data, 2)
	assert.Equal(t, int64(2), payload.Count)
	assert.Equal(t, "Apple", *payload.Data[0].Label)
	assert.Equal(t, "banana", *payload.Data[1].Label)
}

func Test_TagsHandler_Update_ShouldReturnNotFoundForMissingTag(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "tags")
	})

	router := newTestRouter()
	body := `{"name":"Updated","isActive":true}`

	req := httptest.NewRequest(http.MethodPut, "/api/tags/"+uuid.New().String(), bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, strings.ToLower(rec.Body.String()), "tag not found")
}
