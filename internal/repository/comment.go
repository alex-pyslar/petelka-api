package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/redis/go-redis/v9"
)

// CommentRepository управляет доступом к данным комментариев в базе данных и кэше.
type CommentRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewCommentRepository создаёт новый репозиторий для комментариев.
func NewCommentRepository(db *sql.DB, redis *redis.Client) *CommentRepository {
	return &CommentRepository{db: db, redis: redis}
}

// CreateComment создаёт новый комментарий в базе данных.
func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (product_id, user_id, text, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, time.Now()).Scan(&comment.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetComment получает комментарий по ID, используя кэш Redis.
func (r *CommentRepository) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	cacheKey := fmt.Sprintf("comment:%d", id)
	var comment models.Comment

	// Попытка получить данные из кэша
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &comment); err == nil {
			return &comment, nil
		}
	}

	// Если в кэше нет, получаем из БД
	query := `SELECT id, product_id, user_id, text, created_at FROM comments WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.ProductID, &comment.UserID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Сохраняем в кэш
	data, err := json.Marshal(comment)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	return &comment, nil
}

// ListComments получает список всех комментариев.
func (r *CommentRepository) ListComments(ctx context.Context) ([]*models.Comment, error) {
	query := `SELECT id, product_id, user_id, text, created_at FROM comments`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Text, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateComment обновляет существующий комментарий.
func (r *CommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `UPDATE comments SET product_id = $1, user_id = $2, text = $3 WHERE id = $4`
	result, err := r.db.ExecContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, comment.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("comment:%d", comment.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteComment удаляет комментарий по ID.
func (r *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("comment:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}
