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

type ProductRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewProductRepository(db *sql.DB, redis *redis.Client) *ProductRepository {
	return &ProductRepository{db: db, redis: redis}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (name, description, price, category_id, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, time.Now()).Scan(&product.ID)
	return err
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	cacheKey := "product:" + strconv.Itoa(id)
	var product models.Product

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	query := `SELECT id, name, description, price, category_id, created_at FROM products WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(product)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return &product, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]*models.Product, error) {
	query := `SELECT id, name, description, price, category_id, created_at FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, product.ID)
	return err
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
