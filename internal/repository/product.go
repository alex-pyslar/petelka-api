package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
)

// ProductRepository управляет доступом к данным товаров в базе данных и кэше.
type ProductRepository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logger.Logger
}

// NewProductRepository создаёт новый репозиторий для товаров.
func NewProductRepository(db *sql.DB, redis *redis.Client, log *logger.Logger) *ProductRepository {
	return &ProductRepository{db: db, redis: redis, log: log}
}

// CreateProduct создаёт новый товар в базе данных.
func (r *ProductRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	r.log.Infof("Creating product in DB with name: %s", product.Name)
	query := `INSERT INTO products (name, description, price, category_id, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, time.Now()).Scan(&product.ID)
	if err != nil {
		r.log.Errorf("Failed to insert product into DB with name %s: %v", product.Name, err)
		return err
	}
	r.log.Infof("Product created in DB with ID: %d", product.ID)
	return nil
}

// GetProduct получает товар по ID, используя кэш Redis.
func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	cacheKey := "product:" + strconv.Itoa(id)
	r.log.Infof("Fetching product with ID %d from cache", id)
	var product models.Product

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			r.log.Infof("Product with ID %d found in cache", id)
			return &product, nil
		}
		r.log.Warningf("Failed to unmarshal product from cache for ID %d: %v", id, err)
	} else {
		r.log.Infof("Product with ID %d not found in cache: %v", id, err)
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, name, description, price, category_id, created_at FROM products WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.CreatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch product with ID %d from DB: %v", id, err)
		return nil, err
	}
	r.log.Infof("Product with ID %d fetched from DB", id)

	// Сохраняем в кэш
	data, _ := json.Marshal(product)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache product with ID %d: %v", id, err)
	} else {
		r.log.Infof("Product with ID %d cached in Redis", id)
	}
	return &product, nil
}

// ListProducts получает список всех товаров.
func (r *ProductRepository) ListProducts(ctx context.Context) ([]*models.Product, error) {
	r.log.Info("Fetching all products from DB")
	query := `SELECT id, name, description, price, category_id, created_at FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("Failed to fetch products from DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan product row: %v", err)
			return nil, err
		}
		products = append(products, &p)
	}
	r.log.Infof("Fetched %d products from DB", len(products))
	return products, nil
}

// ListProductsByCategory получает список товаров по ID категории.
func (r *ProductRepository) ListProductsByCategory(ctx context.Context, categoryID int) ([]*models.Product, error) {
	r.log.Infof("Fetching products for category ID %d from DB", categoryID)
	query := `SELECT id, name, description, price, category_id, created_at FROM products WHERE category_id = $1`
	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		r.log.Errorf("Failed to fetch products for category ID %d from DB: %v", categoryID, err)
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan product row for category ID %d: %v", categoryID, err)
			return nil, err
		}
		products = append(products, &p)
	}
	r.log.Infof("Fetched %d products for category ID %d from DB", len(products), categoryID)
	return products, nil
}

// UpdateProduct обновляет существующий товар.
func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	r.log.Infof("Updating product with ID %d in DB", product.ID)
	query := `UPDATE products SET name = $1, description = $2, price = $3, category_id = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, product.ID)
	if err != nil {
		r.log.Errorf("Failed to update product with ID %d: %v", product.ID, err)
		return err
	}
	r.log.Infof("Product updated with ID %d", product.ID)
	return nil
}

// DeleteProduct удаляет товар по ID.
func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
	r.log.Infof("Deleting product with ID %d from DB", id)
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("Failed to delete product with ID %d: %v", id, err)
		return err
	}
	r.log.Infof("Product deleted with ID %d", id)
	return nil
}
