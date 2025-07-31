package service

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) VerifyPassword(ctx context.Context, id int, password string) error {
	return s.repo.VerifyPassword(ctx, id, password)
}
