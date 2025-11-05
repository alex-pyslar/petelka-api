package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alex-pyslar/petelka-api/internal/service"
	"github.com/gorilla/mux"
)

// PhotoHandler handles HTTP requests for photos.
type PhotoHandler struct {
	service *service.PhotoService
}

// NewPhotoHandler creates a new PhotoHandler.
func NewPhotoHandler(s *service.PhotoService) *PhotoHandler {
	return &PhotoHandler{service: s}
}

// Upload godoc
// @Summary Upload a new photo
// @Description Uploads an image file (JPG/PNG/JPEG) to MinIO storage and returns objectName and presigned URL
// @Tags photos
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "Image file (max 32MB)"
// @Success 201 {object} map[string]string "objectName and url"
// @Failure 400 {string} string "Invalid file or format"
// @Failure 500 {string} string "Upload failed"
// @Security ApiKeyAuth
// @Router /photos [post]
func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Парсим форму
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "form parsing failed", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Делегируем в сервис (логирование — там!)
	objectName, url, err := h.service.Upload(r.Context(), file, header.Size, header.Filename, header.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"objectName": objectName,
		"url":        url,
	})
}

// Download godoc
// @Summary Get photo download URL
// @Description Returns a presigned URL for downloading the photo (valid for 7 days)
// @Tags photos
// @Produce  plain
// @Param objectName path string true "Object name in MinIO"
// @Success 200 {string} string "Presigned URL"
// @Failure 400 {string} string "Invalid objectName"
// @Failure 404 {string} string "Photo not found"
// @Failure 500 {string} string "Internal error"
// @Router /photos/{objectName} [get]
func (h *PhotoHandler) Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectName := vars["objectName"]
	if objectName == "" {
		http.Error(w, "objectName is required", http.StatusBadRequest)
		return
	}

	url, err := h.service.GetDownloadURL(r.Context(), objectName)
	if err != nil {
		if err.Error() == "photo not found" {
			http.Error(w, "photo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(url))
}
