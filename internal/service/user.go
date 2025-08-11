package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserService предоставляет бизнес-логику для пользователей
type UserService struct {
	repo *repository.UserRepository
	log  *logger.Logger
}

// NewUserService создаёт новый сервис для пользователей
func NewUserService(repo *repository.UserRepository, log *logger.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

// CreateUser создаёт нового пользователя
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	s.log.Infof("Attempting to create user with email: %s", user.Email)

	if user.Password == "" {
		s.log.Errorf("Password is required for user with email: %s", user.Email)
		return fmt.Errorf("password is required")
	}

	// Хешируем пароль перед сохранением в репозиторий
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Errorf("Failed to hash password for email %s: %v", user.Email, err)
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to create user with email %s: %v", user.Email, err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	s.log.Infof("Successfully created user with ID: %d", user.ID)
	return nil
}

// GetUser возвращает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	s.log.Infof("Fetching user with ID: %d", id)

	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("User with ID %d not found", id)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		s.log.Errorf("Failed to fetch user with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	s.log.Infof("Fetched user with ID: %d, Email: %s", user.ID, user.Email)
	return user, nil
}

// GetUserByEmail возвращает пользователя по email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	s.log.Infof("Fetching user with email: %s", email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("User with email %s not found", email)
			return nil, fmt.Errorf("user not found: %w", err)
		}
		s.log.Errorf("Failed to fetch user with email %s: %v", email, err)
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	s.log.Infof("Fetched user with email: %s, ID: %d", user.Email, user.ID)
	return user, nil
}

// ListUsers возвращает всех пользователей
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	s.log.Info("Fetching all users from repository")

	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch all users: %v", err)
		return nil, fmt.Errorf("failed to fetch all users: %w", err)
	}

	s.log.Infof("Successfully fetched %d users", len(users))
	return users, nil
}

// UpdateUser обновляет существующего пользователя
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	s.log.Infof("Updating user with ID: %d", user.ID)

	// Хешируем пароль, если он был передан для обновления
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			s.log.Errorf("Failed to hash password for user ID %d: %v", user.ID, err)
			return fmt.Errorf("failed to hash new password: %w", err)
		}
		user.Password = string(hashedPassword)
	}

	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to update user with ID %d: user not found", user.ID)
			return fmt.Errorf("user with ID %d not found: %w", user.ID, err)
		}
		s.log.Errorf("Failed to update user with ID %d: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	s.log.Infof("Successfully updated user with ID: %d", user.ID)
	return nil
}

// DeleteUser удаляет пользователя по ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	s.log.Infof("Deleting user with ID: %d", id)

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Failed to delete user with ID %d: user not found", id)
			return fmt.Errorf("user with ID %d not found: %w", id, err)
		}
		s.log.Errorf("Failed to delete user with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.log.Infof("Successfully deleted user with ID: %d", id)
	return nil
}

// VerifyPassword проверяет пароль пользователя
func (s *UserService) VerifyPassword(ctx context.Context, id int, password string) error {
	s.log.Infof("Verifying password for user ID: %d", id)

	hashedPassword, err := s.repo.GetUserPassword(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warningf("Password verification failed for user ID %d: user not found", id)
			return fmt.Errorf("user not found")
		}
		s.log.Errorf("Failed to retrieve password for user ID %d: %v", id, err)
		return fmt.Errorf("failed to retrieve user password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		s.log.Warningf("Password mismatch for user ID %d: %v", id, err)
		return fmt.Errorf("invalid password")
	}

	s.log.Infof("Password verified successfully for user ID %d", id)
	return nil
}
