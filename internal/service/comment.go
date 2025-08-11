package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/pkg/errors"
)

// CommentService предоставляет бизнес-логику для комментариев
type CommentService struct {
	repo *repository.CommentRepository
	log  *logger.Logger
}

// NewCommentService создаёт новый сервис для комментариев
func NewCommentService(repo *repository.CommentRepository, log *logger.Logger) *CommentService {
	return &CommentService{repo: repo, log: log}
}

// CreateComment создаёт новый комментарий
func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	s.log.Infof("Attempting to create comment for product ID: %d by user ID: %d", comment.ProductID, comment.UserID)

	err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		s.log.Errorf("Failed to create comment for product ID %d by user ID %d: %v", comment.ProductID, comment.UserID, err)
		return fmt.Errorf("failed to create comment: %w", err)
	}

	s.log.Infof("Successfully created comment with ID: %d", comment.ID)
	return nil
}

// GetComment возвращает комментарий по ID
func (s *CommentService) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	s.log.Infof("Fetching comment with ID: %d", id)

	comment, err := s.repo.GetComment(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Comment with ID %d not found", id)
			return nil, fmt.Errorf("comment not found: %w", err)
		}
		s.log.Errorf("Failed to fetch comment with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch comment: %w", err)
	}

	s.log.Infof("Fetched comment with ID: %d, Product ID: %d", comment.ID, comment.ProductID)
	return comment, nil
}

// ListComments возвращает все комментарии
func (s *CommentService) ListComments(ctx context.Context) ([]*models.Comment, error) {
	s.log.Info("Fetching all comments from repository")

	comments, err := s.repo.ListComments(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch all comments: %v", err)
		return nil, fmt.Errorf("failed to fetch all comments: %w", err)
	}

	s.log.Infof("Successfully fetched %d comments", len(comments))
	return comments, nil
}

// UpdateComment обновляет существующий комментарий
func (s *CommentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	s.log.Infof("Updating comment with ID: %d", comment.ID)

	err := s.repo.UpdateComment(ctx, comment)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to update comment with ID %d: comment not found", comment.ID)
			return fmt.Errorf("comment with ID %d not found: %w", comment.ID, err)
		}
		s.log.Errorf("Failed to update comment with ID %d: %v", comment.ID, err)
		return fmt.Errorf("failed to update comment: %w", err)
	}

	s.log.Infof("Successfully updated comment with ID: %d", comment.ID)
	return nil
}

// DeleteComment удаляет комментарий по ID
func (s *CommentService) DeleteComment(ctx context.Context, id int) error {
	s.log.Infof("Deleting comment with ID: %d", id)

	err := s.repo.DeleteComment(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to delete comment with ID %d: comment not found", id)
			return fmt.Errorf("comment with ID %d not found: %w", id, err)
		}
		s.log.Errorf("Failed to delete comment with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	s.log.Infof("Successfully deleted comment with ID: %d", id)
	return nil
}
