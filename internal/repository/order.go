package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/redis/go-redis/v9"
)

// OrderRepository управляет доступом к данным заказов в базе данных и кэше.
type OrderRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewOrderRepository создаёт новый репозиторий для заказов.
func NewOrderRepository(db *sql.DB, redis *redis.Client) *OrderRepository {
	return &OrderRepository{db: db, redis: redis}
}

// CreateOrder создаёт новый заказ в базе данных.
func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `INSERT INTO orders (user_id, total, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, order.UserID, order.Total, order.Status, time.Now()).Scan(&order.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetOrder получает заказ по ID, используя кэш Redis.
func (r *OrderRepository) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	cacheKey := fmt.Sprintf("order:%d", id)
	var order models.Order

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return &order, nil
		}
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, user_id, total, status, created_at FROM orders WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Сохраняем в кэш
	data, err := json.Marshal(order)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	return &order, nil
}

// ListOrders получает список всех заказов.
func (r *OrderRepository) ListOrders(ctx context.Context) ([]*models.Order, error) {
	query := `SELECT id, user_id, total, status, created_at FROM orders`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrder обновляет существующий заказ.
func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := `UPDATE orders SET user_id = $1, total = $2, status = $3 WHERE id = $4`
	result, err := r.db.ExecContext(ctx, query, order.UserID, order.Total, order.Status, order.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("order:%d", order.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteOrder удаляет заказ по ID.
func (r *OrderRepository) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("order:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}
