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
	"go.opentelemetry.io/otel"
)

type CalendarEventsHandler struct {
	service *services.CalendarEventsService
	logger  *slog.Logger
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

	filter, err := domains.NewCalendarEventDateRangeFilter(r)
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

	var create domains.CalendarEventCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
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

	calendarEventID, err := handler.service.Create(ctx, &create)
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

	var update domains.CalendarEventUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
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

	success, err := handler.service.Update(ctx, idParam, &update)
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
