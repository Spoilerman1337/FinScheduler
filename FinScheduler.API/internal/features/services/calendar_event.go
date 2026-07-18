package services

import (
	"context"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/metrics"
	"finscheduler/internal/persistence"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
)

type CalendarEventsService struct {
	uow    *persistence.UnitOfWork
	logger *slog.Logger
}

const calendarEventsServiceName = "calendar-events"

func NewCalendarEventsService(uow *persistence.UnitOfWork, logger *slog.Logger) *CalendarEventsService {
	return &CalendarEventsService{
		uow:    uow,
		logger: logger,
	}
}

func (service *CalendarEventsService) Create(ctx context.Context, create *domains.CalendarEventCreate) (uuid.UUID, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-service")
	traces.RecordServiceSpan(span, "Create")
	defer span.End()

	if create == nil {
		service.logger.ErrorContext(ctx, "create is nil")
		err := fmt.Errorf("create is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Create", err)
		return uuid.Nil, err
	}

	if err := create.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "create validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Create", err)
		return uuid.Nil, err
	}

	var calendarEventID uuid.UUID
	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var repositoryErr error
		calendarEventID, repositoryErr = repositories.CalendarEvents.Create(ctx, create)
		if repositoryErr != nil {
			return repositoryErr
		}
		if calendarEventID == uuid.Nil {
			return fmt.Errorf("failed to create calendar event: repository returned nil uuid")
		}

		for _, trigger := range create.Triggers {
			_, repositoryErr = repositories.EventTriggerTimes.Create(ctx, &domains.EventTriggerTimeCreate{
				Time:            trigger.Time,
				Commentary:      trigger.Commentary,
				CalendarEventId: calendarEventID,
			})
			if repositoryErr != nil {
				return repositoryErr
			}
		}

		return nil
	})
	if err != nil {
		service.logger.ErrorContext(ctx, "error creating a calendar event", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Create", err)
		return calendarEventID, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return calendarEventID, nil
}

func (service *CalendarEventsService) Update(ctx context.Context, calendarEventID uuid.UUID, update *domains.CalendarEventUpdate) (bool, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-service")
	traces.RecordServiceSpan(span, "Update")
	defer span.End()

	if calendarEventID == uuid.Nil {
		service.logger.ErrorContext(ctx, "calendarEventID is nil")
		err := fmt.Errorf("calendarEventID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Update", err)
		return false, err
	}

	if update == nil {
		service.logger.ErrorContext(ctx, "update is nil")
		err := fmt.Errorf("update is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Update", err)
		return false, err
	}

	if err := update.Validate(); err != nil {
		service.logger.ErrorContext(ctx, "update validation failed", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Update", err)
		return false, err
	}

	var success bool
	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var repositoryErr error
		success, repositoryErr = repositories.CalendarEvents.Update(ctx, calendarEventID, update)
		if repositoryErr != nil {
			return repositoryErr
		}
		if !success {
			return nil
		}

		currentTriggers, repositoryErr := repositories.EventTriggerTimes.GetByCalendarEventID(ctx, calendarEventID)
		if repositoryErr != nil {
			return repositoryErr
		}

		retainedTriggerIDs := make(map[uuid.UUID]struct{}, len(update.Triggers))
		for _, trigger := range update.Triggers {
			if trigger.Id == nil {
				_, repositoryErr = repositories.EventTriggerTimes.Create(ctx, &domains.EventTriggerTimeCreate{
					Time:            trigger.Time,
					Commentary:      trigger.Commentary,
					CalendarEventId: calendarEventID,
				})
				if repositoryErr != nil {
					return repositoryErr
				}
				continue
			}

			retainedTriggerIDs[*trigger.Id] = struct{}{}
			updated, updateErr := repositories.EventTriggerTimes.Update(ctx, *trigger.Id, &domains.EventTriggerTimeUpdate{
				Time:            trigger.Time,
				Commentary:      trigger.Commentary,
				CalendarEventId: calendarEventID,
			})
			if updateErr != nil {
				return updateErr
			}
			if !updated {
				return fmt.Errorf("event trigger time not found: %s", trigger.Id.String())
			}
		}

		for _, currentTrigger := range currentTriggers {
			if _, exists := retainedTriggerIDs[currentTrigger.Id]; exists {
				continue
			}

			deleted, deleteErr := repositories.EventTriggerTimes.Delete(ctx, currentTrigger.Id)
			if deleteErr != nil {
				return deleteErr
			}
			if !deleted {
				return fmt.Errorf("failed to delete event trigger time: %s", currentTrigger.Id.String())
			}
		}

		return nil
	})
	if err != nil {
		service.logger.ErrorContext(ctx, "error updating a calendar event", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Update", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
}

func (service *CalendarEventsService) Delete(ctx context.Context, calendarEventID uuid.UUID) (bool, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-service")
	traces.RecordServiceSpan(span, "Delete")
	defer span.End()

	if calendarEventID == uuid.Nil {
		service.logger.ErrorContext(ctx, "calendarEventID is nil")
		err := fmt.Errorf("calendarEventID is nil")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Delete", err)
		return false, err
	}

	var success bool
	err := service.uow.WithTx(ctx, func(repositories persistence.Repositories) error {
		var repositoryErr error
		success, repositoryErr = repositories.CalendarEvents.Delete(ctx, calendarEventID)
		return repositoryErr
	})
	if err != nil {
		service.logger.ErrorContext(ctx, "error deleting a calendar event", "error", err)
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "Delete", err)
		return success, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return success, nil
}

func (service *CalendarEventsService) GetByDateRange(ctx context.Context, from pgtype.Date, to pgtype.Date) ([]domains.CalendarEvent, error) {
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(ctx, "calendar-events-service")
	traces.RecordServiceSpan(span, "GetByDateRange")
	defer span.End()

	if !from.Valid {
		service.logger.ErrorContext(ctx, "from is invalid")
		err := fmt.Errorf("from is invalid")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "GetByDateRange", err)
		return nil, err
	}

	if !to.Valid {
		service.logger.ErrorContext(ctx, "to is invalid")
		err := fmt.Errorf("to is invalid")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "GetByDateRange", err)
		return nil, err
	}

	if from.Time.After(to.Time) {
		service.logger.ErrorContext(ctx, "from is later than to")
		err := fmt.Errorf("from is later than to")
		traces.EnrichFailedServiceSpan(span, err)
		metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "GetByDateRange", err)
		return nil, err
	}

	calendarEvents := make([]domains.CalendarEvent, 0)
	err := service.uow.WithoutTx(func(repositories persistence.Repositories) error {
		rawCalendarEvents, repositoryErr := repositories.CalendarEvents.GetByDateRange(ctx, from, to)
		if repositoryErr != nil {
			service.logger.ErrorContext(ctx, "Get calendar events by date range failed", "error", repositoryErr)
			traces.EnrichFailedServiceSpan(span, repositoryErr)
			metrics.RecordServiceFailure(ctx, calendarEventsServiceName, "GetByDateRange", repositoryErr)
			return repositoryErr
		}

		calendarEvents = rawCalendarEvents
		if calendarEvents == nil {
			calendarEvents = make([]domains.CalendarEvent, 0)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	traces.EnrichSuccessServiceSpan(span)
	return calendarEvents, nil
}
