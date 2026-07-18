//go:build integration
// +build integration

package featurehttp_test

import (
	"encoding/json"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CalendarEventsHandler_GetByDateRange_ShouldReturnCalendarEvents(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodGet
	target := "/api/calendar-events?from=2026-07-17&to=2026-07-18"
	calendarEventsService := newCalendarEventsService(testDB)

	_, firstCreateErr := calendarEventsService.Create(ctx, &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        newHTTPDateValue(2026, time.July, 17),
	})
	_, secondCreateErr := calendarEventsService.Create(ctx, &domains.CalendarEventCreate{
		Name:        "Doctor",
		Description: "Annual checkup",
		Color:       "blue",
		Date:        newHTTPDateValue(2026, time.July, 18),
	})
	_, outsideRangeCreateErr := calendarEventsService.Create(ctx, &domains.CalendarEventCreate{
		Name:        "Flight",
		Description: "Vacation",
		Color:       "green",
		Date:        newHTTPDateValue(2026, time.July, 19),
	})
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualResponse []domains.CalendarEvent
	decodeErr := json.NewDecoder(response.Body).Decode(&actualResponse)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, outsideRangeCreateErr)
	require.NoError(t, decodeErr)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	require.Len(t, actualResponse, 2)
	assert.Equal(t, "Bills", actualResponse[0].Name)
	assert.Equal(t, "Doctor", actualResponse[1].Name)
}

func Test_CalendarEventsHandler_GetByDateRange_ShouldReturnBadRequestOnInvalidQuery(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodGet
	target := "/api/calendar-events?from=bad-date&to=2026-07-18"
	expectedBodyFragment := `invalid from value "bad-date"`
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

func Test_CalendarEventsHandler_GetByDateRange_ShouldReturnInternalServerErrorOnServiceFailure(t *testing.T) {
	// Arrange
	closedDB := newClosedDB(t)
	app := newTestApplicationWithDB(closedDB)
	method := http.MethodGet
	target := "/api/calendar-events?from=2026-07-17&to=2026-07-18"
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	// Assert
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func Test_CalendarEventsHandler_Create_ShouldReturnCreatedWithLocationAndBindTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPost
	target := "/api/calendar-events"
	locationHeaderName := "Location"
	locationPrefix := "/api/calendar-events/"
	eventTriggerTimesRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	requestBody := `{"name":"Bills","description":"Rent","color":"red","date":"2026-07-17","triggers":[{"time":"08:30:00","commentary":"Morning reminder"},{"time":"20:00:00","commentary":"Evening reminder"}]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	var actualID uuid.UUID
	decodeErr := json.NewDecoder(response.Body).Decode(&actualID)
	actualLocation := response.Header.Get(locationHeaderName)
	eventTriggerTimes, getTriggersErr := eventTriggerTimesRepo.GetByCalendarEventID(ctx, actualID)

	// Assert
	require.NoError(t, decodeErr)
	require.NoError(t, getTriggersErr)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotEqual(t, uuid.Nil, actualID)
	assert.Equal(t, locationPrefix+actualID.String(), actualLocation)
	require.Len(t, eventTriggerTimes, 2)
	assert.Equal(t, "Morning reminder", eventTriggerTimes[0].Commentary)
	assert.Equal(t, "Evening reminder", eventTriggerTimes[1].Commentary)
}

func Test_CalendarEventsHandler_Create_ShouldReturnBadRequestOnInvalidPayload(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPost
	target := "/api/calendar-events"
	requestBody := `{"name":"Bills","description":"Rent","color":"red","date":"2026-07-17","triggers":[{"time":"bad-time","commentary":"Morning reminder"}]}`
	expectedBodyFragment := `invalid time value "bad-time"`
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

func Test_CalendarEventsHandler_Update_ShouldReturnNoContentAndReconcileTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodPut
	calendarEventsService := newCalendarEventsService(testDB)
	eventTriggerTimesRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)

	calendarEventID, createErr := calendarEventsService.Create(ctx, &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        newHTTPDateValue(2026, time.July, 17),
		Triggers: []domains.CalendarEventTriggerCreate{
			{Time: newHTTPTimeValue(8, 30, 0), Commentary: "Keep me"},
			{Time: newHTTPTimeValue(9, 0, 0), Commentary: "Delete me"},
		},
	})
	currentTriggers, getTriggersErr := eventTriggerTimesRepo.GetByCalendarEventID(ctx, calendarEventID)
	require.NoError(t, createErr)
	require.NoError(t, getTriggersErr)
	require.Len(t, currentTriggers, 2)

	keptTriggerID := currentTriggers[0].Id
	target := "/api/calendar-events/" + calendarEventID.String()
	requestBody := `{"name":"Bills updated","description":"Rent updated","color":"orange","date":"2026-07-18","triggers":[{"id":"` + keptTriggerID.String() + `","time":"08:45:00","commentary":"Kept and updated"},{"time":"12:00:00","commentary":"Brand new"}]}`
	request := newJSONRequest(method, target, requestBody)

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	calendarEvents, getEventsErr := calendarEventsService.GetByDateRange(ctx, newHTTPDateValue(2026, time.July, 18), newHTTPDateValue(2026, time.July, 18))
	updatedTriggers, getUpdatedTriggersErr := eventTriggerTimesRepo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, getEventsErr)
	require.NoError(t, getUpdatedTriggersErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	require.Len(t, calendarEvents, 1)
	assert.Equal(t, "Bills updated", calendarEvents[0].Name)
	require.Len(t, updatedTriggers, 2)
	assert.Contains(t, extractHTTPTriggerCommentaries(updatedTriggers), "Kept and updated")
	assert.Contains(t, extractHTTPTriggerCommentaries(updatedTriggers), "Brand new")
	assert.NotContains(t, extractHTTPTriggerCommentaries(updatedTriggers), "Delete me")
}

func Test_CalendarEventsHandler_Update_ShouldReturnNotFoundForMissingEvent(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodPut
	missingID := uuid.New()
	target := "/api/calendar-events/" + missingID.String()
	requestBody := `{"name":"Missing","description":"","color":"red","date":"2026-07-17","triggers":[]}`
	expectedBodyFragment := "calendar event not found"
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

func Test_CalendarEventsHandler_Delete_ShouldReturnNoContentAndCascadeTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	app := newTestApplication()
	ctx := testContext
	method := http.MethodDelete
	calendarEventsService := newCalendarEventsService(testDB)
	eventTriggerTimesRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)

	calendarEventID, createErr := calendarEventsService.Create(ctx, &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        newHTTPDateValue(2026, time.July, 17),
		Triggers: []domains.CalendarEventTriggerCreate{
			{Time: newHTTPTimeValue(8, 30, 0), Commentary: "Morning reminder"},
		},
	})
	target := "/api/calendar-events/" + calendarEventID.String()
	request := newJSONRequest(method, target, "")

	// Act
	recorder := httptest.NewRecorder()
	app.router.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()

	calendarEvents, getEventsErr := calendarEventsService.GetByDateRange(ctx, newHTTPDateValue(2026, time.July, 17), newHTTPDateValue(2026, time.July, 17))
	eventTriggerTimes, getTriggersErr := eventTriggerTimesRepo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getEventsErr)
	require.NoError(t, getTriggersErr)
	assert.Equal(t, http.StatusNoContent, response.StatusCode)
	assert.Empty(t, calendarEvents)
	assert.Empty(t, eventTriggerTimes)
}

func Test_CalendarEventsHandler_Delete_ShouldReturnBadRequestOnInvalidID(t *testing.T) {
	// Arrange
	app := newTestApplication()
	method := http.MethodDelete
	target := "/api/calendar-events/bad-id"
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

func newCalendarEventsService(db *sqlx.DB) *services.CalendarEventsService {
	uow := persistence.NewUnitOfWork(db, testLogger)
	return services.NewCalendarEventsService(uow, testLogger)
}

func newHTTPDateValue(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

func newHTTPTimeValue(hours int, minutes int, seconds int) domains.TimeOnly {
	total := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return domains.TimeOnly{
		Microseconds: int64(total / time.Microsecond),
		Valid:        true,
	}
}

func extractHTTPTriggerCommentaries(triggers []domains.EventTriggerTime) []string {
	commentaries := make([]string, 0, len(triggers))
	for _, trigger := range triggers {
		commentaries = append(commentaries, trigger.Commentary)
	}

	return commentaries
}
