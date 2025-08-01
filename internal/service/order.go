package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

// OrderService предоставляет бизнес-логику для заказов
type OrderService struct {
	repo *repository.OrderRepository
	log  *logger.Logger // Добавлено поле для логирования
}

// NewOrderService создаёт новый сервис для заказов
func NewOrderService(repo *repository.OrderRepository, log *logger.Logger) *OrderService {
	return &OrderService{repo: repo, log: log} // Инициализация логгера
}

// CreateOrder создаёт новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	s.log.Infof("Creating order for user ID: %d", order.UserID) // Логирование
	err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		s.log.Errorf("Failed to create order for user ID %d: %v", order.UserID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Order created with ID: %d for user ID: %d", order.ID, order.UserID) // Логирование успеха
	return nil
}

// GetOrder возвращает заказ по ID
func (s *OrderService) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	s.log.Infof("Fetching order with ID: %d", id) // Логирование
	order, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch order with ID %d: %v", id, err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched order with ID: %d, User ID: %d", order.ID, order.UserID) // Логирование успеха
	return order, err
}

// ListOrders возвращает все заказы
func (s *OrderService) ListOrders(ctx context.Context) ([]*models.Order, error) {
	s.log.Info("Fetching all orders") // Логирование
	orders, err := s.repo.ListOrders(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch orders: %v", err) // Логирование ошибки
		return nil, err
	}
	s.log.Infof("Fetched %d orders", len(orders)) // Логирование успеха
	return orders, err
}

// UpdateOrder обновляет существующий заказ
func (s *OrderService) UpdateOrder(ctx context.Context, order *models.Order) error {
	s.log.Infof("Updating order with ID: %d", order.ID) // Логирование
	err := s.repo.UpdateOrder(ctx, order)
	if err != nil {
		s.log.Errorf("Failed to update order with ID %d: %v", order.ID, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Order updated with ID: %d", order.ID) // Логирование успеха
	return err
}

// DeleteOrder удаляет заказ по ID
func (s *OrderService) DeleteOrder(ctx context.Context, id int) error {
	s.log.Infof("Deleting order with ID: %d", id) // Логирование
	err := s.repo.DeleteOrder(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete order with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	s.log.Infof("Order deleted with ID: %d", id) // Логирование успеха
	return err
}
