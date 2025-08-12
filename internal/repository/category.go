package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/redis/go-redis/v9"
)

// CategoryRepository управляет доступом к данным категорий в базе данных и кэше.
type CategoryRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewCategoryRepository создаёт новый репозиторий для категорий.
func NewCategoryRepository(db *sql.DB, redis *redis.Client) *CategoryRepository {
	return &CategoryRepository{db: db, redis: redis}
}

// CreateCategory создаёт новую категорию в базе данных.
func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name).Scan(&category.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetCategory получает категорию по ID, используя кэш Redis.
func (r *CategoryRepository) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	cacheKey := fmt.Sprintf("category:%d", id)

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var category models.Category
		if err := json.Unmarshal([]byte(cached), &category); err == nil {
			return &category, nil
		}
	}

	var category models.Category
	query := `SELECT id, name FROM categories WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	data, err := json.Marshal(category)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}

	return &category, nil
}

// ListCategories получает список всех категорий.
func (r *CategoryRepository) ListCategories(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory обновляет существующую категорию.
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("category:%d", category.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteCategory удаляет категорию по ID.
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("category:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}
