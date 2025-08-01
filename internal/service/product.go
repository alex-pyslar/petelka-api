package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

// ProductService предоставляет бизнес-логику для товаров
type ProductService struct {
	repo *repository.ProductRepository
	log  *logger.Logger // Добавлено поле для логирования
}

// NewProductService создаёт новый сервис для товаров
func NewProductService(repo *repository.ProductRepository, log *logger.Logger) *ProductService {
	return &ProductService{repo: repo, log: log} // Инициализация логгера
}

// CreateProduct создаёт новый товар
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Creating product with name: %s", product.Name) // Логирование
	err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		s.log.Errorf("Failed to create product with name %s: %v", product.Name, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Product created with name: %s, ID: %d", product.Name, product.ID) // Логирование успеха
	return nil
}

// GetProduct возвращает товар по ID
func (s *ProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	s.log.Infof("Fetching product with ID: %d", id) // Логирование
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch product with ID %d: %v", id, err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched product with ID: %d, Name: %s", product.ID, product.Name) // Логирование успеха
	return product, err
}

// ListProducts возвращает все товары
func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	s.log.Info("Fetching all products") // Логирование
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch products: %v", err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched %d products", len(products)) // Логирование успеха
	return products, err
}

// UpdateProduct обновляет существующий товар
func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Updating product with ID: %d", product.ID) // Логирование
	err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		s.log.Errorf("Failed to update product with ID %d: %v", product.ID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Product updated with ID: %d", product.ID) // Логирование успеха
	return err
}

// DeleteProduct удаляет товар по ID
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	s.log.Infof("Deleting product with ID: %d", id) // Логирование
	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete product with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Product deleted with ID: %d", id) // Логирование успеха
	return err
}
