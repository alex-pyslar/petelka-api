package tests

import (
	"context"
	"github.com/alex-pyslar/petelka-api/internal/models"
	"github.com/alex-pyslar/petelka-api/internal/repository"
	"github.com/alex-pyslar/petelka-api/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	orderRepo := repository.NewOrderRepository(db, redisClient)
	orderService := service.NewOrderService(orderRepo)

	order := &models.Order{
		UserID: 1,
		Total:  100.0,
		Status: "pending",
	}

	err := orderService.CreateOrder(context.Background(), order)
	assert.NoError(t, err)
	assert.NotZero(t, order.ID)
}

func TestGetOrder(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	orderRepo := repository.NewOrderRepository(db, redisClient)
	orderService := service.NewOrderService(orderRepo)

	order := &models.Order{
		UserID: 1,
		Total:  200.0,
		Status: "pending",
	}
	orderService.CreateOrder(context.Background(), order)

	fetchedOrder, err := orderService.GetOrder(context.Background(), order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order.Total, fetchedOrder.Total)
}

func TestListOrders(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	orderRepo := repository.NewOrderRepository(db, redisClient)
	orderService := service.NewOrderService(orderRepo)

	order1 := &models.Order{UserID: 1, Total: 100.0, Status: "pending"}
	order2 := &models.Order{UserID: 1, Total: 200.0, Status: "pending"}
	orderService.CreateOrder(context.Background(), order1)
	orderService.CreateOrder(context.Background(), order2)

	orders, err := orderService.ListOrders(context.Background())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(orders), 2)
}

func TestUpdateOrder(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	orderRepo := repository.NewOrderRepository(db, redisClient)
	orderService := service.NewOrderService(orderRepo)

	order := &models.Order{UserID: 1, Total: 100.0, Status: "pending"}
	orderService.CreateOrder(context.Background(), order)

	order.Status = "completed"
	err := orderService.UpdateOrder(context.Background(), order)
	assert.NoError(t, err)

	fetchedOrder, err := orderService.GetOrder(context.Background(), order.ID)
	assert.NoError(t, err)
	assert.Equal(t, "completed", fetchedOrder.Status)
}

func TestDeleteOrder(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	orderRepo := repository.NewOrderRepository(db, redisClient)
	orderService := service.NewOrderService(orderRepo)

	order := &models.Order{UserID: 1, Total: 100.0, Status: "pending"}
	orderService.CreateOrder(context.Background(), order)

	err := orderService.DeleteOrder(context.Background(), order.ID)
	assert.NoError(t, err)

	_, err = orderService.GetOrder(context.Background(), order.ID)
	assert.Error(t, err)
}
