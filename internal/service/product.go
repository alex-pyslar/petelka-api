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
	log  *logger.Logger
}

// NewProductService создаёт новый сервис для товаров
func NewProductService(repo *repository.ProductRepository, log *logger.Logger) *ProductService {
	return &ProductService{repo: repo, log: log}
}

// CreateProduct создаёт новый товар
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Creating product with name: %s", product.Name)
	err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		s.log.Errorf("Failed to create product with name %s: %v", product.Name, err)
		return err
	}
	s.log.Infof("Product created with name: %s, ID: %d", product.Name, product.ID)
	return nil
}

// GetProduct возвращает товар по ID
func (s *ProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	s.log.Infof("Fetching product with ID: %d", id)
	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch product with ID %d: %v", id, err)
		return nil, err
	}
	s.log.Infof("Fetched product with ID: %d, Name: %s", product.ID, product.Name)
	return product, err
}

// ListProducts возвращает все товары
func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	s.log.Info("Fetching all products")
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch products: %v", err)
		return nil, err
	}
	s.log.Infof("Fetched %d products", len(products))
	return products, err
}

// ListProductsByCategory возвращает товары по ID категории
func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int) ([]*models.Product, error) {
	s.log.Infof("Fetching products for category ID: %d", categoryID)
	products, err := s.repo.ListProductsByCategory(ctx, categoryID)
	if err != nil {
		s.log.Errorf("Failed to fetch products for category ID %d: %v", categoryID, err)
		return nil, err
	}
	s.log.Infof("Fetched %d products for category ID: %d", len(products), categoryID)
	return products, err
}

// UpdateProduct обновляет существующий товар
func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Updating product with ID: %d", product.ID)
	err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		s.log.Errorf("Failed to update product with ID %d: %v", product.ID, err)
		return err
	}
	s.log.Infof("Product updated with ID: %d", product.ID)
	return err
}

// DeleteProduct удаляет товар по ID
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	s.log.Infof("Deleting product with ID: %d", id)
	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete product with ID %d: %v", id, err)
		return err
	}
	s.log.Infof("Product deleted with ID: %d", id)
	return err
}
