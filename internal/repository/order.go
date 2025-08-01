package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/alex-pyslar/online-store/internal/logger" // Импортируем логгер
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
)

// OrderRepository управляет доступом к данным заказов в базе данных и кэше.
type OrderRepository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logger.Logger // Добавлено поле для логирования
}

// NewOrderRepository создаёт новый репозиторий для заказов.
func NewOrderRepository(db *sql.DB, redis *redis.Client, log *logger.Logger) *OrderRepository {
	return &OrderRepository{db: db, redis: redis, log: log} // Инициализация логгера
}

// CreateOrder создаёт новый заказ в базе данных.
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	r.log.Infof("Creating order in DB for user ID: %d", order.UserID) // Логирование
	query := `INSERT INTO orders (user_id, total, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, order.UserID, order.Total, order.Status, time.Now()).Scan(&order.ID)
	if err != nil {
		r.log.Errorf("Failed to insert order into DB for user ID %d: %v", order.UserID, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Order created in DB with ID: %d", order.ID) // Логирование успеха
	return nil
}

// GetOrder получает заказ по ID, используя кэш Redis.
func (r *OrderRepository) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	cacheKey := "order:" + strconv.Itoa(id)
	r.log.Infof("Fetching order with ID %d from cache", id) // Логирование
	var order models.Order

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			r.log.Infof("Order with ID %d found in cache", id) // Логирование
			return &order, nil
		}
		r.log.Warningf("Failed to unmarshal order from cache for ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Order with ID %d not found in cache: %v", id, err) // Логирование
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, user_id, total, status, created_at FROM orders WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch order with ID %d from DB: %v", id, err) // Логирование ошибки
		return nil, err
	}
	r.log.Infof("Order with ID %d fetched from DB", id) // Логирование успеха

	// Сохраняем в кэш
	data, _ := json.Marshal(order)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache order with ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Order with ID %d cached in Redis", id) // Логирование
	}
	return &order, nil
}

// ListOrders получает список всех заказов.
func (r *OrderRepository) ListOrders(ctx context.Context) ([]*models.Order, error) {
	r.log.Info("Fetching all orders from DB") // Логирование
	query := `SELECT id, user_id, total, status, created_at FROM orders`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("Failed to fetch orders from DB: %v", err) // Логирование ошибки
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.Status, &o.CreatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan order row: %v", err) // Логирование ошибки
			return nil, err
		}
		orders = append(orders, &o)
	}
	r.log.Infof("Fetched %d orders from DB", len(orders)) // Логирование успеха
	return orders, nil
}

// UpdateOrder обновляет существующий заказ.
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	r.log.Infof("Updating order with ID %d in DB", order.ID) // Логирование
	query := `UPDATE orders SET user_id = $1, total = $2, status = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Total, order.Status, order.ID)
	if err != nil {
		r.log.Errorf("Failed to update order with ID %d: %v", order.ID, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Order updated with ID %d", order.ID) // Логирование успеха
	return nil
}

// DeleteOrder удаляет заказ по ID.
func (r *OrderRepository) DeleteOrder(ctx context.Context, id int) error {
	r.log.Infof("Deleting order with ID %d from DB", id) // Логирование
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("Failed to delete order with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Order deleted with ID %d", id) // Логирование успеха
	return nil
}
