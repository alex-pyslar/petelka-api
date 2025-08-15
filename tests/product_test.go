package tests

import (
	"context"
	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/alex-pyslar/petelka-api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(db, redisClient)
	productService := service.NewProductService(productRepo)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		CategoryID:  1,
	}

	err := productService.CreateProduct(context.Background(), product)
	assert.NoError(t, err)
	assert.NotZero(t, product.ID)
}

func TestGetProduct(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(db, redisClient)
	productService := service.NewProductService(productRepo)

	product := &models.Product{
		Name:        "Test Product 2",
		Description: "Test Description 2",
		Price:       200.0,
		CategoryID:  1,
	}
	productService.CreateProduct(context.Background(), product)

	fetchedProduct, err := productService.GetProduct(context.Background(), product.ID)
	assert.NoError(t, err)
	assert.Equal(t, product.Name, fetchedProduct.Name)
}

func TestListProducts(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(db, redisClient)
	productService := service.NewProductService(productRepo)

	product1 := &models.Product{Name: "Product 1", Description: "Desc 1", Price: 100.0, CategoryID: 1}
	product2 := &models.Product{Name: "Product 2", Description: "Desc 2", Price: 200.0, CategoryID: 1}
	productService.CreateProduct(context.Background(), product1)
	productService.CreateProduct(context.Background(), product2)

	products, err := productService.ListProducts(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(products), 2)
}

func TestUpdateProduct(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(db, redisClient)
	productService := service.NewProductService(productRepo)

	product := &models.Product{Name: "Update Product", Description: "Old Desc", Price: 100.0, CategoryID: 1}
	productService.CreateProduct(context.Background(), product)

	product.Description = "New Desc"
	err := productService.UpdateProduct(context.Background(), product)
	assert.NoError(t, err)

	fetchedProduct, err := productService.GetProduct(context.Background(), product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "New Desc", fetchedProduct.Description)
}

func TestDeleteProduct(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	productRepo := repository.NewProductRepository(db, redisClient)
	productService := service.NewProductService(productRepo)

	product := &models.Product{Name: "Delete Product", Description: "Delete Desc", Price: 100.0, CategoryID: 1}
	productService.CreateProduct(context.Background(), product)

	err := productService.DeleteProduct(context.Background(), product.ID)
	assert.NoError(t, err)

	_, err = productService.GetProduct(context.Background(), product.ID)
	assert.Error(t, err)
}
