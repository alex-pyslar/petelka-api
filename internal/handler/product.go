package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// ProductHandler handles requests to products.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} models.Product "Product created successfully"
// @Failure 400 {string} string "Invalid request body or product type"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка поля Type
	if product.Type != "yarn" && product.Type != "garment" {
		http.Error(w, "Invalid product type: must be 'yarn' or 'garment'", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduct(r.Context(), &product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get details of a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product "Product found"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
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

	product, err := h.service.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// ListProducts godoc
// @Summary List all products
// @Description Retrieve a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// SearchProducts godoc
// @Summary Search products with filters and pagination
// @Description Search products by name, type, category_id, and color with pagination
// @Tags products
// @Produce json
// @Param name query string false "Product name (partial match)"
// @Param type query string false "Product type (yarn or garment)"
// @Param category_id query int false "Category ID"
// @Param color query string false "Product color (partial match)"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} map[string]interface{} "List of products with total count"
// @Failure 400 {string} string "Invalid query parameters"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	name := q.Get("name")
	productType := q.Get("type")
	color := q.Get("color")

	var categoryID int
	if categoryIDStr := q.Get("category_id"); categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, "Invalid category_id format", http.StatusBadRequest)
			return
		}
		categoryID = id
	}

	page := 1
	if pageStr := q.Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page format", http.StatusBadRequest)
			return
		}
		page = p
	}

	limit := 10
	if limitStr := q.Get("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit format", http.StatusBadRequest)
			return
		}
		limit = l
	}

	products, totalCount, err := h.service.SearchProducts(r.Context(), name, productType, categoryID, color, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"products":    products,
		"total_count": totalCount,
		"page":        page,
		"limit":       limit,
	}
	json.NewEncoder(w).Encode(resp)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object with updated fields"
// @Success 200 {object} models.Product "Product updated successfully"
// @Failure 400 {string} string "Invalid request body or product type"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product.ID = id

	// Проверка поля Type
	if product.Type != "yarn" && product.Type != "garment" {
		http.Error(w, "Invalid product type: must be 'yarn' or 'garment'", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateProduct(r.Context(), &product); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
