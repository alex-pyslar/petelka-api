package tests

import (
	"context"
	"database/sql"
	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/alex-pyslar/petelka-api/internal/service"
	"github.com/redis/go-redis/v9"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/ecommerce_test?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	teardown := func() {
		db.Close()
	}
	return db, teardown
}

func setupTestRedis(t *testing.T) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return redisClient
}

func TestCreateUser(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user := &models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		OAuthID:  "test_oauth_id",
		Password: "password123",
	}

	err := userService.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Password) // Проверяем, что пароль захеширован
}

func TestGetUser(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user := &models.User{
		Email:    "test2@example.com",
		Name:     "Test User 2",
		OAuthID:  "test_oauth_id_2",
		Password: "password123",
	}
	userService.CreateUser(context.Background(), user)

	fetchedUser, err := userService.GetUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
	assert.NotEmpty(t, fetchedUser.Password)
}

func TestListUsers(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user1 := &models.User{Email: "user1@example.com", Name: "User 1", OAuthID: "oauth1", Password: "pass1"}
	user2 := &models.User{Email: "user2@example.com", Name: "User 2", OAuthID: "oauth2", Password: "pass2"}
	userService.CreateUser(context.Background(), user1)
	userService.CreateUser(context.Background(), user2)

	users, err := userService.ListUsers(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 2)
}

func TestUpdateUser(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user := &models.User{Email: "update@example.com", Name: "Update User", OAuthID: "update_oauth", Password: "oldpass"}
	userService.CreateUser(context.Background(), user)

	user.Name = "Updated User"
	user.Password = "newpass"
	err := userService.UpdateUser(context.Background(), user)
	assert.NoError(t, err)

	fetchedUser, err := userService.GetUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", fetchedUser.Name)
	assert.NotEqual(t, "newpass", fetchedUser.Password) // Проверяем, что пароль захеширован
}

func TestDeleteUser(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user := &models.User{Email: "delete@example.com", Name: "Delete User", OAuthID: "delete_oauth", Password: "pass123"}
	userService.CreateUser(context.Background(), user)

	err := userService.DeleteUser(context.Background(), user.ID)
	assert.NoError(t, err)

	_, err = userService.GetUser(context.Background(), user.ID)
	assert.Error(t, err)
}

func TestVerifyPassword(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db, redisClient)
	userService := service.NewUserService(userRepo)

	user := &models.User{Email: "verify@example.com", Name: "Verify User", Password: "password123"}
	userService.CreateUser(context.Background(), user)

	err := userService.VerifyPassword(context.Background(), user.ID, "password123")
	assert.NoError(t, err)

	err = userService.VerifyPassword(context.Background(), user.ID, "wrongpass")
	assert.Error(t, err)
}
