package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/alex-pyslar/online-store/internal/config"
	"github.com/alex-pyslar/online-store/internal/handler"
	"github.com/alex-pyslar/online-store/internal/logger"
	"github.com/alex-pyslar/online-store/internal/repository"
	"github.com/alex-pyslar/online-store/internal/service"

	_ "github.com/alex-pyslar/online-store/docs"
)

// @title Online Store API
// @version 1.0
// @description REST API for an ecommerce platform
// @host localhost:8080
// @BasePath /api
func main() {
	// Инициализация логгера
	log, err := logger.NewLogger()
	if err != nil {
		panic(err) // или используйте другой способ обработки ошибки
	}

	// Инициализация конфигурации
	cfg, err := config.NewConfig(log)
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}
	defer cfg.DB.Close()
	defer cfg.Redis.Close()

	// Логирование успешного подключения
	log.Info("Connected to PostgreSQL and Redis successfully")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(cfg.DB, cfg.Redis, log)
	productRepo := repository.NewProductRepository(cfg.DB, cfg.Redis, log)
	categoryRepo := repository.NewCategoryRepository(cfg.DB, cfg.Redis, log)
	orderRepo := repository.NewOrderRepository(cfg.DB, cfg.Redis, log)
	commentRepo := repository.NewCommentRepository(cfg.DB, cfg.Redis, log)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, log)
	productService := service.NewProductService(productRepo, log)
	categoryService := service.NewCategoryService(categoryRepo, log)
	orderService := service.NewOrderService(orderRepo, log)
	commentService := service.NewCommentService(commentRepo, log)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService, log)
	productHandler := handler.NewProductHandler(productService, log)
	categoryHandler := handler.NewCategoryHandler(categoryService, log)
	orderHandler := handler.NewOrderHandler(orderService, log)
	commentHandler := handler.NewCommentHandler(commentService, log)
	authHandler := handler.NewAuthHandler(userService, log)

	// Настройка роутера
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Применяем JWTMiddleware к защищенным маршрутам
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authHandler.JWTMiddleware)

	// Аутентификация
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Пользователи
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	protected.HandleFunc("/users", userHandler.ListUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Товары
	protected.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	protected.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	protected.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	protected.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	protected.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Категории
	protected.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	protected.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	protected.HandleFunc("/categories", categoryHandler.ListCategories).Methods("GET")
	protected.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")

	// Заказы
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	protected.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	protected.HandleFunc("/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	protected.HandleFunc("/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")

	// Комментарии
	protected.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
	protected.HandleFunc("/comments/{id}", commentHandler.GetComment).Methods("GET")
	protected.HandleFunc("/comments", commentHandler.ListComments).Methods("GET")
	protected.HandleFunc("/comments/{id}", commentHandler.UpdateComment).Methods("PUT")
	protected.HandleFunc("/comments/{id}", commentHandler.DeleteComment).Methods("DELETE")

	// Метрики Prometheus
	router.Handle("/metrics", promhttp.Handler())

	// Маршрут для Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Логирование старта сервера
	log.Infof("Server starting on :8080 at %s", time.Now().Format("2006-01-02 15:04:05"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
