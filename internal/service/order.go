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

// OrderService предоставляет бизнес-логику для заказов
type OrderService struct {
	repo *repository.OrderRepository
	log  *logger.Logger
}

// NewOrderService создаёт новый сервис для заказов
func NewOrderService(repo *repository.OrderRepository, log *logger.Logger) *OrderService {
	return &OrderService{repo: repo, log: log}
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	s.log.Infof("Attempting to create order for user ID: %d", order.UserID)

	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		s.log.Errorf("Failed to create order for user ID %d: %v", order.UserID, err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	s.log.Infof("Successfully created order with ID: %d", order.ID)
	return nil
}

// GetOrder возвращает заказ по ID
func (s *OrderService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	s.log.Infof("Fetching order with ID: %d", id)

	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Order with ID %d not found", id)
			return nil, fmt.Errorf("order not found: %w", err)
		}
		s.log.Errorf("Failed to fetch order with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	s.log.Infof("Fetched order with ID: %d, User ID: %d", order.ID, order.UserID)
	return order, nil
}

// ListOrders возвращает все заказы
func (s *OrderService) ListOrders(ctx context.Context) ([]*models.Order, error) {
	s.log.Info("Fetching all orders from repository")

	orders, err := s.repo.ListOrders(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch all orders: %v", err)
		return nil, fmt.Errorf("failed to fetch all orders: %w", err)
	}

	s.log.Infof("Successfully fetched %d orders", len(orders))
	return orders, nil
}

// UpdateOrder обновляет существующий заказ
func (s *OrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	s.log.Infof("Updating order with ID: %d", order.ID)

	err := s.repo.UpdateOrder(ctx, order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to update order with ID %d: order not found", order.ID)
			return fmt.Errorf("order with ID %d not found: %w", order.ID, err)
		}
		s.log.Errorf("Failed to update order with ID %d: %v", order.ID, err)
		return fmt.Errorf("failed to update order: %w", err)
	}

	s.log.Infof("Successfully updated order with ID: %d", order.ID)
	return nil
}

// DeleteOrder удаляет заказ по ID
func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	s.log.Infof("Deleting order with ID: %d", id)

	err := s.repo.DeleteOrder(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to delete order with ID %d: order not found", id)
			return fmt.Errorf("order with ID %d not found: %w", id, err)
		}
		s.log.Errorf("Failed to delete order with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete order: %w", err)
	}

	s.log.Infof("Successfully deleted order with ID: %d", id)
	return nil
}
