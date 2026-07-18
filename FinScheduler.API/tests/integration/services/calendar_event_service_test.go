//go:build integration
// +build integration

package services_test

import (
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/repositories"
	"finscheduler/internal/features/services"
	"finscheduler/internal/persistence"
	"finscheduler/tests/internal/testsupport"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalendarEventsServiceCreate_ShouldCreateEventWithTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	eventTriggersRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	date := newServiceDateValue(2026, time.July, 17)
	create := &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        date,
		Triggers: []domains.CalendarEventTriggerCreate{
			{Time: newServiceTimeValue(8, 30, 0), Commentary: "Morning reminder"},
			{Time: newServiceTimeValue(20, 0, 0), Commentary: "Evening reminder"},
		},
	}

	// Act
	calendarEventID, createErr := service.Create(ctx, create)
	calendarEvents, getEventsErr := service.GetByDateRange(ctx, date, date)
	eventTriggerTimes, getTriggersErr := eventTriggersRepo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, createErr)
	require.NoError(t, getEventsErr)
	require.NoError(t, getTriggersErr)
	require.NotEqual(t, uuid.Nil, calendarEventID)
	require.Len(t, calendarEvents, 1)
	require.Len(t, eventTriggerTimes, 2)
	assert.Equal(t, calendarEventID, calendarEvents[0].Id)
	assert.Equal(t, create.Name, calendarEvents[0].Name)
	assert.Equal(t, []string{"Morning reminder", "Evening reminder"}, []string{eventTriggerTimes[0].Commentary, eventTriggerTimes[1].Commentary})
}

func TestCalendarEventsServiceCreate_ShouldRollbackEventWhenTriggerInsertFails(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	date := newServiceDateValue(2026, time.July, 17)
	create := &domains.CalendarEventCreate{
		Name:        "Bills",
		Description: "Rent",
		Color:       "red",
		Date:        date,
		Triggers: []domains.CalendarEventTriggerCreate{
			{Time: newServiceTimeValue(8, 30, 0), Commentary: "First reminder"},
			{Time: newServiceTimeValue(8, 30, 0), Commentary: "Duplicated reminder"},
		},
	}

	// Act
	calendarEventID, createErr := service.Create(ctx, create)
	calendarEvents, getEventsErr := service.GetByDateRange(ctx, date, date)

	// Assert
	require.Error(t, createErr)
	require.NoError(t, getEventsErr)
	assert.NotEqual(t, uuid.Nil, calendarEventID)
	assert.Empty(t, calendarEvents)
}

func TestCalendarEventsServiceUpdate_ShouldReconcileTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	calendarEventsRepo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	eventTriggersRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	originalDate := newServiceDateValue(2026, time.July, 17)
	updatedDate := newServiceDateValue(2026, time.July, 18)
	firstEventCreate := &domains.CalendarEventCreate{Name: "Bills", Description: "Rent", Color: "red", Date: originalDate}
	secondEventCreate := &domains.CalendarEventCreate{Name: "Doctor", Description: "Annual", Color: "blue", Date: originalDate}

	firstEventID, firstEventCreateErr := calendarEventsRepo.Create(ctx, firstEventCreate)
	secondEventID, secondEventCreateErr := calendarEventsRepo.Create(ctx, secondEventCreate)
	keptTriggerID, keptTriggerErr := eventTriggersRepo.Create(ctx, &domains.EventTriggerTimeCreate{Time: newServiceTimeValue(8, 30, 0), Commentary: "Keep me", CalendarEventId: firstEventID})
	deletedTriggerID, deletedTriggerErr := eventTriggersRepo.Create(ctx, &domains.EventTriggerTimeCreate{Time: newServiceTimeValue(9, 0, 0), Commentary: "Delete me", CalendarEventId: firstEventID})
	attachedTriggerID, attachedTriggerErr := eventTriggersRepo.Create(ctx, &domains.EventTriggerTimeCreate{Time: newServiceTimeValue(10, 0, 0), Commentary: "Attach me", CalendarEventId: secondEventID})

	update := &domains.CalendarEventUpdate{
		Name:        "Bills updated",
		Description: "Rent updated",
		Color:       "orange",
		Date:        updatedDate,
		Triggers: []domains.CalendarEventTriggerUpdate{
			{Id: &keptTriggerID, Time: newServiceTimeValue(8, 45, 0), Commentary: "Kept and updated"},
			{Id: &attachedTriggerID, Time: newServiceTimeValue(10, 30, 0), Commentary: "Attached to first"},
			{Time: newServiceTimeValue(12, 0, 0), Commentary: "Brand new"},
		},
	}

	// Act
	updated, updateErr := service.Update(ctx, firstEventID, update)
	updatedEvents, getEventsErr := service.GetByDateRange(ctx, updatedDate, updatedDate)
	firstEventTriggers, firstTriggersErr := eventTriggersRepo.GetByCalendarEventID(ctx, firstEventID)
	secondEventTriggers, secondTriggersErr := eventTriggersRepo.GetByCalendarEventID(ctx, secondEventID)

	// Assert
	require.NoError(t, firstEventCreateErr)
	require.NoError(t, secondEventCreateErr)
	require.NoError(t, keptTriggerErr)
	require.NoError(t, deletedTriggerErr)
	require.NoError(t, attachedTriggerErr)
	require.NoError(t, updateErr)
	require.NoError(t, getEventsErr)
	require.NoError(t, firstTriggersErr)
	require.NoError(t, secondTriggersErr)
	require.True(t, updated)
	require.Len(t, updatedEvents, 1)
	require.Len(t, firstEventTriggers, 3)
	assert.Equal(t, firstEventID, updatedEvents[0].Id)
	assert.Equal(t, update.Name, updatedEvents[0].Name)
	assert.Empty(t, secondEventTriggers)
	assert.NotContains(t, extractTriggerIDs(firstEventTriggers), deletedTriggerID)
	assert.Contains(t, extractTriggerIDs(firstEventTriggers), keptTriggerID)
	assert.Contains(t, extractTriggerIDs(firstEventTriggers), attachedTriggerID)
	assert.True(t, containsTriggerCommentary(firstEventTriggers, "Brand new"))
	assert.True(t, containsTriggerCommentary(firstEventTriggers, "Kept and updated"))
	assert.True(t, containsTriggerCommentary(firstEventTriggers, "Attached to first"))
}

func TestCalendarEventsServiceDelete_ShouldCascadeTriggers(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	calendarEventsRepo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	eventTriggersRepo := repositories.NewEventTriggerTimesRepository(testDB, testLogger)
	date := newServiceDateValue(2026, time.July, 17)
	calendarEventCreate := &domains.CalendarEventCreate{Name: "Bills", Description: "Rent", Color: "red", Date: date}

	calendarEventID, eventCreateErr := calendarEventsRepo.Create(ctx, calendarEventCreate)
	_, triggerCreateErr := eventTriggersRepo.Create(ctx, &domains.EventTriggerTimeCreate{Time: newServiceTimeValue(8, 30, 0), Commentary: "Morning reminder", CalendarEventId: calendarEventID})

	// Act
	deleted, deleteErr := service.Delete(ctx, calendarEventID)
	calendarEvents, getEventsErr := service.GetByDateRange(ctx, date, date)
	eventTriggerTimes, getTriggersErr := eventTriggersRepo.GetByCalendarEventID(ctx, calendarEventID)

	// Assert
	require.NoError(t, eventCreateErr)
	require.NoError(t, triggerCreateErr)
	require.NoError(t, deleteErr)
	require.NoError(t, getEventsErr)
	require.NoError(t, getTriggersErr)
	require.True(t, deleted)
	assert.Empty(t, calendarEvents)
	assert.Empty(t, eventTriggerTimes)
}

func TestCalendarEventsServiceGetByDateRange_ShouldReturnEventsInRange(t *testing.T) {
	// Arrange
	t.Cleanup(func() {
		testsupport.Truncate(t, testDB, "calendar_event")
	})

	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	calendarEventsRepo := repositories.NewCalendarEventsRepository(testDB, testLogger)
	from := newServiceDateValue(2026, time.July, 17)
	middle := newServiceDateValue(2026, time.July, 18)
	to := newServiceDateValue(2026, time.July, 19)
	outside := newServiceDateValue(2026, time.July, 20)

	_, firstCreateErr := calendarEventsRepo.Create(ctx, &domains.CalendarEventCreate{Name: "Bills", Description: "Rent", Color: "red", Date: from})
	_, secondCreateErr := calendarEventsRepo.Create(ctx, &domains.CalendarEventCreate{Name: "Doctor", Description: "Annual", Color: "blue", Date: middle})
	_, thirdCreateErr := calendarEventsRepo.Create(ctx, &domains.CalendarEventCreate{Name: "Travel", Description: "Flight", Color: "green", Date: outside})

	// Act
	calendarEvents, getErr := service.GetByDateRange(ctx, from, to)

	// Assert
	require.NoError(t, firstCreateErr)
	require.NoError(t, secondCreateErr)
	require.NoError(t, thirdCreateErr)
	require.NoError(t, getErr)
	require.Len(t, calendarEvents, 2)
	assert.Equal(t, []string{"Bills", "Doctor"}, []string{calendarEvents[0].Name, calendarEvents[1].Name})
}

func TestCalendarEventsServiceUpdateMissing_ShouldReturnFalseWithoutErr(t *testing.T) {
	// Arrange
	ctx := testContext
	uow := persistence.NewUnitOfWork(testDB, testLogger)
	service := services.NewCalendarEventsService(uow, testLogger)
	missingID := uuid.New()
	update := &domains.CalendarEventUpdate{Name: "Missing", Color: "red", Date: newServiceDateValue(2026, time.July, 17)}

	// Act
	updated, err := service.Update(ctx, missingID, update)

	// Assert
	require.NoError(t, err)
	assert.False(t, updated)
}

func extractTriggerIDs(triggers []domains.EventTriggerTime) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(triggers))
	for _, trigger := range triggers {
		ids = append(ids, trigger.Id)
	}

	return ids
}

func containsTriggerCommentary(triggers []domains.EventTriggerTime, commentary string) bool {
	for _, trigger := range triggers {
		if trigger.Commentary == commentary {
			return true
		}
	}

	return false
}

func newServiceDateValue(year int, month time.Month, day int) pgtype.Date {
	return pgtype.Date{
		Time:  time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
}

func newServiceTimeValue(hours int, minutes int, seconds int) pgtype.Time {
	total := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	return pgtype.Time{
		Microseconds: int64(total / time.Microsecond),
		Valid:        true,
	}
}
