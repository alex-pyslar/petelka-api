package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alex-pyslar/petelka-api/internal/logger"
	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/pkg/errors"
)

// CategoryService предоставляет бизнес-логику для категорий
type CategoryService struct {
	repo *repository.CategoryRepository
	log  *logger.Logger
}

// NewCategoryService создаёт новый сервис для категорий
func NewCategoryService(repo *repository.CategoryRepository, log *logger.Logger) *CategoryService {
	return &CategoryService{repo: repo, log: log}
}

// CreateCategory создаёт новую категорию
func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	s.log.Infof("Attempting to create category with name: %s", category.Name)

	err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		s.log.Errorf("Failed to create category: %v", err)
		return fmt.Errorf("failed to create category: %w", err)
	}

	s.log.Infof("Successfully created category with ID: %d", category.ID)
	return nil
}

// GetCategory возвращает категорию по ID
func (s *CategoryService) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	s.log.Infof("Fetching category with ID: %d", id)

	category, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Category with ID %d not found", id)
			return nil, fmt.Errorf("category not found: %w", err)
		}
		s.log.Errorf("Failed to fetch category with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}

	s.log.Infof("Fetched category with ID: %d", category.ID)
	return category, nil
}

// ListCategories возвращает все категории
func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.Category, error) {
	s.log.Info("Fetching all categories from repository")

	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch all categories: %v", err)
		return nil, fmt.Errorf("failed to fetch all categories: %w", err)
	}

	s.log.Infof("Successfully fetched %d categories", len(categories))
	return categories, nil
}

// UpdateCategory обновляет существующую категорию
func (s *CategoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	s.log.Infof("Updating category with ID: %d", category.ID)

	err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to update category with ID %d: category not found", category.ID)
			return fmt.Errorf("category with ID %d not found: %w", category.ID, err)
		}
		s.log.Errorf("Failed to update category with ID %d: %v", category.ID, err)
		return fmt.Errorf("failed to update category: %w", err)
	}

	s.log.Infof("Successfully updated category with ID: %d", category.ID)
	return nil
}

// DeleteCategory удаляет категорию по ID
func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	s.log.Infof("Deleting category with ID: %d", id)

	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to delete category with ID %d: category not found", id)
			return fmt.Errorf("category with ID %d not found: %w", id, err)
		}
		s.log.Errorf("Failed to delete category with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete category: %w", err)
	}

	s.log.Infof("Successfully deleted category with ID: %d", id)
	return nil
}
