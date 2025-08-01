package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/alex-pyslar/online-store/internal/logger"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/alex-pyslar/online-store/internal/models"
)

// UserRepository управляет доступом к данным пользователей в базе данных и кэше.
type UserRepository struct {
	db    *sql.DB
	redis *redis.Client
	log   *logger.Logger
}

// NewUserRepository создаёт новый репозиторий для пользователей.
func NewUserRepository(db *sql.DB, redis *redis.Client, log *logger.Logger) *UserRepository {
	return &UserRepository{db: db, redis: redis, log: log}
}

// CreateUser создаёт нового пользователя в базе данных.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	r.log.Infof("Creating user in DB with email: %s", user.Email)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			r.log.Errorf("Failed to hash password for email %s: %v", user.Email, err)
			return err
		}
		user.Password = string(hashedPassword)
	}

	query := `INSERT INTO users (email, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Password, time.Now(), time.Now()).Scan(&user.ID)
	if err != nil {
		r.log.Errorf("Failed to insert user into DB with email %s: %v", user.Email, err)
		return err
	}
	r.log.Infof("User created in DB with ID: %d", user.ID)
	return nil
}

// GetUser получает пользователя по ID, используя кэш Redis.
func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	cacheKey := "user:" + strconv.Itoa(id)
	r.log.Infof("Fetching user with ID %d from cache", id)
	var user models.User

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			r.log.Infof("User with ID %d found in cache", id)
			return &user, nil
		}
		r.log.Warningf("Failed to unmarshal user from cache for ID %d: %v", id, err)
	} else {
		r.log.Infof("User with ID %d not found in cache: %v", id, err)
	}

	query := `SELECT id, email, name, password, created_at, updated_at FROM users WHERE id = $1`
	err = r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch user with ID %d from DB: %v", id, err)
		return nil, err
	}
	r.log.Infof("User with ID %d fetched from DB", id)

	data, _ := json.Marshal(user)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache user with ID %d: %v", id, err)
	} else {
		r.log.Infof("User with ID %d cached in Redis", id)
	}
	return &user, nil
}

// GetUserByEmail получает пользователя по email, используя кэш Redis.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	cacheKey := "user:email:" + email
	r.log.Infof("Fetching user with email %s from cache", email)
	var user models.User

	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			r.log.Infof("User with email %s found in cache", email)
			return &user, nil
		}
		r.log.Warningf("Failed to unmarshal user from cache for email %s: %v", email, err)
	} else {
		r.log.Infof("User with email %s not found in cache: %v", email, err)
	}

	query := `SELECT id, email, name, password, created_at, updated_at FROM users WHERE email = $1`
	err = r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.log.Errorf("Failed to fetch user with email %s from DB: %v", email, err)
		return nil, err
	}
	r.log.Infof("User with email %s fetched from DB", email)

	data, _ := json.Marshal(user)
	if err := r.redis.Set(ctx, cacheKey, data, 10*time.Minute).Err(); err != nil {
		r.log.Warningf("Failed to cache user with email %s: %v", email, err)
	} else {
		r.log.Infof("User with email %s cached in Redis", email)
	}
	return &user, nil
}

// ListUsers получает список всех пользователей.
func (r *UserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	r.log.Info("Fetching all users from DB")
	query := `SELECT id, email, name, password, created_at, updated_at FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Errorf("Failed to fetch users from DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Password, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			r.log.Errorf("Failed to scan user row: %v", err)
			return nil, err
		}
		users = append(users, &u)
	}
	r.log.Infof("Fetched %d users from DB", len(users))
	return users, nil
}

// UpdateUser обновляет существующего пользователя.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	r.log.Infof("Updating user with ID %d in DB", user.ID)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			r.log.Errorf("Failed to hash password for user ID %d: %v", user.ID, err)
			return err
		}
		user.Password = string(hashedPassword)
	}

	query := `UPDATE users SET email = $1, name = $2, password = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Name, user.Password, time.Now(), user.ID)
	if err != nil {
		r.log.Errorf("Failed to update user with ID %d: %v", user.ID, err)
		return err
	}
	r.log.Infof("User updated with ID %d", user.ID)
	return nil
}

// DeleteUser удаляет пользователя по ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	r.log.Infof("Deleting user with ID %d from DB", id)
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.Errorf("Failed to delete user with ID %d: %v", id, err)
		return err
	}
	r.log.Infof("User deleted with ID %d", id)
	return nil
}

// VerifyPassword проверяет пароль пользователя.
func (r *UserRepository) VerifyPassword(ctx context.Context, id int, password string) error {
	r.log.Infof("Verifying password for user ID %d", id)
	var hashedPassword string
	query := `SELECT password FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&hashedPassword)
	if err != nil {
		r.log.Errorf("Failed to fetch password for user ID %d: %v", id, err)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		r.log.Errorf("Password verification failed for user ID %d: %v", id, err)
		return err
	}
	r.log.Infof("Password verified successfully for user ID %d", id)
	return nil
}
