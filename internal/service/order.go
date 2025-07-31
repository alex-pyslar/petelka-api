package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	return s.repo.CreateOrder(ctx, order)
}

func (s *OrderService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) ([]*models.Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *OrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	return s.repo.UpdateOrder(ctx, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	return s.repo.DeleteOrder(ctx, id)
}
