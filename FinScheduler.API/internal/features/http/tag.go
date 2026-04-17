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

type TagsHandler struct {
	service *services.TagsService
	logger  *slog.Logger
}

func NewTagsHandler(service *services.TagsService, logger *slog.Logger) *TagsHandler {
	return &TagsHandler{
		service: service,
		logger:  logger,
	}
}

func (handler *TagsHandler) RegisterEndpoints(router chi.Router) {
	router.Get("/", handler.Get)
	router.Get("/{id}", handler.GetById)
	router.Get("/lookup", handler.GetLookup)
	router.Post("/", handler.Create)
	router.Put("/{id}", handler.Update)
}

func (handler *TagsHandler) Get(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(r.Context(), "tags-http")
	traces.RecordHttpSpan(span, r, "/tags")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /tags", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	filter := domains.NewTagsFilter(r)

	if err := filter.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	tags, count, err := handler.service.Get(ctx, &filter)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Tags filtering ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	err = json.NewEncoder(w).Encode(domains.NewPaginatedList(tags, count))
	if err != nil {
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *TagsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(r.Context(), "tags-http")
	traces.RecordHttpSpan(span, r, "/tags/{id}")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /tags/{id}", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to parse id to uuid", "id", id, "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	tag, err := handler.service.GetById(ctx, idParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			statusCode = http.StatusNotFound
			traces.EnrichFailedHttpSpan(span, err, statusCode)
			http.Error(w, "tag not found", statusCode)
			return
		}

		handler.logger.ErrorContext(ctx, "Fetching tag ended in failure", "error", err)
		statusCode = http.StatusNotFound
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(tag)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *TagsHandler) GetLookup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(r.Context(), "tags-http")
	traces.RecordHttpSpan(span, r, "/tags/lookup")
	defer func() {
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "GET /tags/lookup", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	filter := domains.NewTagsFilter(r)

	if err := filter.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		statusCode = http.StatusBadRequest
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	tags, count, err := handler.service.GetLookup(ctx, &filter)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Fetching tags lookup ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	err = json.NewEncoder(w).Encode(domains.NewPaginatedList(tags, count))
	if err != nil {
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *TagsHandler) Create(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusCreated
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(r.Context(), "tags-http")
	traces.RecordHttpSpan(span, r, "/tags")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "POST /tags", statusCode)

		if statusCode < 400 {
			traces.EnrichSuccessHttpSpan(span, statusCode)
		}
		span.End()
	}()

	w.Header().Set("Content-Type", "application/json")

	var create domains.TagCreate
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

	newTagID, err := handler.service.Create(ctx, &create)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Tag creation ended in failure", "error", err)
		statusCode = http.StatusInternalServerError
		traces.EnrichFailedHttpSpan(span, err, statusCode)
		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", r.URL.String(), newTagID))
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(newTagID); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *TagsHandler) Update(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusNoContent
	tracer := otel.Tracer("tags")
	ctx, span := tracer.Start(r.Context(), "tags-http")
	traces.RecordHttpSpan(span, r, "/tags/{id}")
	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
		metrics.RecordHTTPDuration(ctx, start)
		metrics.RecordHTTPRequest(ctx, r, "PUT /tags/{id}", statusCode)

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

	var update domains.TagUpdate
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
		statusCode = http.StatusInternalServerError
		http.Error(w, err.Error(), statusCode)
		return
	}

	if !success {
		statusCode = http.StatusNotFound
		http.Error(w, "tag not found", statusCode)
		return
	}

	w.WriteHeader(statusCode)
}
