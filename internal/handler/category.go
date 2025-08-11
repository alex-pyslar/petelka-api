package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/alex-pyslar/online-store/internal/models"
	"net/http"
	"strconv"

	"github.com/alex-pyslar/online-store/internal/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// CategoryHandler handles requests to categories.
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance.
func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
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
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCategory(r.Context(), &category); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	category.ID = id

	if err := h.service.UpdateCategory(r.Context(), &category); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
