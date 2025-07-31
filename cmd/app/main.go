package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/alex-pyslar/online-store/internal/config"
	"github.com/alex-pyslar/online-store/internal/handler"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/alex-pyslar/online-store/internal/service"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	//_ "online-store/docs" // Для Swagger
)

// @title Online Store API
// @version 1.0
// @description REST API for an ecommerce platform
// @host localhost:8080
// @BasePath /
func main() {
	// Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}
	defer cfg.DB.Close()
	defer cfg.Redis.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(cfg.DB, cfg.Redis)
	productRepo := repository.NewProductRepository(cfg.DB, cfg.Redis)
	categoryRepo := repository.NewCategoryRepository(cfg.DB, cfg.Redis)
	orderRepo := repository.NewOrderRepository(cfg.DB, cfg.Redis)
	commentRepo := repository.NewCommentRepository(cfg.DB, cfg.Redis)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	orderService := service.NewOrderService(orderRepo)
	commentService := service.NewCommentService(commentRepo)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	orderHandler := handler.NewOrderHandler(orderService)
	commentHandler := handler.NewCommentHandler(commentService)
	authHandler := handler.NewAuthHandler(cfg.OAuthConfig, userService)

	// Настройка роутера
	router := mux.NewRouter()
	// Авторизация
	router.HandleFunc("/auth/google/login", authHandler.HandleGoogleLogin).Methods("GET")
	router.HandleFunc("/auth/google/callback", authHandler.HandleGoogleCallback).Methods("GET")

	// Пользователи
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", userHandler.ListUsers).Methods("GET")

	// Товары
	router.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	router.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	router.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Категории
	router.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	router.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Заказы
	router.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	router.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	router.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	router.HandleFunc("/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")

	// Комментарии
	router.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	router.HandleFunc("/comments/{id}", commentHandler.GetComment).Methods("GET")
	router.HandleFunc("/comments", commentHandler.ListComments).Methods("GET")
	router.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	router.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")

	// Метрики Prometheus
	router.Handle("/metrics", promhttp.Handler())

	// Swagger документация
	router.PathPrefix("/swagger/").Handler(ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Логирование старта сервера с текущей датой и временем
	currentTime := time.Now()
	log.Printf("Server starting on :8080 at %s", currentTime.Format("2006-01-02 15:04:05 MST"))
	log.Fatal(http.ListenAndServe(":8080", router))
}
