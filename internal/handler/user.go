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

// UserHandler handles requests to users.
type UserHandler struct {
	service *service.UserService
	log     *logger.Logger
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(s *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{service: s, log: log}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Create user attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Errorf("Error decoding create user request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(ctx, &user); err != nil {
		h.log.Errorf("Failed to create user, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	h.log.Infof("User created with ID: %d, request_id: %s", user.ID, requestID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a user's details by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User found"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid user ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Fetching user with ID: %d, request_id: %s", id, requestID)
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to fetch user with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.log.Infof("User fetched with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(user)
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieve a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User "List of users"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Fetching all users, request_id: %s", requestID)
	users, err := h.service.ListUsers(ctx)
	if err != nil {
		h.log.Errorf("Failed to fetch users, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Fetched %d users, request_id: %s", len(users), requestID)
	json.NewEncoder(w).Encode(users)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User object with updated fields"
// @Success 200 {object} models.User "User updated successfully"
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid user ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Updating user with ID: %d, request_id: %s", id, requestID)
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Errorf("Error decoding update user request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user.ID = id

	if err := h.service.UpdateUser(ctx, &user); err != nil {
		h.log.Errorf("Failed to update user with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	h.log.Infof("User updated with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid user ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Deleting user with ID: %d, request_id: %s", id, requestID)
	if err := h.service.DeleteUser(ctx, id); err != nil {
		h.log.Errorf("Failed to delete user with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	h.log.Infof("User deleted with ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusNoContent)
}

// VerifyPassword godoc
// @Summary Verify user password
// @Description Verify a user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param password body VerifyPasswordRequest true "Password to verify"
// @Success 200 {string} string "Password verified"
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 401 {string} string "Invalid password"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /users/{id}/verify-password [post]
func (h *UserHandler) VerifyPassword(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid user ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input VerifyPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Errorf("Error decoding verify password request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.log.Infof("Verifying password for user ID: %d, request_id: %s", id, requestID)
	if err := h.service.VerifyPassword(ctx, id, input.Password); err != nil {
		h.log.Errorf("Password verification failed for user ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	h.log.Infof("Password verified for user ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password verified"))
}

// VerifyPasswordRequest is the request body for password verification.
type VerifyPasswordRequest struct {
	Password string `json:"password"`
}
