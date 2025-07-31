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

type CommentRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewCommentRepository(db *sql.DB, redis *redis.Client) *CommentRepository {
	return &CommentRepository{db: db, redis: redis}
}

func (r *CommentRepository) CreateComment(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (product_id, user_id, text, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, time.Now()).Scan(&comment.ID)
	return err
}

func (r *CommentRepository) GetComment(ctx context.Context, id int) (*models.Comment, error) {
	cacheKey := "comment:" + strconv.Itoa(id)
	var comment models.Comment

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &comment); err == nil {
			return &comment, nil
		}
	}

	query := `SELECT id, product_id, user_id, text, created_at FROM comments WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.ProductID, &comment.UserID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(comment)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return &comment, nil
}

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
	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, comment *models.Comment) error {
	query := `UPDATE comments SET product_id = $1, user_id = $2, text = $3 WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, comment.ProductID, comment.UserID, comment.Text, comment.ID)
	return err
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
