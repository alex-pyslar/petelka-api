package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) error {
	return s.repo.CreateComment(ctx, comment)
}

func (s *CommentService) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	return s.repo.GetComment(ctx, id)
}

func (s *CommentService) ListComments(ctx context.Context) ([]*models.Comment, error) {
	return s.repo.ListComments(ctx)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *models.Comment) error {
	return s.repo.UpdateComment(ctx, comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, id int) error {
	return s.repo.DeleteComment(ctx, id)
}
