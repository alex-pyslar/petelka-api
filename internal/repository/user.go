package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/alex-pyslar/online-store/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	query := `INSERT INTO users (email, name, oauth_id, password, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.OAuthID, user.Password, time.Now()).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	cacheKey := "user:" + strconv.Itoa(id)
	var user models.User

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	query := `SELECT id, email, name, oauth_id, password, created_at FROM users WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name, &user.OAuthID, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(user)
	r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, email, name, oauth_id, password, created_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.OAuthID, &u.Password, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	query := `UPDATE users SET email = $1, name = $2, oauth_id = $3, password = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Name, user.OAuthID, user.Password, user.ID)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) VerifyPassword(ctx context.Context, id int, password string) error {
	user, err := r.GetUser(ctx, id)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
