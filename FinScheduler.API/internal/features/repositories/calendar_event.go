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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

type CalendarEventsRepository struct {
	db     DBTX
	logger *slog.Logger
}

func NewCalendarEventsRepository(db DBTX, logger *slog.Logger) *CalendarEventsRepository {
	return &CalendarEventsRepository{db: db, logger: logger}
}

func (repository *CalendarEventsRepository) Create(ctx context.Context, create *domains.CalendarEventCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationInsert)
	defer span.End()

	if create == nil {
		repository.logger.ErrorContext(ctx, "create should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("create should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	if !create.Date.Valid {
		repository.logger.ErrorContext(ctx, "date should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("date should be valid")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	newID, err := uuid.NewV7()
	if err != nil {
		repository.logger.ErrorContext(ctx, "uuid generation error", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	query := `INSERT INTO public.calendar_event (id, name, description, color, "date")
			  VALUES (?, ?, ?, ?, ?)`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, newID, create.Name, create.Description, create.Color, create.Date)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, calendarEventsTableName, err == nil, metrics.DatabaseOperationInsert)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on INSERT operation", "error", err, "id", newID, "name", create.Name)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationInsert)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return uuid.Nil, err
	}

	affected, _ := result.RowsAffected()
	metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, true, metrics.DatabaseOperationInsert)
	traces.EnrichSuccessRepositorySpanWrite(span, affected)
	return newID, nil
}

func (repository *CalendarEventsRepository) GetByDateRange(ctx context.Context, from pgtype.Date, to pgtype.Date) ([]domains.CalendarEvent, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationSelect)
	defer span.End()

	if !from.Valid {
		repository.logger.ErrorContext(ctx, "from should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("from should be valid")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	if !to.Valid {
		repository.logger.ErrorContext(ctx, "to should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("to should be valid")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	if from.Time.After(to.Time) {
		repository.logger.ErrorContext(ctx, "from should not be later than to")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("from should not be later than to")
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	calendarEvents := make([]domains.CalendarEvent, 0)
	query := `SELECT id,
					   name,
					   COALESCE(description, '') AS description,
					   color,
					   "date"
			  FROM public.calendar_event
			  WHERE "date" >= ? AND "date" <= ?
			  ORDER BY "date" ASC, name ASC, id ASC`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query,
		"from", from.Time.Format("2006-01-02"),
		"to", to.Time.Format("2006-01-02"),
	)
	start := time.Now()
	err := sqlx.SelectContext(ctx, repository.db, &calendarEvents, query, from, to)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, calendarEventsTableName, err == nil, metrics.DatabaseOperationSelect)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on SELECT operation", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationSelect)
		traces.EnrichFailedRepositorySpanRead(span, err, 0)
		return nil, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, true, metrics.DatabaseOperationSelect)
	traces.EnrichSuccessRepositorySpanRead(span, int64(len(calendarEvents)))
	return calendarEvents, nil
}

func (repository *CalendarEventsRepository) Update(ctx context.Context, calendarEventID uuid.UUID, update *domains.CalendarEventUpdate) (bool, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationUpdate)
	defer span.End()

	if calendarEventID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "calendarEventID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("calendarEventID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	if update == nil {
		repository.logger.ErrorContext(ctx, "update should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("update should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	if !update.Date.Valid {
		repository.logger.ErrorContext(ctx, "date should be valid")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("date should be valid")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := `UPDATE public.calendar_event
			  SET name = ?, description = ?, color = ?, "date" = ?
			  WHERE id = ?`
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "calendarEventID", calendarEventID)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, update.Name, update.Description, update.Color, update.Date, calendarEventID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, calendarEventsTableName, err == nil, metrics.DatabaseOperationUpdate)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on UPDATE operation", "error", err, "calendarEventID", calendarEventID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationUpdate)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, true, metrics.DatabaseOperationUpdate)
	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return rowsAffected > 0, nil
}

func (repository *CalendarEventsRepository) Delete(ctx context.Context, calendarEventID uuid.UUID) (bool, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-repository")
	traces.RecordRepositorySpan(span, databaseDriver, metrics.DatabaseOperationDelete)
	defer span.End()

	if calendarEventID == uuid.Nil {
		repository.logger.ErrorContext(ctx, "calendarEventID should not be nil")
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationNone)

		err := fmt.Errorf("calendarEventID should not be nil")
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	query := "DELETE FROM public.calendar_event WHERE id = ?"
	query = repository.db.Rebind(query)

	repository.logger.InfoContext(ctx, "executing operation:", "query", query, "calendarEventID", calendarEventID)
	start := time.Now()
	result, err := repository.db.ExecContext(ctx, query, calendarEventID)
	metrics.RecordDatabaseDuration(ctx, start, databaseDriver, calendarEventsTableName, err == nil, metrics.DatabaseOperationDelete)
	if err != nil {
		repository.logger.ErrorContext(ctx, "error on DELETE operation", "error", err, "calendarEventID", calendarEventID)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		repository.logger.ErrorContext(ctx, "error fetching affected rows", "error", err)
		metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, false, metrics.DatabaseOperationDelete)
		traces.EnrichFailedRepositorySpanWrite(span, err, 0)
		return false, err
	}

	metrics.RecordDatabaseRequest(ctx, databaseDriver, calendarEventsTableName, true, metrics.DatabaseOperationDelete)
	traces.EnrichSuccessRepositorySpanWrite(span, rowsAffected)
	return rowsAffected > 0, nil
}
