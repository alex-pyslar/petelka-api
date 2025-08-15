package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/redis/go-redis/v9"
)

// UserRepository управляет доступом к данным пользователей в базе данных и кэше.
type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
}

// NewUserRepository создаёт новый репозиторий для пользователей.
func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{db: db, redis: redis}
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, name, password, created_at, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Password, time.Now(), user.Role).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser gets a user by ID, using Redis cache.
func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	var user models.User

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	query := `SELECT id, email, name, role, password, created_at FROM users WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	data, err := json.Marshal(user)
	if err == nil {
		r.redis.Set(ctx, cacheKey, data, 10*time.Minute)
	}
	return &user, nil
}

// GetUserByEmail gets a user by email.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, email, name, role, password, created_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers gets a list of all users.
func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, email, name, role, password, created_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates an existing user.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET email = $1, name = $2, role = $3, password = $4 WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query, user.Email, user.Name, user.Role, user.Password, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("user:%d", user.ID)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// DeleteUser удаляет пользователя по ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	cacheKey := fmt.Sprintf("user:%d", id)
	r.redis.Del(ctx, cacheKey)

	return nil
}

// GetUserPassword получает хешированный пароль пользователя по ID.
func (r *UserRepository) GetUserPassword(ctx context.Context, id int) (string, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}
	return hashedPassword, nil
}
