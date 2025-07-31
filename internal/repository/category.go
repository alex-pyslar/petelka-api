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

type CategoryRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewCategoryRepository(db *sql.DB, redis *redis.Client) *CategoryRepository {
	return &CategoryRepository{db: db, redis: redis}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO categories (name, created_at) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name, time.Now()).Scan(&category.ID)
	return err
}

func (r *CategoryRepository) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	cacheKey := "category:" + strconv.Itoa(id)
	var category models.Category

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &category); err == nil {
			return &category, nil
		}
	}

	query := `SELECT id, name, created_at FROM categories WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(category)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return &category, nil
}

func (r *CategoryRepository) ListCategories(ctx context.Context) ([]*models.Category, error) {
	query := `SELECT id, name, created_at FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	return err
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
