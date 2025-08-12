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
	query := `INSERT INTO products (name, description, price, category_id, image, type, composition, country_of_origin, length_in_100g, size, garment_length, color) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryID,
		product.Image,
		product.Type,
		product.Composition,
		product.CountryOfOrigin,
		product.LengthIn100g,
		product.Size,
		product.GarmentLength,
		product.Color,
	).Scan(&product.ID)
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
	query := `SELECT id, name, description, price, category_id, image, type, composition, country_of_origin, length_in_100g, size, garment_length, color 
	          FROM products WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CategoryID,
		&product.Image,
		&product.Type,
		&product.Composition,
		&product.CountryOfOrigin,
		&product.LengthIn100g,
		&product.Size,
		&product.GarmentLength,
		&product.Color,
	)
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
	query := `SELECT id, name, description, price, category_id, image, type, composition, country_of_origin, length_in_100g, size, garment_length, color 
	          FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CategoryID,
			&p.Image,
			&p.Type,
			&p.Composition,
			&p.CountryOfOrigin,
			&p.LengthIn100g,
			&p.Size,
			&p.GarmentLength,
			&p.Color,
		); err != nil {
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
	query := `SELECT id, name, description, price, category_id, image, type, composition, country_of_origin, length_in_100g, size, garment_length, color 
	          FROM products WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CategoryID,
			&p.Image,
			&p.Type,
			&p.Composition,
			&p.CountryOfOrigin,
			&p.LengthIn100g,
			&p.Size,
			&p.GarmentLength,
			&p.Color,
		); err != nil {
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
	query := `UPDATE products SET name = $1, description = $2, price = $3, category_id = $4, image = $5, type = $6, 
	          composition = $7, country_of_origin = $8, length_in_100g = $9, size = $10, garment_length = $11, color = $12 
	          WHERE id = $13`
	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.CategoryID,
		product.Image,
		product.Type,
		product.Composition,
		product.CountryOfOrigin,
		product.LengthIn100g,
		product.Size,
		product.GarmentLength,
		product.Color,
		product.ID,
	)
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
