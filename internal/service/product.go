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

// ProductService предоставляет бизнес-логику для товаров.
type ProductService struct {
	repo *repository.ProductRepository
	log  *logger.Logger
}

// NewProductService создаёт новый сервис для товаров.
func NewProductService(repo *repository.ProductRepository, log *logger.Logger) *ProductService {
	return &ProductService{repo: repo, log: log}
}

// validateProduct проверяет корректность полей продукта в зависимости от типа.
func (s *ProductService) validateProduct(product *models.Product) error {
	// Проверка общих полей
	if product.Name == "" {
		return fmt.Errorf("name is required")
	}
	if product.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if product.Image == "" {
		return fmt.Errorf("image is required")
	}
	if product.CategoryID <= 0 {
		return fmt.Errorf("category_id must be greater than 0")
	}

	// Проверка типа
	if product.Type != "yarn" && product.Type != "garment" {
		return fmt.Errorf("type must be either 'yarn' or 'garment'")
	}

	// Проверка специфичных полей для yarn
	if product.Type == "yarn" {
		if product.Composition == "" {
			return fmt.Errorf("composition is required for yarn products")
		}
		if product.CountryOfOrigin == "" {
			return fmt.Errorf("country_of_origin is required for yarn products")
		}
		if product.LengthIn100g <= 0 {
			return fmt.Errorf("length_in_100g must be greater than 0 for yarn products")
		}
		if product.Color == "" {
			return fmt.Errorf("color is required for yarn products")
		}
	}

	// Проверка специфичных полей для garment
	if product.Type == "garment" {
		if product.Composition == "" {
			return fmt.Errorf("composition is required for garment products")
		}
		if product.Size == "" {
			return fmt.Errorf("size is required for garment products")
		}
		if product.GarmentLength == "" {
			return fmt.Errorf("garment_length is required for garment products")
		}
		if product.Color == "" {
			return fmt.Errorf("color is required for garment products")
		}
	}

	return nil
}

// CreateProduct создаёт новый товар.
func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Attempting to create product with name: %s, type: %s", product.Name, product.Type)

	// Валидация продукта
	if err := s.validateProduct(product); err != nil {
		s.log.Errorf("Validation failed for product '%s': %v", product.Name, err)
		return fmt.Errorf("validation failed: %w", err)
	}

	err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		s.log.Errorf("Failed to create product '%s': %v", product.Name, err)
		return fmt.Errorf("failed to create product: %w", err)
	}

	s.log.Infof("Successfully created product with ID: %d", product.ID)
	return nil
}

// GetProduct возвращает товар по ID.
func (s *ProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	s.log.Infof("Fetching product with ID: %d", id)

	product, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Product with ID %d not found", id)
			return nil, fmt.Errorf("product not found: %w", err)
		}
		s.log.Errorf("Failed to fetch product with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	s.log.Infof("Fetched product with ID: %d, Name: %s, Type: %s", product.ID, product.Name, product.Type)
	return product, nil
}

// ListProducts возвращает все товары.
func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	s.log.Info("Fetching all products")

	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch all products: %v", err)
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}

	s.log.Infof("Successfully fetched %d products", len(products))
	return products, nil
}

// ListProductsByCategory возвращает товары по ID категории.
func (s *ProductService) ListProductsByCategory(ctx context.Context, categoryID int) ([]*models.Product, error) {
	s.log.Infof("Fetching products for category ID: %d", categoryID)

	products, err := s.repo.ListProductsByCategory(ctx, categoryID)
	if err != nil {
		s.log.Errorf("Failed to fetch products for category ID %d: %v", categoryID, err)
		return nil, fmt.Errorf("failed to fetch products by category: %w", err)
	}

	s.log.Infof("Successfully fetched %d products for category ID: %d", len(products), categoryID)
	return products, nil
}

// UpdateProduct обновляет существующий товар.
func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	s.log.Infof("Updating product with ID: %d, Type: %s", product.ID, product.Type)

	// Валидация продукта
	if err := s.validateProduct(product); err != nil {
		s.log.Errorf("Validation failed for product ID %d: %v", product.ID, err)
		return fmt.Errorf("validation failed: %w", err)
	}

	err := s.repo.UpdateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to update product with ID %d: product not found", product.ID)
			return fmt.Errorf("product with ID %d not found: %w", product.ID, err)
		}
		s.log.Errorf("Failed to update product with ID %d: %v", product.ID, err)
		return fmt.Errorf("failed to update product: %w", err)
	}

	s.log.Infof("Successfully updated product with ID: %d", product.ID)
	return nil
}

// DeleteProduct удаляет товар по ID.
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	s.log.Infof("Deleting product with ID: %d", id)

	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to delete product with ID %d: product not found", id)
			return fmt.Errorf("product with ID %d not found: %w", id, err)
		}
		s.log.Errorf("Failed to delete product with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	s.log.Infof("Successfully deleted product with ID: %d", id)
	return nil
}
