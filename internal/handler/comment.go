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

// CommentHandler handles requests to comments.
type CommentHandler struct {
	service *service.CommentService
	log     *logger.Logger
}

// NewCommentHandler creates a new CommentHandler instance.
func NewCommentHandler(s *service.CommentService, log *logger.Logger) *CommentHandler {
	return &CommentHandler{service: s, log: log}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment with the input payload
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment object"
// @Success 201 {object} models.Comment "Comment created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /comments [post]
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Create comment attempt from IP: %s, request_id: %s", r.RemoteAddr, requestID)

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.log.Errorf("Error decoding create comment request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateComment(ctx, &comment); err != nil {
		h.log.Errorf("Failed to create comment, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Comment created with ID: %d, request_id: %s", comment.ID, requestID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetComment godoc
// @Summary Get a comment by ID
// @Description Get details of a comment by its ID
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.Comment "Comment found"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Comment not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid comment ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Fetching comment with ID: %d, request_id: %s", id, requestID)
	comment, err := h.service.GetComment(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to fetch comment with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	h.log.Infof("Comment fetched with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(comment)
}

// ListComments godoc
// @Summary List all comments
// @Description Retrieve a list of all comments
// @Tags comments
// @Produce json
// @Success 200 {array} models.Comment "List of comments"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /comments [get]
func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	h.log.Infof("Fetching all comments, request_id: %s", requestID)
	comments, err := h.service.ListComments(ctx)
	if err != nil {
		h.log.Errorf("Failed to fetch comments, request_id: %s: %v", requestID, err)
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Fetched %d comments, request_id: %s", len(comments), requestID)
	json.NewEncoder(w).Encode(comments)
}

// UpdateComment godoc
// @Summary Update an existing comment
// @Description Update comment details by ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param comment body models.Comment true "Comment object with updated fields"
// @Success 200 {object} models.Comment "Comment updated successfully"
// @Failure 400 {string} string "Invalid request body or ID"
// @Failure 404 {string} string "Comment not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid comment ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Updating comment with ID: %d, request_id: %s", id, requestID)
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.log.Errorf("Error decoding update comment request body, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	comment.ID = id

	if err := h.service.UpdateComment(ctx, &comment); err != nil {
		h.log.Errorf("Failed to update comment with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Comment updated with ID: %d, request_id: %s", id, requestID)
	json.NewEncoder(w).Encode(comment)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags comments
// @Param id path int true "Comment ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID format"
// @Failure 404 {string} string "Comment not found"
// @Failure 500 {string} string "Internal server error"
// @Security ApiKeyAuth
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(r.Context(), "request_id", requestID)

	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id", h.log, requestID)
	if err != nil {
		h.log.Errorf("Invalid comment ID in request, request_id: %s: %v", requestID, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	h.log.Infof("Deleting comment with ID: %d, request_id: %s", id, requestID)
	if err := h.service.DeleteComment(ctx, id); err != nil {
		h.log.Errorf("Failed to delete comment with ID %d, request_id: %s: %v", id, requestID, err)
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	h.log.Infof("Comment deleted with ID: %d, request_id: %s", id, requestID)
	w.WriteHeader(http.StatusNoContent)
}
