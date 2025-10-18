package items

import (
	"encoding/json"
	"finscheduler/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type ItemsHandler struct {
	service *ItemsService
}

func NewItemsHandler(service *ItemsService) *ItemsHandler {
	return &ItemsHandler{service: service}
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
	w.Header().Set("Content-Type", "application/json")

	filter := NewItemFilter(r)

	items, count, err := handler.service.Get(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(shared.NewPaginatedList(items, count))
}

func (handler *ItemsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	idParam, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := handler.service.GetById(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

func (handler *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var create ItemCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := handler.service.Create(&create)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (handler *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	idParam, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var update ItemUpdate
	if err = json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := handler.service.Update(idParam, &update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}

func (handler *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")

	idParam, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := handler.service.Delete(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
}
