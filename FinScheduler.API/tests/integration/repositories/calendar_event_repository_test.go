//go:build integration
// +build integration

package repositories_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/tests/internal/testsupport"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarEventsRepositoryCreateAndGetByDateRange_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	date := newDateValue(2026, time.July, 17)
	create := &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        date,
	}

	// Act
	calendarEventID, createErr := repo.Create(ctx, create)
	calendarEvents, getErr := repo.GetByDateRange(ctx, date, date)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.NotEqual(t, uuid.Nil, calendarEventID)
	require.Len(t, calendarEvents, 1)
	assert.Equal(t, calendarEventID, calendarEvents[0].Id)
	assert.Equal(t, create.Name, calendarEvents[0].Name)
	assert.Equal(t, create.Description, calendarEvents[0].Description)
	assert.Equal(t, create.Color, calendarEvents[0].Color)
	assert.True(t, calendarEvents[0].Date.Valid)
	assert.Equal(t, create.Date.Time.Format("2006-01-02"), calendarEvents[0].Date.Time.Format("2006-01-02"))
}

func TestCalendarEventsRepositoryGetByDateRange_ShouldReturnEventsWithinRange(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	firstID := uuid.New()
	secondID := uuid.New()
	thirdID := uuid.New()
	from := newDateValue(2026, time.July, 17)
	to := newDateValue(2026, time.July, 18)
	outside := newDateValue(2026, time.July, 19)

	insertCalendarEvent(t, firstID, "Bills", "Rent", "red", from)
	insertCalendarEvent(t, secondID, "Doctor", "Annual", "blue", to)
	insertCalendarEvent(t, thirdID, "Travel", "Flight", "green", outside)

	// Act
	calendarEvents, err := repo.GetByDateRange(ctx, from, to)

	// Assert
	require.NoError(t, err)
	require.Len(t, calendarEvents, 2)
	assert.Equal(t, firstID, calendarEvents[0].Id)
	assert.Equal(t, secondID, calendarEvents[1].Id)
	assert.NotEqual(t, thirdID, calendarEvents[0].Id)
	assert.Equal(t, from.Time.Format("2006-01-02"), calendarEvents[0].Date.Time.Format("2006-01-02"))
	assert.Equal(t, to.Time.Format("2006-01-02"), calendarEvents[1].Date.Time.Format("2006-01-02"))
}

func TestCalendarEventsRepositoryUpdate_ShouldMoveEventToNewDate(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	calendarEventID := uuid.New()
	originalDate := newDateValue(2026, time.July, 17)
	updatedDate := newDateValue(2026, time.July, 18)
	update := &domains.CalendarEventUpdate{
		Name:        "Updated Bills",
		Description: "Updated Rent",
		Color:       "orange",
		Date:        updatedDate,
	}

	insertCalendarEvent(t, calendarEventID, "Bills", "Rent", "red", originalDate)

	// Act
	updated, updateErr := repo.Update(ctx, calendarEventID, update)
	originalDateEvents, originalDateErr := repo.GetByDateRange(ctx, originalDate, originalDate)
	updatedDateEvents, updatedDateErr := repo.GetByDateRange(ctx, updatedDate, updatedDate)

	// Assert
	require.NoError(t, updateErr)
	require.NoError(t, originalDateErr)
	require.NoError(t, updatedDateErr)
	assert.True(t, updated)
	assert.Empty(t, originalDateEvents)
	require.Len(t, updatedDateEvents, 1)
	assert.Equal(t, calendarEventID, updatedDateEvents[0].Id)
	assert.Equal(t, update.Name, updatedDateEvents[0].Name)
	assert.Equal(t, update.Description, updatedDateEvents[0].Description)
	assert.Equal(t, update.Color, updatedDateEvents[0].Color)
	assert.Equal(t, update.Date.Time.Format("2006-01-02"), updatedDateEvents[0].Date.Time.Format("2006-01-02"))
}

func TestCalendarEventsRepositoryDelete_ShouldCascadeToTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	calendarEventsRepo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	eventTriggerTimesRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	calendarEventID := uuid.New()
	triggerID := uuid.New()
	date := newDateValue(2026, time.July, 17)

	insertCalendarEvent(t, calendarEventID, "Bills", "Rent", "red", date)
	insertEventTriggerTime(t, triggerID, newTimeValue(8, 30, 0), "Morning reminder", calendarEventID)

	// Act
	deleted, deleteErr := calendarEventsRepo.Delete(ctx, calendarEventID)
	calendarEvents, getEventsErr := calendarEventsRepo.GetByDateRange(ctx, date, date)
	eventTriggerTimes, getTriggersErr := eventTriggerTimesRepo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, deleteErr)
	require.NoError(t, getEventsErr)
	require.NoError(t, getTriggersErr)
	assert.True(t, deleted)
	assert.Empty(t, calendarEvents)
	assert.Empty(t, eventTriggerTimes)
}

func TestCalendarEventsRepositoryGetByDateRange_ShouldReturnErrorOnInvalidRange(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	from := newDateValue(2026, time.July, 18)
	to := newDateValue(2026, time.July, 17)

	// Act
	calendarEvents, err := repo.GetByDateRange(ctx, from, to)

	// Assert
	require.EqualError(t, err, "from should not be later than to")
	assert.Nil(t, calendarEvents)
}

func TestCalendarEventsRepositoryGetByDateRange_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewCalendarEventsRepository(closedDB, testLogger)
	from := newDateValue(2026, time.July, 17)
	to := newDateValue(2026, time.July, 18)

	// Act
	calendarEvents, err := repo.GetByDateRange(ctx, from, to)

	// Assert
	require.Error(t, err)
	assert.Nil(t, calendarEvents)
}

func insertCalendarEvent(t testing.TB, id uuid.UUID, name string, description string, color string, date pgtype.Date) {
	t.Helper()

	query := `INSERT INTO calendar_event (id, name, description, color, "date")
			  VALUES ($1, $2, $3, $4, $5)`
	_, err := testDB.Exec(query, id, name, description, color, date)
	require.NoError(t, err)
}

func newDateValue(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}
