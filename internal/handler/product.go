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

// ProductHandler handles requests to products.
type ProductHandler struct {
	service *service.ProductService
	log     *logger.Logger
}

// NewProductHandler creates a new ProductHandler instance.
func NewProductHandler(s *service.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{service: s, log: log}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} models.Product "Product created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Create product attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.log.Errorf("Error decoding create product request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateProduct(ctx, &product); err != nil {
		h.log.Errorf("Failed to create product, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Product created with ID: %d, request_id: %s", product.ID, requestID)
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
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid product ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Fetching product with ID: %d, request_id: %s", id, requestID)
	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to fetch product with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	h.log.Infof("Product fetched with ID: %d, request_id: %s", id, requestID)
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
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Fetching all products, request_id: %s", requestID)
	products, err := h.service.ListProducts(ctx)
	if err != nil {
		h.log.Errorf("Failed to fetch products, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Fetched %d products, request_id: %s", len(products), requestID)
	json.NewEncoder(w).Encode(products)
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
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid product ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Updating product with ID: %d, request_id: %s", id, requestID)
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		h.log.Errorf("Error decoding update product request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	product.ID = id

	if err := h.service.UpdateProduct(ctx, &product); err != nil {
		h.log.Errorf("Failed to update product with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Product updated with ID: %d, request_id: %s", id, requestID)
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
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid product ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Deleting product with ID: %d, request_id: %s", id, requestID)
	if err := h.service.DeleteProduct(ctx, id); err != nil {
		h.log.Errorf("Failed to delete product with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Product deleted with ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusNoContent)
}
