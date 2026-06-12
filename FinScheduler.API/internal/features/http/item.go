package featurehttp

import (
	"database/sql"
	"encoding/json"
	"errors"
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

type ItemsHandler struct {
	service *services.ItemsService
	logger  *slog.Logger
}

func NewItemsHandler(service *services.ItemsService, logger *slog.Logger) *ItemsHandler {
	return &ItemsHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *ItemsHandler) RegisterEndpoints(router chi.Router) {
	router.Get("/", handler.GetListingInfo)
	router.Get("/{id}", handler.GetDetailedInfo)
	router.Post("/", handler.Create)
	router.Put("/{id}", handler.Update)
	router.Delete("/{id}", handler.Delete)
}

func (handler *ItemsHandler) GetListingInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(r.Context(), "items-http")
	traces.RecordHttpSpan(span, r, "/items")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /items", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	filter, err := domains.NewItemFilter(r)
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

	items, count, err := handler.service.GetListingInfo(ctx, &filter)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Items filtering ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	err = json.NewEncoder(w).Encode(domains.NewPaginatedList(items, count))
	if err != nil {
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) GetDetailedInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(r.Context(), "items-http")
	traces.RecordHttpSpan(span, r, "/items/{id}")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /items/{id}", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to parse item id", "id", id, "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	item, err := handler.service.GetDetailedInfo(ctx, idParam)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Get item by id ended in failure", "id", id, "error", err)

		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
			notFoundErr := fmt.Errorf("item not found")
			traces.EnrichFailedHttpSpan(span, notFoundErr, statusCode)
			http.Error(w, notFoundErr.Error(), statusCode)
			return
		}

		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusCreated
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(r.Context(), "items-http")
	traces.RecordHttpSpan(span, r, "/items")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "POST /items", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	var create domains.ItemCreate
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

	newItemID, err := handler.service.Create(ctx, &create)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Item creation ended in failure", "error", err)
		if errors.Is(err, domains.ErrInvalidReference) {
			statusCode = http.StatusBadRequest
			traces.EnrichFailedHttpSpan(span, err, statusCode)
			http.Error(w, err.Error(), statusCode)
			return
		}

		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", r.URL.String(), newItemID))
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(newItemID); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusNoContent
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(r.Context(), "items-http")
	traces.RecordHttpSpan(span, r, "/items/{id}")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "PUT /items/{id}", statusCode)

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

	var update domains.ItemUpdate
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
		handler.logger.ErrorContext(ctx, "database error", "error", err)
		if errors.Is(err, domains.ErrInvalidReference) {
			statusCode = http.StatusBadRequest
			traces.EnrichFailedHttpSpan(span, err, statusCode)
			http.Error(w, err.Error(), statusCode)
			return
		}

		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	if !success {
		statusCode = http.StatusNotFound
		http.Error(w, "item not found", statusCode)
		return
	}

	w.WriteHeader(statusCode)
}

func (handler *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusNoContent
	tracer := otel.Tracer("items")
	ctx, span := tracer.Start(r.Context(), "items-http")
	traces.RecordHttpSpan(span, r, "/items/{id}")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "DELETE /items/{id}", statusCode)

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
		handler.logger.ErrorContext(ctx, "database error", "error", err)
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	if !success {
		statusCode = http.StatusNotFound
		http.Error(w, "item not found", statusCode)
		return
	}

	w.WriteHeader(statusCode)
}
