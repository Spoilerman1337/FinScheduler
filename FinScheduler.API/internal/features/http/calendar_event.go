package featurehttp

import (
	"encoding/json"
	"finscheduler/internal/features/domains"
	"finscheduler/internal/features/services"
	"finscheduler/internal/metrics"
	"finscheduler/internal/traces"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
)

const calendarEventDateFormat = "2006-01-02"
const calendarEventTimeFormat = "15:04:05"

type CalendarEventsHandler struct {
	service *services.CalendarEventsService
	logger  *slog.Logger
}

type calendarEventDateRangeFilter struct {
	From pgtype.Date
	To   pgtype.Date
}

type calendarEventCreateRequest struct {
	Name        string                              `json:"name"`
	Description string                              `json:"description"`
	Color       string                              `json:"color"`
	Date        string                              `json:"date"`
	Triggers    []calendarEventTriggerCreateRequest `json:"triggers"`
}

type calendarEventTriggerCreateRequest struct {
	Time       string `json:"time"`
	Commentary string `json:"commentary"`
}

type calendarEventUpdateRequest struct {
	Name        string                              `json:"name"`
	Description string                              `json:"description"`
	Color       string                              `json:"color"`
	Date        string                              `json:"date"`
	Triggers    []calendarEventTriggerUpdateRequest `json:"triggers"`
}

type calendarEventTriggerUpdateRequest struct {
	Id         *string `json:"id"`
	Time       string  `json:"time"`
	Commentary string  `json:"commentary"`
}

func NewCalendarEventsHandler(service *services.CalendarEventsService, logger *slog.Logger) *CalendarEventsHandler {
	return &CalendarEventsHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *CalendarEventsHandler) RegisterEndpoints(router chi.Router) {
	router.Get("/", handler.GetByDateRange)
	router.Post("/", handler.Create)
	router.Put("/{id}", handler.Update)
	router.Delete("/{id}", handler.Delete)
}

func (handler *CalendarEventsHandler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(r.Context(), "calendar-events-http")
	traces.RecordHttpSpan(span, r, "/calendar-events")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /calendar-events", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	filter, err := newCalendarEventDateRangeFilter(r)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to parse query", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := filter.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	calendarEvents, err := handler.service.GetByDateRange(ctx, filter.From, filter.To)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Calendar events range query ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(calendarEvents); err != nil {
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *CalendarEventsHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusCreated
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(r.Context(), "calendar-events-http")
	traces.RecordHttpSpan(span, r, "/calendar-events")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "POST /calendar-events", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	var createRequest calendarEventCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	create, err := createRequest.ToDomain()
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to map request to domain", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := create.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	calendarEventID, err := handler.service.Create(ctx, create)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Calendar event creation ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", r.URL.String(), calendarEventID))
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(calendarEventID); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *CalendarEventsHandler) Update(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusNoContent
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(r.Context(), "calendar-events-http")
	traces.RecordHttpSpan(span, r, "/calendar-events/{id}")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "PUT /calendar-events/{id}", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to fetch updated entity", "id", id, "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	var updateRequest calendarEventUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	update, err := updateRequest.ToDomain()
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to map request to domain", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := update.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	success, err := handler.service.Update(ctx, idParam, update)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Calendar event update ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if !success {
		statusCode = http.StatusNotFound
		http.Error(w, "calendar event not found", statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

func (handler *CalendarEventsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusNoContent
	tracer := otel.Tracer("calendar-events")
	ctx, span := tracer.Start(r.Context(), "calendar-events-http")
	traces.RecordHttpSpan(span, r, "/calendar-events/{id}")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "DELETE /calendar-events/{id}", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to fetch deleted entity", "id", id, "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	success, err := handler.service.Delete(ctx, idParam)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Calendar event deletion ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if !success {
		statusCode = http.StatusNotFound
		http.Error(w, "calendar event not found", statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

func newCalendarEventDateRangeFilter(r *http.Request) (*calendarEventDateRangeFilter, error) {
	queryParams := r.URL.Query()

	fromRaw := queryParams.Get("from")
	if fromRaw == "" {
		return nil, fmt.Errorf("from is required")
	}

	toRaw := queryParams.Get("to")
	if toRaw == "" {
		return nil, fmt.Errorf("to is required")
	}

	from, err := parseCalendarEventDate(fromRaw, "from")
	if err != nil {
		return nil, err
	}

	to, err := parseCalendarEventDate(toRaw, "to")
	if err != nil {
		return nil, err
	}

	return &calendarEventDateRangeFilter{
		From: from,
		To:   to,
	}, nil
}

func (filter *calendarEventDateRangeFilter) Validate() error {
	if filter.From.Time.After(filter.To.Time) {
		return fmt.Errorf("from should not be later than to")
	}

	return nil
}

func (request *calendarEventCreateRequest) ToDomain() (*domains.CalendarEventCreate, error) {
	date, err := parseCalendarEventDate(request.Date, "date")
	if err != nil {
		return nil, err
	}

	triggers := make([]domains.CalendarEventTriggerCreate, 0, len(request.Triggers))
	for _, triggerRequest := range request.Triggers {
		triggerTime, parseErr := parseCalendarEventTime(triggerRequest.Time, "time")
		if parseErr != nil {
			return nil, parseErr
		}

		triggers = append(triggers, domains.CalendarEventTriggerCreate{
			Time:       triggerTime,
			Commentary: triggerRequest.Commentary,
		})
	}

	return &domains.CalendarEventCreate{
		Name:        request.Name,
		Description: request.Description,
		Color:       request.Color,
		Date:        date,
		Triggers:    triggers,
	}, nil
}

func (request *calendarEventUpdateRequest) ToDomain() (*domains.CalendarEventUpdate, error) {
	date, err := parseCalendarEventDate(request.Date, "date")
	if err != nil {
		return nil, err
	}

	triggers := make([]domains.CalendarEventTriggerUpdate, 0, len(request.Triggers))
	for _, triggerRequest := range request.Triggers {
		triggerTime, parseErr := parseCalendarEventTime(triggerRequest.Time, "time")
		if parseErr != nil {
			return nil, parseErr
		}

		var triggerID *uuid.UUID
		if triggerRequest.Id != nil {
			if *triggerRequest.Id == "" {
				return nil, fmt.Errorf("trigger id is empty")
			}

			parsedTriggerID, parseUUIDErr := uuid.Parse(*triggerRequest.Id)
			if parseUUIDErr != nil {
				return nil, fmt.Errorf("trigger id is invalid: %w", parseUUIDErr)
			}

			triggerID = &parsedTriggerID
		}

		triggers = append(triggers, domains.CalendarEventTriggerUpdate{
			Id:         triggerID,
			Time:       triggerTime,
			Commentary: triggerRequest.Commentary,
		})
	}

	return &domains.CalendarEventUpdate{
		Name:        request.Name,
		Description: request.Description,
		Color:       request.Color,
		Date:        date,
		Triggers:    triggers,
	}, nil
}

func parseCalendarEventDate(value string, fieldName string) (pgtype.Date, error) {
	parsedDate, err := time.Parse(calendarEventDateFormat, value)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("invalid %s value %q: %w", fieldName, value, err)
	}

	return pgtype.Date{
		Time:  parsedDate,
		Valid: true,
	}, nil
}

func parseCalendarEventTime(value string, fieldName string) (pgtype.Time, error) {
	parsedTime, err := time.Parse(calendarEventTimeFormat, value)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("invalid %s value %q: %w", fieldName, value, err)
	}

	total := time.Duration(parsedTime.Hour())*time.Hour +
		time.Duration(parsedTime.Minute())*time.Minute +
		time.Duration(parsedTime.Second())*time.Second

	return pgtype.Time{
		Microseconds: int64(total / time.Microsecond),
		Valid:        true,
	}, nil
}
