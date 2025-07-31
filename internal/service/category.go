package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	return s.repo.CreateCategory(ctx, category)
}

func (s *CategoryService) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	return s.repo.UpdateCategory(ctx, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}
