package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.CreateProduct(ctx, product)
}

func (s *ProductService) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context) ([]*models.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) error {
	return s.repo.UpdateProduct(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.DeleteProduct(ctx, id)
}
