package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

// CategoryService предоставляет бизнес-логику для категорий
type CategoryService struct {
	repo *repository.CategoryRepository
	log  *logger.Logger // Добавлено поле для логирования
}

// NewCategoryService создаёт новый сервис для категорий
func NewCategoryService(repo *repository.CategoryRepository, log *logger.Logger) *CategoryService {
	return &CategoryService{repo: repo, log: log} // Инициализация логгера
}

// CreateCategory создаёт новую категорию
func (s *CategoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	s.log.Infof("Creating category with name: %s", category.Name) // Логирование
	err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		s.log.Errorf("Failed to create category with name %s: %v", category.Name, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Category created with name: %s, ID: %d", category.Name, category.ID) // Логирование успеха
	return nil
}

// GetCategory возвращает категорию по ID
func (s *CategoryService) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	s.log.Infof("Fetching category with ID: %d", id) // Логирование
	category, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch category with ID %d: %v", id, err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched category with ID: %d, Name: %s", category.ID, category.Name) // Логирование успеха
	return category, err
}

// ListCategories возвращает все категории
func (s *CategoryService) ListCategories(ctx context.Context) ([]*models.Category, error) {
	s.log.Info("Fetching all categories") // Логирование
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch categories: %v", err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched %d categories", len(categories)) // Логирование успеха
	return categories, err
}

// UpdateCategory обновляет существующую категорию
func (s *CategoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	s.log.Infof("Updating category with ID: %d", category.ID) // Логирование
	err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		s.log.Errorf("Failed to update category with ID %d: %v", category.ID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Category updated with ID: %d", category.ID) // Логирование успеха
	return err
}

// DeleteCategory удаляет категорию по ID
func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	s.log.Infof("Deleting category with ID: %d", id) // Логирование
	err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete category with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Category deleted with ID: %d", id) // Логирование успеха
	return err
}
