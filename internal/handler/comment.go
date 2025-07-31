package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/gorilla/mux"
)

// @Summary Create a new comment
// @Description Create a new comment with the input payload
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment object"
// @Success 201 {object} models.Comment
// @Router /comments [post]
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// @Summary Get a comment by ID
// @Description Get details of a comment by ID
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} models.Comment
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	comment, err := h.service.GetComment(r.Context(), id)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// @Summary List all comments
// @Description Get a list of all comments
// @Tags comments
// @Produce json
// @Success 200 {array} models.Comment
// @Router /comments [get]
func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.service.ListComments(r.Context())
	if err != nil {
		http.Error(w, "Failed to list comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// @Summary Update a comment
// @Description Update details of a comment by ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param comment body models.Comment true "Comment object"
// @Success 200 {object} models.Comment
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	comment.ID = id

	if err := h.service.UpdateComment(r.Context(), &comment); err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags comments
// @Param id path int true "Comment ID"
// @Success 204
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := getIDFromVars(vars, "id")
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteComment(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
