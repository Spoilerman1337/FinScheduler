package repositories

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type EventTriggerTimesRepository struct {
	db     DBTX
	logger *slog.Logger
}

func NewEventTriggerTimesRepository(db DBTX, logger *slog.Logger) *EventTriggerTimesRepository {
	return &EventTriggerTimesRepository{db: db, logger: logger}
}

func (repository *EventTriggerTimesRepository) Create(ctx context.Context, create *domains.EventTriggerTimeCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("event-trigger-times")
	ctx, span := tracer.Start(ctx, "event-trigger-times-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationInsert)
	defer span.End()

	if create == nil {
		repository.logger.ErrorContext(ctx, "create should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("create should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	if !create.Time.Valid {
		repository.logger.ErrorContext(ctx, "time should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("time should be valid")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	if create.CalendarEventId == uuid.Nil {
		repository.logger.ErrorContext(ctx, "calendarEventID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("calendarEventID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	query := `INSERT INTO public.event_trigger_time (id, "time", commentary, calendar_event_id)
			  VALUES (?, ?, ?, ?)`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, newID, create.Time, create.Commentary, create.CalendarEventId)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, eventTriggerTimesTableName, err == nil, metrics.DatabaseOperationInsert)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "id", newID, "calendarEventID", create.CalendarEventId)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	affected, _ := result.RowsAffected()
	metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, true, metrics.DatabaseOperationInsert)
	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return newID, nil
}

func (repository *EventTriggerTimesRepository) GetByCalendarEventID(ctx context.Context, calendarEventID uuid.UUID) ([]domains.EventTriggerTime, error) {
	tracer := otel.Tracer("event-trigger-times")
	ctx, span := tracer.Start(ctx, "event-trigger-times-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	if calendarEventID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "calendarEventID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("calendarEventID should not be nil")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	eventTriggerTimes := make([]domains.EventTriggerTime, 0)
	query := `SELECT id,
					   "time",
					   COALESCE(commentary, '') AS commentary,
					   calendar_event_id
			  FROM public.event_trigger_time
			  WHERE calendar_event_id = ?
			  ORDER BY "time" ASC, id ASC`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "calendarEventID", calendarEventID)
	start := time.Now()
	err := sqlx.SelectContext(ctx, repository.db, &eventTriggerTimes, query, calendarEventID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, eventTriggerTimesTableName, err == nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err, "calendarEventID", calendarEventID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, true, metrics.DatabaseOperationSelect)
	traces.EnrichSuccessRepositorySpanRead(span, int64(len(eventTriggerTimes)))
	return eventTriggerTimes, nil
}

func (repository *EventTriggerTimesRepository) Update(ctx context.Context, eventTriggerTimeID uuid.UUID, update *domains.EventTriggerTimeUpdate) (bool, error) {
	tracer := otel.Tracer("event-trigger-times")
	ctx, span := tracer.Start(ctx, "event-trigger-times-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	if eventTriggerTimeID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "eventTriggerTimeID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("eventTriggerTimeID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	if update == nil {
		repository.logger.ErrorContext(ctx, "update should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("update should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	if !update.Time.Valid {
		repository.logger.ErrorContext(ctx, "time should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("time should be valid")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	if update.CalendarEventId == uuid.Nil {
		repository.logger.ErrorContext(ctx, "calendarEventID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("calendarEventID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := `UPDATE public.event_trigger_time
			  SET "time" = ?, commentary = ?, calendar_event_id = ?
			  WHERE id = ?`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "eventTriggerTimeID", eventTriggerTimeID)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, update.Time, update.Commentary, update.CalendarEventId, eventTriggerTimeID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, eventTriggerTimesTableName, err == nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "eventTriggerTimeID", eventTriggerTimeID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, true, metrics.DatabaseOperationUpdate)
	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return rowsAffected > 0, nil
}

func (repository *EventTriggerTimesRepository) Delete(ctx context.Context, eventTriggerTimeID uuid.UUID) (bool, error) {
	tracer := otel.Tracer("event-trigger-times")
	ctx, span := tracer.Start(ctx, "event-trigger-times-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationDelete)
	defer span.End()

	if eventTriggerTimeID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "eventTriggerTimeID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("eventTriggerTimeID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := "DELETE FROM public.event_trigger_time WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "eventTriggerTimeID", eventTriggerTimeID)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, eventTriggerTimeID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, eventTriggerTimesTableName, err == nil, metrics.DatabaseOperationDelete)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "error", err, "eventTriggerTimeID", eventTriggerTimeID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, eventTriggerTimesTableName, true, metrics.DatabaseOperationDelete)
	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return rowsAffected > 0, nil
}
