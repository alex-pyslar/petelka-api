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

// CommentHandler handles requests to comments.
type CommentHandler struct {
	service *service.CommentService
}

// NewCommentHandler creates a new CommentHandler instance.
func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
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
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	comment, err := h.service.GetComment(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	comments, err := h.service.ListComments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	comment.ID = id

	if err := h.service.UpdateComment(r.Context(), &comment); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	if err := h.service.DeleteComment(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
