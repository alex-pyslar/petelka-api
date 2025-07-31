package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
)

type OrderRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewOrderRepository(db *sql.DB, redis *redis.Client) *OrderRepository {
	return &OrderRepository{db: db, redis: redis}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `INSERT INTO orders (user_id, total, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, order.UserID, order.Total, order.Status, time.Now()).Scan(&order.ID)
	return err
}

func (r *OrderRepository) GetOrder(ctx context.Context, id int) (*models.Order, error) {
	cacheKey := "order:" + strconv.Itoa(id)
	var order models.Order

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			return &order, nil
		}
	}

	query := `SELECT id, user_id, total, status, created_at FROM orders WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(order)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return &order, nil
}

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
		err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	query := `UPDATE orders SET user_id = $1, total = $2, status = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Total, order.Status, order.ID)
	return err
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id int) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
