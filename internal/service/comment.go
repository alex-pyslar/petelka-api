package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

// CommentService предоставляет бизнес-логику для комментариев
type CommentService struct {
	repo *repository.CommentRepository
	log  *logger.Logger // Добавлено поле для логирования
}

// NewCommentService создаёт новый сервис для комментариев
func NewCommentService(repo *repository.CommentRepository, log *logger.Logger) *CommentService {
	return &CommentService{repo: repo, log: log} // Инициализация логгера
}

// CreateComment создаёт новый комментарий
func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	s.log.Infof("Creating comment for product ID: %d by user ID: %d", comment.ProductID, comment.UserID) // Логирование
	err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		s.log.Errorf("Failed to create comment for product ID %d by user ID %d: %v", comment.ProductID, comment.UserID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Comment created with ID: %d for product ID: %d", comment.ID, comment.ProductID) // Логирование успеха
	return nil
}

// GetComment возвращает комментарий по ID
func (s *CommentService) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	s.log.Infof("Fetching comment with ID: %d", id) // Логирование
	comment, err := s.repo.GetComment(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch comment with ID %d: %v", id, err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched comment with ID: %d, Product ID: %d", comment.ID, comment.ProductID) // Логирование успеха
	return comment, err
}

// ListComments возвращает все комментарии
func (s *CommentService) ListComments(ctx context.Context) ([]*models.Comment, error) {
	s.log.Info("Fetching all comments") // Логирование
	comments, err := s.repo.ListComments(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch comments: %v", err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched %d comments", len(comments)) // Логирование успеха
	return comments, err
}

// UpdateComment обновляет существующий комментарий
func (s *CommentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	s.log.Infof("Updating comment with ID: %d", comment.ID) // Логирование
	err := s.repo.UpdateComment(ctx, comment)
	if err != nil {
		s.log.Errorf("Failed to update comment with ID %d: %v", comment.ID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Comment updated with ID: %d", comment.ID) // Логирование успеха
	return err
}

// DeleteComment удаляет комментарий по ID
func (s *CommentService) DeleteComment(ctx context.Context, id int) error {
	s.log.Infof("Deleting comment with ID: %d", id) // Логирование
	err := s.repo.DeleteComment(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete comment with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Comment deleted with ID: %d", id) // Логирование успеха
	return err
}
