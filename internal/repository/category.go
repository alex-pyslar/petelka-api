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

// CategoryRepository управляет доступом к данным категорий в базе данных и кэше.
type CategoryRepository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logger.Logger // Добавлено поле для логирования
}

// NewCategoryRepository создаёт новый репозиторий для категорий.
func NewCategoryRepository(db *sql.DB, redis *redis.Client, log *logger.Logger) *CategoryRepository {
	return &CategoryRepository{db: db, redis: redis, log: log} // Инициализация логгера
}

// CreateCategory создаёт новую категорию в базе данных.
func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	r.log.Infof("Creating category in DB with name: %s", category.Name) // Логирование
	query := `INSERT INTO categories (name, created_at) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, category.Name, time.Now()).Scan(&category.ID)
	if err != nil {
		r.log.Errorf("Failed to insert category into DB with name %s: %v", category.Name, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Category created in DB with ID: %d", category.ID) // Логирование успеха
	return nil
}

// GetCategory получает категорию по ID, используя кэш Redis.
func (r *CategoryRepository) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	cacheKey := "category:" + strconv.Itoa(id)
	r.log.Infof("Fetching category with ID %d from cache", id) // Логирование
	var category models.Category

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &category); err == nil {
			r.log.Infof("Category with ID %d found in cache", id) // Логирование
			return &category, nil
		}
		r.log.Warningf("Failed to unmarshal category from cache for ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Category with ID %d not found in cache: %v", id, err) // Логирование
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, name, created_at FROM categories WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch category with ID %d from DB: %v", id, err) // Логирование ошибки
		return nil, err
	}
	r.log.Infof("Category with ID %d fetched from DB", id) // Логирование успеха

	// Сохраняем в кэш
	data, _ := json.Marshal(category)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache category with ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Category with ID %d cached in Redis", id) // Логирование
	}
	return &category, nil
}

// ListCategories получает список всех категорий.
func (r *CategoryRepository) ListCategories(ctx context.Context) ([]*models.Category, error) {
	r.log.Info("Fetching all categories from DB") // Логирование
	query := `SELECT id, name, created_at FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("Failed to fetch categories from DB: %v", err) // Логирование ошибки
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan category row: %v", err) // Логирование ошибки
			return nil, err
		}
		categories = append(categories, &c)
	}
	r.log.Infof("Fetched %d categories from DB", len(categories)) // Логирование успеха
	return categories, nil
}

// UpdateCategory обновляет существующую категорию.
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	r.log.Infof("Updating category with ID %d in DB", category.ID) // Логирование
	query := `UPDATE categories SET name = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		r.log.Errorf("Failed to update category with ID %d: %v", category.ID, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Category updated with ID %d", category.ID) // Логирование успеха
	return nil
}

// DeleteCategory удаляет категорию по ID.
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	r.log.Infof("Deleting category with ID %d from DB", id) // Логирование
	query := `DELETE FROM categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("Failed to delete category with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Category deleted with ID %d", id) // Логирование успеха
	return nil
}
