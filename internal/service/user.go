package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
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
	s.log.Infof("Creating user with email: %s", user.Email)
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to create user with email %s: %v", user.Email, err)
		return err
	}
	s.log.Infof("User created with email: %s, ID: %d", user.Email, user.ID)
	return nil
}

// GetUser возвращает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	s.log.Infof("Fetching user with ID: %d", id)
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to fetch user with ID %d: %v", id, err)
		return nil, err
	}
	s.log.Infof("Fetched user with ID: %d, Email: %s", user.ID, user.Email)
	return user, err
}

// GetUserByEmail возвращает пользователя по email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	s.log.Infof("Fetching user with email: %s", email)
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.log.Errorf("Failed to fetch user with email %s: %v", email, err)
		return nil, err
	}
	s.log.Infof("Fetched user with email: %s, ID: %d", user.Email, user.ID)
	return user, err
}

// ListUsers возвращает всех пользователей
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	s.log.Info("Fetching all users")
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		s.log.Errorf("Failed to fetch users: %v", err)
		return nil, err
	}
	s.log.Infof("Fetched %d users", len(users))
	return users, err
}

// UpdateUser обновляет существующего пользователя
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	s.log.Infof("Updating user with ID: %d", user.ID)
	err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		s.log.Errorf("Failed to update user with ID %d: %v", user.ID, err)
		return err
	}
	s.log.Infof("User updated with ID: %d", user.ID)
	return err
}

// DeleteUser удаляет пользователя по ID
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	s.log.Infof("Deleting user with ID: %d", id)
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		s.log.Errorf("Failed to delete user with ID %d: %v", id, err)
		return err
	}
	s.log.Infof("User deleted with ID: %d", id)
	return err
}

// VerifyPassword проверяет пароль пользователя
func (s *UserService) VerifyPassword(ctx context.Context, id int, password string) error {
	s.log.Infof("Verifying password for user ID: %d", id)
	err := s.repo.VerifyPassword(ctx, id, password)
	if err != nil {
		s.log.Errorf("Password verification failed for user ID %d: %v", id, err)
		return err
	}
	s.log.Infof("Password verified successfully for user ID %d", id)
	return nil
}
