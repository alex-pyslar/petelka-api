package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/alex-pyslar/petelka-api/internal/logger"
	"github.com/alex-pyslar/petelka-api/internal/repository"
)

type PhotoService struct {
	repo *repository.PhotoRepository
	log  *logger.Logger
}

func NewPhotoService(repo *repository.PhotoRepository, log *logger.Logger) *PhotoService {
	return &PhotoService{repo: repo, log: log}
}

func (s *PhotoService) Upload(ctx context.Context, file io.Reader, size int64, filename, contentType string) (string, string, error) {
	s.log.Infof("Attempting to upload photo: %s (size: %d bytes)", filename, size)

	// Валидация
	if size <= 0 {
		s.log.Warningf("Upload rejected: invalid size %d", size)
		return "", "", fmt.Errorf("invalid file size")
	}
	if size > 32<<20 {
		s.log.Warningf("Upload rejected: file too large (%d bytes)", size)
		return "", "", fmt.Errorf("file too large: max 32MB")
	}
	ext := filepath.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		s.log.Warningf("Upload rejected: invalid file type %s", ext)
		return "", "", fmt.Errorf("invalid file type: only jpg, png, jpeg")
	}

	objectName, url, err := s.repo.Upload(ctx, file, size, filename, contentType)
	if err != nil {
		s.log.Errorf("Failed to upload photo %s: %v", filename, err)
		return "", "", fmt.Errorf("upload failed: %w", err)
	}

	s.log.Infof("Successfully uploaded photo: objectName=%s", objectName)
	return objectName, url, nil
}

func (s *PhotoService) GetDownloadURL(ctx context.Context, objectName string) (string, error) {
	s.log.Infof("Fetching download URL for objectName: %s", objectName)

	url, err := s.repo.GetPresignedURL(ctx, objectName)
	if err != nil {
		if err.Error() == "The specified key does not exist." {
			s.log.Warningf("Photo not found in MinIO: %s", objectName)
			return "", fmt.Errorf("photo not found")
		}
		s.log.Errorf("Failed to generate URL for %s: %v", objectName, err)
		return "", fmt.Errorf("failed to get download URL: %w", err)
	}

	s.log.Infof("Generated download URL for %s", objectName)
	return url, nil
}
