package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CategoryHandler handles requests to categories.
type CategoryHandler struct {
	service *service.CategoryService
	log     *logger.Logger
}

// NewCategoryHandler creates a new CategoryHandler instance.
func NewCategoryHandler(s *service.CategoryService, log *logger.Logger) *CategoryHandler {
	return &CategoryHandler{service: s, log: log}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the input payload
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} models.Category "Category created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Create category attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.log.Errorf("Error decoding create category request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCategory(ctx, &category); err != nil {
		h.log.Errorf("Failed to create category, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Category created with ID: %d, request_id: %s", category.ID, requestID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetCategory godoc
// @Summary Get a category by ID
// @Description Get details of a category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category "Category found"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid category ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Fetching category with ID: %d, request_id: %s", id, requestID)
	category, err := h.service.GetCategory(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to fetch category with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	h.log.Infof("Category fetched with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(category)
}

// ListCategories godoc
// @Summary List all categories
// @Description Retrieve a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category "List of categories"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Fetching all categories, request_id: %s", requestID)
	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		h.log.Errorf("Failed to fetch categories, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Fetched %d categories, request_id: %s", len(categories), requestID)
	json.NewEncoder(w).Encode(categories)
}

// UpdateCategory godoc
// @Summary Update an existing category
// @Description Update category details by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category object with updated fields"
// @Success 200 {object} models.Category "Category updated successfully"
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid category ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Updating category with ID: %d, request_id: %s", id, requestID)
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.log.Errorf("Error decoding update category request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	category.ID = id

	if err := h.service.UpdateCategory(ctx, &category); err != nil {
		h.log.Errorf("Failed to update category with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Category updated with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Category not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid category ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Deleting category with ID: %d, request_id: %s", id, requestID)
	if err := h.service.DeleteCategory(ctx, id); err != nil {
		h.log.Errorf("Failed to delete category with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Category deleted with ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusNoContent)
}
