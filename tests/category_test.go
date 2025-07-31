package tests

import (
	"context"
	"github.com/alex-pyslar/online-store/internal/models"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/alex-pyslar/online-store/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Test Category",
	}

	err := categoryService.CreateCategory(context.Background(), category)
	assert.NoError(t, err)
	assert.NotZero(t, category.ID)
}

func TestGetCategory(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	category := &models.Category{
		Name: "Test Category 2",
	}
	categoryService.CreateCategory(context.Background(), category)

	fetchedCategory, err := categoryService.GetCategory(context.Background(), category.ID)
	assert.NoError(t, err)
	assert.Equal(t, category.Name, fetchedCategory.Name)
}

func TestListCategories(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	category1 := &models.Category{Name: "Category 1"}
	category2 := &models.Category{Name: "Category 2"}
	categoryService.CreateCategory(context.Background(), category1)
	categoryService.CreateCategory(context.Background(), category2)

	categories, err := categoryService.ListCategories(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(categories), 2)
}

func TestUpdateCategory(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	category := &models.Category{Name: "Update Category"}
	categoryService.CreateCategory(context.Background(), category)

	category.Name = "Updated Category"
	err := categoryService.UpdateCategory(context.Background(), category)
	assert.NoError(t, err)

	fetchedCategory, err := categoryService.GetCategory(context.Background(), category.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Category", fetchedCategory.Name)
}

func TestDeleteCategory(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	categoryRepo := repository.NewCategoryRepository(db, redisClient)
	categoryService := service.NewCategoryService(categoryRepo)

	category := &models.Category{Name: "Delete Category"}
	categoryService.CreateCategory(context.Background(), category)

	err := categoryService.DeleteCategory(context.Background(), category.ID)
	assert.NoError(t, err)

	_, err = categoryService.GetCategory(context.Background(), category.ID)
	assert.Error(t, err)
}
