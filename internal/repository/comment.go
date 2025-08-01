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

// CommentRepository управляет доступом к данным комментариев в базе данных и кэше.
type CommentRepository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logger.Logger // Добавлено поле для логирования
}

// NewCommentRepository создаёт новый репозиторий для комментариев.
func NewCommentRepository(db *sql.DB, redis *redis.Client, log *logger.Logger) *CommentRepository {
	return &CommentRepository{db: db, redis: redis, log: log} // Инициализация логгера
}

// CreateComment создаёт новый комментарий в базе данных.
func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	r.log.Infof("Creating comment in DB for product ID: %d by user ID: %d", comment.ProductID, comment.UserID) // Логирование
	query := `INSERT INTO comments (product_id, user_id, text, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, time.Now()).Scan(&comment.ID)
	if err != nil {
		r.log.Errorf("Failed to insert comment into DB for product ID %d by user ID %d: %v", comment.ProductID, comment.UserID, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Comment created in DB with ID: %d", comment.ID) // Логирование успеха
	return nil
}

// GetComment получает комментарий по ID, используя кэш Redis.
func (r *CommentRepository) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	cacheKey := "comment:" + strconv.Itoa(id)
	r.log.Infof("Fetching comment with ID %d from cache", id) // Логирование
	var comment models.Comment

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &comment); err == nil {
			r.log.Infof("Comment with ID %d found in cache", id) // Логирование
			return &comment, nil
		}
		r.log.Warningf("Failed to unmarshal comment from cache for ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Comment with ID %d not found in cache: %v", id, err) // Логирование
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, product_id, user_id, text, created_at FROM comments WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.ProductID, &comment.UserID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch comment with ID %d from DB: %v", id, err) // Логирование ошибки
		return nil, err
	}
	r.log.Infof("Comment with ID %d fetched from DB", id) // Логирование успеха

	// Сохраняем в кэш
	data, _ := json.Marshal(comment)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache comment with ID %d: %v", id, err) // Логирование предупреждения
	} else {
		r.log.Infof("Comment with ID %d cached in Redis", id) // Логирование
	}
	return &comment, nil
}

// ListComments получает список всех комментариев.
func (r *CommentRepository) ListComments(ctx context.Context) ([]*models.Comment, error) {
	r.log.Info("Fetching all comments from DB") // Логирование
	query := `SELECT id, product_id, user_id, text, created_at FROM comments`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("Failed to fetch comments from DB: %v", err) // Логирование ошибки
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Text, &c.CreatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan comment row: %v", err) // Логирование ошибки
			return nil, err
		}
		comments = append(comments, &c)
	}
	r.log.Infof("Fetched %d comments from DB", len(comments)) // Логирование успеха
	return comments, nil
}

// UpdateComment обновляет существующий комментарий.
func (r *CommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	r.log.Infof("Updating comment with ID %d in DB", comment.ID) // Логирование
	query := `UPDATE comments SET product_id = $1, user_id = $2, text = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, comment.ID)
	if err != nil {
		r.log.Errorf("Failed to update comment with ID %d: %v", comment.ID, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Comment updated with ID %d", comment.ID) // Логирование успеха
	return nil
}

// DeleteComment удаляет комментарий по ID.
func (r *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	r.log.Infof("Deleting comment with ID %d from DB", id) // Логирование
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("Failed to delete comment with ID %d: %v", id, err) // Логирование ошибки
		return err
	}
	r.log.Infof("Comment deleted with ID %d", id) // Логирование успеха
	return nil
}
