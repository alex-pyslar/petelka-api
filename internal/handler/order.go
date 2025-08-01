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

// OrderHandler handles requests to orders.
type OrderHandler struct {
	service *service.OrderService
	log     *logger.Logger
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(s *service.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{service: s, log: log}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with the input payload
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order object"
// @Success 201 {object} models.Order "Order created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Create order attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.log.Errorf("Error decoding create order request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateOrder(ctx, &order); err != nil {
		h.log.Errorf("Failed to create order, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Order created with ID: %d, request_id: %s", order.ID, requestID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get details of an order by its ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order "Order found"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid order ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Fetching order with ID: %d, request_id: %s", id, requestID)
	order, err := h.service.GetOrder(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to fetch order with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	h.log.Infof("Order fetched with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(order)
}

// ListOrders godoc
// @Summary List all orders
// @Description Retrieve a list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order "List of orders"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Fetching all orders, request_id: %s", requestID)
	orders, err := h.service.ListOrders(ctx)
	if err != nil {
		h.log.Errorf("Failed to fetch orders, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Fetched %d orders, request_id: %s", len(orders), requestID)
	json.NewEncoder(w).Encode(orders)
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update order details by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body models.Order true "Order object with updated fields"
// @Success 200 {object} models.Order "Order updated successfully"
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid order ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Updating order with ID: %d, request_id: %s", id, requestID)
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.log.Errorf("Error decoding update order request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order.ID = id

	if err := h.service.UpdateOrder(ctx, &order); err != nil {
		h.log.Errorf("Failed to update order with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Order updated with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(order)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Param id path int true "Order ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid order ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Deleting order with ID: %d, request_id: %s", id, requestID)
	if err := h.service.DeleteOrder(ctx, id); err != nil {
		h.log.Errorf("Failed to delete order with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Order deleted with ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusNoContent)
}
