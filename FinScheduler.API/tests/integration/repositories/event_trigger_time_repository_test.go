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

func TestEventTriggerTimesRepositoryCreateAndGetByCalendarEventID_ShouldNotErr(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	calendarEventID := uuid.New()
	create := &domains.EventTriggerTimeCreate{
		Time:            newTimeValue(8, 30, 0),
		Commentary:      "Morning reminder",
		CalendarEventId: calendarEventID,
	}

	insertCalendarEvent(t, calendarEventID, "Bills", "Rent", "red", newDateValue(2026, time.July, 17))

	// Act
	eventTriggerTimeID, createErr := repo.Create(ctx, create)
	eventTriggerTimes, getErr := repo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getErr)
	require.NotEqual(t, uuid.Nil, eventTriggerTimeID)
	require.Len(t, eventTriggerTimes, 1)
	assert.Equal(t, eventTriggerTimeID, eventTriggerTimes[0].Id)
	assert.Equal(t, create.Time.Microseconds, eventTriggerTimes[0].Time.Microseconds)
	assert.True(t, eventTriggerTimes[0].Time.Valid)
	assert.Equal(t, create.Commentary, eventTriggerTimes[0].Commentary)
	assert.Equal(t, create.CalendarEventId, eventTriggerTimes[0].CalendarEventId)
}

func TestEventTriggerTimesRepositoryGetByCalendarEventID_ShouldReturnTriggersForRequestedEvent(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	firstEventID := uuid.New()
	secondEventID := uuid.New()
	firstTriggerID := uuid.New()
	secondTriggerID := uuid.New()
	thirdTriggerID := uuid.New()

	insertCalendarEvent(t, firstEventID, "Bills", "Rent", "red", newDateValue(2026, time.July, 17))
	insertCalendarEvent(t, secondEventID, "Doctor", "Annual", "blue", newDateValue(2026, time.July, 17))
	insertEventTriggerTime(t, firstTriggerID, newTimeValue(8, 30, 0), "Morning reminder", firstEventID)
	insertEventTriggerTime(t, secondTriggerID, newTimeValue(10, 0, 0), "Second reminder", firstEventID)
	insertEventTriggerTime(t, thirdTriggerID, newTimeValue(9, 0, 0), "Other event", secondEventID)

	// Act
	eventTriggerTimes, err := repo.GetByCalendarEventID(ctx, firstEventID)

	// Assert
	require.NoError(t, err)
	require.Len(t, eventTriggerTimes, 2)
	assert.Equal(t, firstTriggerID, eventTriggerTimes[0].Id)
	assert.Equal(t, newTimeValue(8, 30, 0).Microseconds, eventTriggerTimes[0].Time.Microseconds)
	assert.Equal(t, "Morning reminder", eventTriggerTimes[0].Commentary)
	assert.Equal(t, secondTriggerID, eventTriggerTimes[1].Id)
	assert.Equal(t, firstEventID, eventTriggerTimes[1].CalendarEventId)
}

func TestEventTriggerTimesRepositoryUpdate_ShouldModifyTrigger(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	calendarEventID := uuid.New()
	triggerID := uuid.New()
	update := &domains.EventTriggerTimeUpdate{
		Time:            newTimeValue(9, 15, 0),
		Commentary:      "Updated reminder",
		CalendarEventId: calendarEventID,
	}

	insertCalendarEvent(t, calendarEventID, "Bills", "Rent", "red", newDateValue(2026, time.July, 17))
	insertEventTriggerTime(t, triggerID, newTimeValue(8, 30, 0), "Morning reminder", calendarEventID)

	// Act
	updated, updateErr := repo.Update(ctx, triggerID, update)
	eventTriggerTimes, getErr := repo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, updateErr)
	require.NoError(t, getErr)
	assert.True(t, updated)
	require.Len(t, eventTriggerTimes, 1)
	assert.Equal(t, triggerID, eventTriggerTimes[0].Id)
	assert.Equal(t, update.Time.Microseconds, eventTriggerTimes[0].Time.Microseconds)
	assert.Equal(t, update.Commentary, eventTriggerTimes[0].Commentary)
	assert.Equal(t, update.CalendarEventId, eventTriggerTimes[0].CalendarEventId)
}

func TestEventTriggerTimesRepositoryDelete_ShouldRemoveTrigger(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	repo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	calendarEventID := uuid.New()
	triggerID := uuid.New()

	insertCalendarEvent(t, calendarEventID, "Bills", "Rent", "red", newDateValue(2026, time.July, 17))
	insertEventTriggerTime(t, triggerID, newTimeValue(8, 30, 0), "Morning reminder", calendarEventID)

	// Act
	deleted, deleteErr := repo.Delete(ctx, triggerID)
	eventTriggerTimes, getErr := repo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, deleteErr)
	require.NoError(t, getErr)
	assert.True(t, deleted)
	assert.Empty(t, eventTriggerTimes)
}

func TestEventTriggerTimesRepositoryGetByCalendarEventID_ShouldReturnErrorOnNilID(t *testing.T) {
	// Arrange
	ctx := testContext
	repo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)

	// Act
	eventTriggerTimes, err := repo.GetByCalendarEventID(ctx, uuid.Nil)

	// Assert
	require.EqualError(t, err, "calendarEventID should not be nil")
	assert.Nil(t, eventTriggerTimes)
}

func TestEventTriggerTimesRepositoryGetByCalendarEventID_ShouldReturnErrorWhenDatabaseIsClosed(t *testing.T) {
	// Arrange
	ctx := testContext
	closedDB := newClosedDB(t)
	repo := repositories.NewEventTriggerTimesRepository(closedDB, testLogger)
	eventID := uuid.New()

	// Act
	eventTriggerTimes, err := repo.GetByCalendarEventID(ctx, eventID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, eventTriggerTimes)
}

func insertEventTriggerTime(t testing.TB, id uuid.UUID, triggerTime pgtype.Time, commentary string, calendarEventID uuid.UUID) {
	t.Helper()

	query := `INSERT INTO event_trigger_time (id, "time", commentary, calendar_event_id)
			  VALUES ($1, $2, $3, $4)`
	_, err := testDB.Exec(query, id, triggerTime, commentary, calendarEventID)
	require.NoError(t, err)
}

func newTimeValue(hours int, minutes int, seconds int) pgtype.Time {
	total := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return pgtype.Time{
		Microseconds: int64(total / time.Microsecond),
		Valid:        true,
	}
}
