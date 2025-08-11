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

// ProductRepository управляет доступом к данным товаров в базе данных и кэше.
type ProductRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewProductRepository создаёт новый репозиторий для товаров.
func NewProductRepository(db *sql.DB, redis *redis.Client) *ProductRepository {
	return &ProductRepository{db: db, redis: redis}
}

// CreateProduct создаёт новый товар в базе данных.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	query := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, time.Now()).Scan(&product.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetProduct получает товар по ID, используя кэш Redis.
func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	var product models.Product

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, name, description, price FROM products WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Сохраняем в кэш
	data, err := json.Marshal(product)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	return &product, nil
}

// ListProducts получает список всех товаров.
func (r *ProductRepository) ListProducts(ctx context.Context) ([]*models.Product, error) {
	query := `SELECT id, name, description, price, category_id FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// ListProductsByCategory получает список товаров по ID категории.
func (r *ProductRepository) ListProductsByCategory(ctx context.Context, categoryID int) ([]*models.Product, error) {
	query := `SELECT id, name, description, price, category_id FROM products WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct обновляет существующий товар.
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("product:%d", product.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteProduct удаляет товар по ID.
func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}
