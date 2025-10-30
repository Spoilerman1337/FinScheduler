package items

import (
	"encoding/json"
	"finscheduler/internal/shared"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type ItemsHandler struct {
	service *ItemsService
	logger  *slog.Logger
}

func NewItemsHandler(service *ItemsService, logger *slog.Logger) *ItemsHandler {
	return &ItemsHandler{service: service, logger: logger}
}

func (handler *ItemsHandler) RegisterEndpoints() chi.Router {
	router := chi.NewRouter()
	router.Get("/", handler.Get)
	router.Get("/{id}", handler.GetById)
	router.Post("/", handler.Create)
	router.Put("/{id}", handler.Update)
	router.Delete("/{id}", handler.Delete)

	return router
}

func (handler *ItemsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	filter := NewItemFilter(r)

	if err := filter.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, count, err := handler.service.Get(ctx, &filter)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Items filtering ended in failure", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(shared.NewPaginatedList(items, count))
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to parse id to uuid", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := handler.service.GetById(ctx, idParam)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Fetching item ended in failure", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	defer func() {
		err := r.Body.Close()
		if err != nil {
			handler.logger.ErrorContext(ctx, "Failed to close request body", "error", err)
		}
	}()

	var create ItemCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := create.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newItemID, err := handler.service.Create(ctx, &create)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Item creation ended in failure", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", r.URL.String(), newItemID))
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newItemID); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to encode result", "error", err)
		return
	}
}

func (handler *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to fetch updated entity", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var update ItemUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		handler.logger.ErrorContext(ctx, "Failed to decode body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := update.Validate(); err != nil {
		handler.logger.ErrorContext(ctx, "Validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success, err := handler.service.Update(ctx, idParam, &update)
	if err != nil || !success {
		handler.logger.ErrorContext(ctx, "Item updated ended in failure", "error", err, "success", success)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	idParam, err := uuid.Parse(id)
	if err != nil {
		handler.logger.ErrorContext(ctx, "Failed to fetch deleted entity", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success, err := handler.service.Delete(ctx, idParam)
	if err != nil || !success {
		handler.logger.ErrorContext(ctx, "Item deletion ended in failure", "error", err, "success", success)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
